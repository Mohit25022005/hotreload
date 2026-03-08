package runner

import (
	"hotreload/internal/logx"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type Runner struct {
	command   string
	cmd       *exec.Cmd
	logger    *slog.Logger
	lastStart time.Time
	mu        sync.Mutex
}

func New(command string, logger *slog.Logger) *Runner {
	return &Runner{
		command: command,
		logger:  logger,
	}
}

func (r *Runner) Restart() {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.stopLocked()

	// crash-loop protection
	if time.Since(r.lastStart) < 2*time.Second {
		r.logger.Warn("server restarted too quickly, delaying")
		time.Sleep(2 * time.Second)
	}

	logx.Server("starting server")

	cmd := exec.Command(r.command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		logx.Error("failed to start server")
		r.logger.Error("failed to start server", "error", err)
		return
	}

	r.cmd = cmd
	r.lastStart = time.Now()

	// monitor process exit
	go func() {
		err := cmd.Wait()
		if err != nil {
			r.logger.Warn("server exited with error", "error", err)
		} else {
			r.logger.Info("server stopped")
		}
	}()
}

func (r *Runner) Stop() {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.stopLocked()
}

func (r *Runner) stopLocked() {

	if r.cmd == nil || r.cmd.Process == nil {
		return
	}

	logx.Server("stopping server")

	// try graceful shutdown first
	err := r.cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		r.logger.Warn("failed to send SIGTERM, killing process", "error", err)
		r.forceKill()
		return
	}

	done := make(chan error, 1)

	go func() {
		done <- r.cmd.Wait()
	}()

	select {

	case <-time.After(3 * time.Second):
		r.logger.Warn("server did not stop in time, forcing kill")
		r.forceKill()

	case <-done:
		r.logger.Info("server stopped gracefully")
	}

	r.cmd = nil
}

func (r *Runner) forceKill() {

	if r.cmd == nil || r.cmd.Process == nil {
		return
	}

	err := r.cmd.Process.Kill()
	if err != nil {
		r.logger.Warn("failed to kill process", "error", err)
	}

	_, _ = r.cmd.Process.Wait()
}