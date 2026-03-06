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

func New(root string, trigger func(), logger *slog.Logger) (*Watcher, error) {

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

			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 {

				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
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

func shouldIgnore(path string) bool {

	ignore := []string{
		".git",
		"node_modules",
		"bin",
		".swp",
		".tmp",
	}

	for _, i := range ignore {
		if strings.Contains(path, i) {
			return true
		}
	}

	return false
}