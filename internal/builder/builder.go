package builder

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
)

type Builder struct {
	command string
	logger  *slog.Logger
}

func New(cmd string, logger *slog.Logger) *Builder {
	return &Builder{
		command: cmd,
		logger:  logger,
	}
}

func (b *Builder) Build(ctx context.Context) error {

	b.logger.Info("building project")

	cmd := exec.CommandContext(ctx, "cmd", "/C", b.command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}