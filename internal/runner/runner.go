package runner

import (
	"log/slog"
	"os"
	"os/exec"
	"sync"
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

	// crash loop protection
	if time.Since(r.lastStart) < 2*time.Second {
		r.logger.Warn("server restarted too quickly, delaying")
		time.Sleep(2 * time.Second)
	}

	r.logger.Info("starting server")

	cmd := exec.Command(r.command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		r.logger.Error("failed to start server", "error", err)
		return
	}

	r.cmd = cmd
	r.lastStart = time.Now()
}

func (r *Runner) stopLocked() {

	if r.cmd == nil || r.cmd.Process == nil {
		return
	}

	r.logger.Info("stopping server")

	err := r.cmd.Process.Kill()
	if err != nil {
		r.logger.Warn("process already stopped or cannot be killed", "error", err)
	}

	_, _ = r.cmd.Process.Wait()
	r.cmd = nil
}