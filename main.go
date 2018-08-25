package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Server represents the server environment (router and database)
type Server struct {
	DB   *DB
	http *http.Server
}

func main() {

	// get flag values
	var port int
	var dbhost string
	var dbuser string
	var dbpass string
	var dbname string

	flag.IntVar(&port, "port", 8000, "the port to listen for http connections on")
	flag.StringVar(&dbhost, "db-host", "localhost", "database hostname (default: localhost)")
	flag.StringVar(&dbuser, "db-user", "", "database username")
	flag.StringVar(&dbpass, "db-pass", "", "database user password")
	flag.StringVar(&dbname, "db-name", "", "database name")
	flag.Parse()

	// wait for database connection (see database.go)
	db, err := waitForDB(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbuser, dbpass, dbhost, dbname))
	if err != nil {
		log.Panic("Database connection error:", err)
	}

	// create a pointer to a new http server on specified port
	httpServer := &http.Server{Addr: fmt.Sprintf(":%v", port), Handler: nil}

	api := &Server{db, httpServer}

	// start http server
	log.Printf("Starting HTTP server on port %v.", port)
	log.Println("Press CTRL+C to stop.")
	http.HandleFunc("/health", api.health)
	log.Fatal(httpServer.ListenAndServe())
}

// health is a simple handler that returns 200 OK status and text "OK"
// it can be used for readiness/health checks
func (*Server) health(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, http.StatusText(405), 405)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}
