package window

import (
	"os"
	"sync"
)

type Window struct {
	mu     *sync.Mutex
	size   int64
	window map[int64]int64
}

func (w *Window) Count(epoch int64) int64 {
	var sum int64
	start := epoch - int64(w.size)

	w.mu.Lock()
	// get total requests inside window
	for timestamp, total := range w.window {
		if timestamp > start {
			sum += total
			continue
		}
		// reduce memory footprint
		delete(w.window, timestamp)
	}
	w.mu.Unlock()
	return sum
}

func (w *Window) Add(epoch int64) {
	w.mu.Lock()
	w.window[epoch]++
	w.mu.Unlock()
}

func (w *Window) Persist(path string) error {
	w.mu.Lock()

	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	w.mu.Unlock()
	return nil
}

func New(size int64) *Window {
	return &Window{
		mu:     &sync.Mutex{},
		size:   size,
		window: make(map[int64]int64),
	}
}
