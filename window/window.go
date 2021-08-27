package window

import (
	"encoding/json"
	"os"
	"sync"
)

type Window struct {
	mu     *sync.Mutex
	size   int64
	path   string
	window map[int64]int64
}

// Count sums the number of requests in the window
func (w *Window) Count(epoch int64) int64 {
	var sum int64
	start := epoch - w.size

	w.mu.Lock()
	defer w.mu.Unlock()
	for timestamp, total := range w.window {
		if timestamp > start {
			sum += total
			continue
		}
		// reduce memory footprint
		delete(w.window, timestamp)
	}
	return sum
}

// Add increments the counter for the specific second inside the window
func (w *Window) Add(epoch int64) {
	w.mu.Lock()
	w.window[epoch]++
	w.mu.Unlock()
}

// Persist persists the window to a file
func (w *Window) Persist() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	data, err := json.Marshal(w.window)
	if err != nil {
		return err
	}
	if err = os.WriteFile(w.path, data, 0644); err != nil {
		return err
	}
	return nil
}

func (w *Window) restore() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	data, err := os.ReadFile(w.path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &w.window); err != nil {
		return err
	}
	return nil
}

func New(size int, path string, restore bool) *Window {
	w := &Window{
		mu:     &sync.Mutex{},
		size:   int64(size),
		path:   path,
		window: make(map[int64]int64),
	}
	if restore {
		w.restore()
	}
	return w
}
