package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"hotreload/internal/builder"
	"hotreload/internal/debounce"
	"hotreload/internal/event"
	"hotreload/internal/runner"
	"hotreload/internal/watcher"
)

func main() {

	root := flag.String("root", ".", "project root directory")
	buildCmd := flag.String("build", "", "build command")
	execCmd := flag.String("exec", "", "run command")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if *buildCmd == "" || *execCmd == "" {
		logger.Error("build and exec commands are required")
		os.Exit(1)
	}

	logger.Info("🔥 HotReload started", "root", *root)

	// components
	builder := builder.New(*buildCmd, logger)
	runner := runner.New(*execCmd, logger)
	debouncer := debounce.New(500)
	queue := event.NewQueue()

	// trigger pipeline
	trigger := func() {

		logger.Info("processing change")

		err := builder.Build()
		if err != nil {
			return
		}

		runner.Restart()
	}

	// watcher
	w, err := watcher.New(*root, logger)
	if err != nil {
		logger.Error("failed to initialize watcher", "error", err)
		os.Exit(1)
	}

	// watcher → queue
	go func() {

		for range w.Events {

			logger.Info("file change detected")

			queue.Push()
		}

	}()

	// queue → debounce → build pipeline
	go func() {

		for range queue.Events() {

			debouncer.Trigger(trigger)

		}

	}()

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {

		<-sig

		logger.Info("shutting down")

		runner.Stop()

		os.Exit(0)

	}()

	// initial build
	trigger()

	// start watcher loop
	w.Start()
}