package main

import (
	"os"
	"time"
)

// Event is information about what file changed
type Event struct {
	Size int64
}

// Watcher watches a given file for changes
type Watcher struct {
	file     string
	lastSize int64
	ticker   *time.Ticker
	Events   chan Event
	Errors   chan error
}

// NewWatcher creates a new watcher
func NewWatcher(pathToFile string, tick time.Duration) *Watcher {
	return &Watcher{
		lastSize: int64(-1),
		file:     pathToFile,
		ticker:   time.NewTicker(tick),
		Events:   make(chan Event),
		Errors:   make(chan error),
	}
}

// Start starts the watcher
func (w *Watcher) Start() {
	go func() {
		w.tick()
		for range w.ticker.C {
			w.tick()
		}
	}()
}

func (w *Watcher) tick() {
	size, err := w.size()
	if err != nil {
		w.Errors <- err
	} else if size != w.lastSize {
		w.lastSize = size
		w.Events <- Event{
			Size: size,
		}
	}
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	w.ticker.Stop()
}

func (w *Watcher) size() (int64, error) {
	fi, err := os.Stat(w.file)
	if err != nil {
		return -1, err
	}
	return fi.Size(), nil
}
