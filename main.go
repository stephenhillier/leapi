package main

import (
	"io"
	"log"
	"net/http"
)

func main() {

	port := 8000

	log.Printf("Starting HTTP server on port %v.", port)
	log.Println("Press CTRL+C to stop.")
	http.HandleFunc("/health", health)
	log.Fatal(http.ListenAndServe(":8000", nil))
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
