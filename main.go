package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"time"

	"github.com/felipebool/skid/window"
)

var size = flag.Int("size", 60, "window size")
var path = flag.String("path", ".window.json", "path to persist window")

//var restore = flag.Bool("restore", false, "restore state from file")

func counter(sliding *window.Window) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("request received")
		epoch := time.Now().Unix()
		sliding.Add(epoch)
		fmt.Fprintf(w, "%d", sliding.Count(epoch))
	}
}

func run(windowSize int) error {
	shutdown := make(chan os.Signal, 1)
	serverError := make(chan error)

	// signal trap to gracefully stop application
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	sliding := window.New(int64(windowSize))

	// avoid double request when testing on browser
	http.HandleFunc("/favicon.ico", func(rw http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/", counter(sliding))

	// keep main goroutine clean to receive errors
	go func() {
		serverError <- http.ListenAndServe(":8080", nil)
	}()

	select {
	case err := <-serverError:
		return err
	case <-shutdown:
		log.Println("shutdown started")
		log.Println("trying to persist window to file")
		if err := sliding.Persist(*path); err != nil {
			log.Println("unable to persist window to file")
			return err
		}
		log.Println("window persisted successfully to", *path)
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
