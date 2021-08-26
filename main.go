package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var size = flag.Int("size", 60, "window size")

// add signal to save to file when ctrl-C
type window struct {
	mu     *sync.Mutex
	size   int64
	window map[int64]int64
}

func new(size int64) *window {
	return &window{
		mu:     &sync.Mutex{},
		size:   size,
		window: make(map[int64]int64),
	}
}

func (w *window) count(epoch int64) int64 {
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

func (w *window) add(epoch int64) {
	w.mu.Lock()
	w.window[epoch]++
	w.mu.Unlock()
}

func counter(sliding *window) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		epoch := time.Now().Unix()
		sliding.add(epoch)
		fmt.Fprintf(w, "%d", sliding.count(epoch))
	}
}

func run(windowSize int) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	http.HandleFunc("/", counter(new(int64(windowSize))))
	return http.ListenAndServe(":8080", nil)
}

func main() {
	flag.Parse()
	if err := run(*size); err != nil {
		log.Fatal(err)
	}
}
