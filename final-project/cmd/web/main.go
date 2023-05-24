package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

func main() {
	// connect to a database
	db := initDB()
	db.Ping()
	// create sessions

	// create channels

	// create a waitgroup

	// set the application config

	// set up mail

	// listen to web connection

}

func initDB() *sql.DB {
	// try to connect to db repeatedly if necessarily
	conn := connectToDB()

	if conn == nil {
		log.Panic() // can't resolve
	}

	return conn
}

func connectToDB() *sql.DB {
	// try to connect a fixed number of time, give up if things fails

	counts := 0

	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet...")

		} else {
			log.Println("Connected to database.")
			return connection
		}

		if counts > 10 {
			return nil
		}

		log.Println("Backing off for 1 second.")
		time.Sleep(1 * time.Second)

		counts++
		continue
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
