package builder

import (
	"context"
	"hotreload/internal/logx"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"time"
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

	// cancel previous build if running
	if b.cancel != nil {
		b.logger.Warn("canceling previous build")
		b.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel

	b.mu.Unlock()

	start := time.Now()

	logx.Build("building project")

	cmd := exec.CommandContext(ctx, "cmd", "/C", b.command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	// clear cancel after build finishes
	b.mu.Lock()
	b.cancel = nil
	b.mu.Unlock()

	if ctx.Err() == context.Canceled {
		b.logger.Warn("build canceled due to new changes")
		return ctx.Err()
	}

	if err != nil {
		logx.Error("build failed")
		b.logger.Error("build failed", "error", err)
		return err
	}

	duration := time.Since(start)

	logx.Build("build completed")
	b.logger.Info("build finished", "duration", duration)

	return nil
}