package builder

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"sync"
)

type Builder struct {
	command string
	logger  *slog.Logger

	mu     sync.Mutex
	cancel context.CancelFunc
}

func New(cmd string, logger *slog.Logger) *Builder {
	return &Builder{
		command: cmd,
		logger:  logger,
	}
}

func (b *Builder) Build() error {

	b.mu.Lock()

	if b.cancel != nil {
		b.logger.Warn("canceling previous build")
		b.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel

	b.mu.Unlock()

	b.logger.Info("building project")

	cmd := exec.CommandContext(ctx, "cmd", "/C", b.command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if ctx.Err() == context.Canceled {
		b.logger.Warn("build canceled due to new changes")
		return ctx.Err()
	}

	if err != nil {
		b.logger.Error("build failed", "error", err)
		return err
	}

	b.logger.Info("build completed")

	return nil
}