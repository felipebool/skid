package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"sync"
	"time"
)

var size = flag.Int("size", 60, "window size")
var path = flag.String("path", ".window.json", "path to persist window")

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

func (w *window) persist(path string) error {
	log.Println("window persisted successfully to", path)
	return nil
}

func counter(sliding *window) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("request received")
		epoch := time.Now().Unix()
		sliding.add(epoch)
		fmt.Fprintf(w, "%d", sliding.count(epoch))
	}
}

func run(windowSize int) error {
	shutdown := make(chan os.Signal, 1)
	serverError := make(chan error)

	// signal trap to gracefully stop application
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	sliding := new(int64(windowSize))

	// avoid double request when testing on browser
	http.HandleFunc("/favicon.ico", func(rw http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/", counter(sliding))

	// start server in another goroutine
	go func() {
		serverError <- http.ListenAndServe(":8080", nil)
	}()

	select {
	case err := <-serverError:
		return err
	case <-shutdown:
		log.Println("shutdown started")
		log.Println("trying to persist window to file")
		// try to persist data to file
		if err := sliding.persist(*path); err != nil {
			log.Println("unable to persist to file")
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	log.Println("server started")
	if err := run(*size); err != nil {
		log.Fatal(err)
	}
}
