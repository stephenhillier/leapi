package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
}

func main() {
	// define the port that the server will run on
	port := 8000

	// set up a channel that accepts an interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	s := http.Server{Addr: ":8000", Handler: nil}

	// start server and continue
	go func() {
		log.Printf("Starting HTTP server on port %v.", port)
		log.Println("Press CTRL+C to stop.")
		http.HandleFunc("/health", health)
		log.Fatal(s.ListenAndServe())
	}()

	// server is running, now wait for an interrupt signal
	<-stop

	// interrupt signal received: stop server.
	log.Println("Stopping server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Shutdown(ctx)
	log.Println("Server shutdown. Come again soon!")
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
