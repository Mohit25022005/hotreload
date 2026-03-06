package runner

import (
	"log/slog"
	"os"
	"os/exec"
)

type Runner struct {
	command string
	cmd     *exec.Cmd
	logger  *slog.Logger
}

func New(cmd string, logger *slog.Logger) *Runner {
	return &Runner{
		command: cmd,
		logger:  logger,
	}
}

func (r *Runner) Restart() {

	r.Stop()

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
}

func (r *Runner) Stop() {

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