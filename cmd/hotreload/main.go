package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"hotreload/internal/builder"
	"hotreload/internal/debounce"
	"hotreload/internal/runner"
	"hotreload/internal/watcher"
)

func main() {

	root := flag.String("root", ".", "project root")
	buildCmd := flag.String("build", "", "build command")
	execCmd := flag.String("exec", "", "run command")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if *buildCmd == "" || *execCmd == "" {
		logger.Error("build and exec commands required")
		os.Exit(1)
	}

	buildManager := builder.New(*buildCmd, logger)
	runner := runner.New(*execCmd, logger)

	debouncer := debounce.New(500)

	trigger := func() {

		logger.Info("change detected")

		ctx := context.Background()

		err := buildManager.Build(ctx)
		if err != nil {
			logger.Error("build failed", "error", err)
			return
		}

		runner.Restart()
	}

	w, err := watcher.New(*root, trigger, logger)
	if err != nil {
		logger.Error("watcher failed", "error", err)
		os.Exit(1)
	}

	go func() {
		for range w.Events {
			debouncer.Trigger(trigger)
		}
	}()

	// initial build
	trigger()

	w.Start()
}