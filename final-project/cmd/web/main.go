package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
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
	session := initSession()
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

// initSession returns a session manager
func initSession() *scs.SessionManager {

	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour              // session lasts for one day
	session.Cookie.Persist = true                  // let cookies persist after sessions are closed
	session.Cookie.SameSite = http.SameSiteLaxMode // what does this do?
	session.Cookie.Secure = true

	return session
}

func initRedis() *redis.Pool {

	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS")) // variable specified in makefil
		},
	}

	return redisPool
}
