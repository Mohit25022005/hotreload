package builder

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func TestBuilderInitialization(t *testing.T) {

	logger := newTestLogger()

	b := New("echo hello", logger)

	if b == nil {
		t.Fatal("builder should not be nil")
	}

	if b.command != "echo hello" {
		t.Fatalf("expected command 'echo hello', got %s", b.command)
	}
}

func TestBuildSuccess(t *testing.T) {

	logger := newTestLogger()

	b := New("echo build success", logger)

	err := b.Build()

	if err != nil {
		t.Fatalf("expected build success, got error: %v", err)
	}
}

func TestBuildFailure(t *testing.T) {

	logger := newTestLogger()

	// exit 1 simulates a failing build
	b := New("exit 1", logger)

	err := b.Build()

	if err == nil {
		t.Fatal("expected build failure but got success")
	}
}

func TestBuildCancellation(t *testing.T) {

	logger := newTestLogger()

	// simulate a long-running build
	b := New("ping 127.0.0.1 -n 5 > nul", logger)

	go func() {
		_ = b.Build()
	}()

	time.Sleep(500 * time.Millisecond)

	// start new build which should cancel the previous one
	err := b.Build()

	if err != nil && err.Error() == "context canceled" {
		t.Log("build cancellation worked correctly")
	}
}