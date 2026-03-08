package watcher

import (
	"hotreload/internal/logx"
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

	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if d.IsDir() && !shouldIgnore(path) {

			err := w.Add(path)
			if err != nil {
				logger.Warn("failed to watch directory", "path", path, "error", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return wr, nil
}

func (w *Watcher) Start() {

	logx.Watch("watching for file changes")

	for {

		select {

		case event, ok := <-w.watcher.Events:

			if !ok {
				return
			}

			if shouldIgnore(event.Name) {
				continue
			}

			// handle new directories
			if event.Op&fsnotify.Create != 0 {

				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {

					err := w.watcher.Add(event.Name)
					if err == nil {
						w.logger.Info("new directory detected", "path", event.Name)
					}

					continue
				}
			}

			if !shouldWatch(event.Name) {
				continue
			}

			// editor safe detection
			if event.Op&(fsnotify.Write|
				fsnotify.Create|
				fsnotify.Remove|
				fsnotify.Rename|
				fsnotify.Chmod) != 0 {
				logx.Watch("file change detected")

				select {
				case w.Events <- struct{}{}:
				default:
				}
			}

		case err, ok := <-w.watcher.Errors:

			if !ok {
				return
			}

			w.logger.Error("watch error", "error", err)
		}
	}
}

func (w *Watcher) Close() error {
	return w.watcher.Close()
}

func shouldWatch(path string) bool {

	if shouldIgnore(path) {
		return false
	}

	// ignore editor temporary files
	if strings.HasSuffix(path, "~") ||
		strings.HasSuffix(path, ".swp") ||
		strings.HasSuffix(path, ".tmp") {
		return false
	}

	return true
}

func shouldIgnore(path string) bool {

	ignore := []string{
		".git",
		"node_modules",
		"bin",
		"tmp",
		".idea",
		".vscode",
	}

	for _, dir := range ignore {
		if strings.Contains(path, dir) {
			return true
		}
	}

	if strings.HasSuffix(path, "~") ||
		strings.HasSuffix(path, ".swp") ||
		strings.HasSuffix(path, ".tmp") ||
		strings.HasSuffix(path, ".log") {
		return true
	}

	return false
}
