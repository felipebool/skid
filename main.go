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

var size = flag.Int("size", 60, "window size (default: 60 seconds)")
var path = flag.String("path", ".window.json", "path to persist window (default: .window.json)")
var restore = flag.Bool("restore", false, "restore state from file (default: false)")
var port = flag.Int("port", 8080, "port to use (default: 8080)")

func counter(sliding *window.Window) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("request received")
		epoch := time.Now().Unix()
		sliding.Add(epoch)
		fmt.Fprintf(w, "%d", sliding.Count(epoch))
	}
}

func run(size, port int, path string, restore bool) error {
	shutdown := make(chan os.Signal, 1)
	serverError := make(chan error)

	// signal trap to gracefully stop application
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	sliding := window.New(size, path, restore)

	// avoid double request when testing on browser
	http.HandleFunc("/favicon.ico", func(rw http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/", counter(sliding))

	// keep main goroutine clean to receive errors
	go func() {
		serverError <- http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()

	select {
	case err := <-serverError:
		return err
	case <-shutdown:
		log.Println("shutdown started")
		log.Println("trying to persist window to file")
		if err := sliding.Persist(); err != nil {
			log.Println("unable to persist window to file")
			return err
		}
		log.Println("window persisted successfully to", path)
	}
	return nil
}

func main() {
	flag.Parse()
	log.Println("server started")
	if err := run(*size, *port, *path, *restore); err != nil {
		log.Fatal(err)
	}
}
