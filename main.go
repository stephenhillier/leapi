package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	port := 8000

	// add a channel that accepts an interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	server := http.Server{Addr: fmt.Sprintf(":%v", port), Handler: nil}

	// start server and continue
	go func() {
		log.Printf("Starting HTTP server on port %v.", port)
		log.Println("Press CTRL+C to stop.")
		http.HandleFunc("/health", health)
		log.Fatal(server.ListenAndServe())
	}()

	// wait for interrupt signal
	<-stop

	// interrupt signal received- shut down http server.
	log.Println("Stopping server.")
	server.Shutdown(context.Background())

}

func health(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, http.StatusText(405), 405)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}
