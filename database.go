package main

import (
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	// load postgres driver
	_ "github.com/lib/pq"
)

// DB represents a database connection
type DB struct {
	*sqlx.DB
}

// waitForDB returns a new database handle (type DB).
// It will wait for a connection to become available
// rather than panic/exit on error.
func waitForDB(connectionConfig string) (*DB, error) {
	var db *DB
	var ver string
	for {
		conn, err := sqlx.Connect("postgres", connectionConfig)

		// if db connection not available, wait and retry
		if err != nil && strings.HasSuffix(err.Error(), ": connection refused") {
			log.Println(err, "- retrying...")
			time.Sleep(5 * time.Second)
			continue
		} else if err != nil {
			log.Panic(err)
		}

		db = &DB{conn}

		ver, err = db.health()
		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}

	log.Printf("Database ready (%s).", ver)
	return db, nil
}

// health performs a simple query against the database and
// returns an error (nil if no error)
func (db *DB) health() (string, error) {
	var version string
	row := db.QueryRow("SELECT VERSION()")
	err := row.Scan(&version)
	return version, err
}
