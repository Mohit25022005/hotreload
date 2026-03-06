package builder

import (
	"log/slog"
	"os"
	"testing"
)

func TestBuilderInitialization(t *testing.T) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	b := New("echo test", logger)

	if b == nil {
		t.Fatal("builder should not be nil")
	}
}