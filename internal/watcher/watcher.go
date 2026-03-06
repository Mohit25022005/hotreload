package watcher

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	root    string
	logger  *slog.Logger
	watcher *fsnotify.Watcher
	Events  chan struct{}
}

func New(root string, logger *slog.Logger) (*Watcher, error) {

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	wr := &Watcher{
		root:    root,
		logger:  logger,
		watcher: w,
		Events:  make(chan struct{}, 1),
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() && !shouldIgnore(path) {
			w.Add(path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return wr, nil
}

func (w *Watcher) Start() {

	for {

		select {

		case event := <-w.watcher.Events:

			if shouldIgnore(event.Name) {
				continue
			}

			if !shouldWatch(event.Name) {
				continue
			}

			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 {

				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					w.logger.Info("new directory detected", "path", event.Name)
					w.watcher.Add(event.Name)
				}

				select {
				case w.Events <- struct{}{}:
				default:
				}
			}

		case err := <-w.watcher.Errors:
			w.logger.Error("watch error", "error", err)
		}
	}
}

func shouldWatch(path string) bool {

	ext := filepath.Ext(path)

	switch ext {
	case ".go", ".mod", ".sum":
		return true
	}

	return false
}

func shouldIgnore(path string) bool {

	ignore := []string{
		".git",
		"node_modules",
		"bin",
		"tmp",
	}

	for _, dir := range ignore {
		if strings.Contains(path, dir) {
			return true
		}
	}

	if strings.HasSuffix(path, "~") ||
		strings.HasSuffix(path, ".swp") ||
		strings.HasSuffix(path, ".tmp") {
		return true
	}

	return false
}