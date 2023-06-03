package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

func (app *Config) serve() {
	// start http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Println("Starting web server...")

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
func main() {
	// connect to a database
	db := initDB()

	// create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// log.Lshortfile is a flag that records a message coming from a file, and the name of the file in short form

	// create sessions
	session := initSession()

	// create channels
	// later
	// create a waitgroup
	wg := sync.WaitGroup{}

	// set the application config
	app := Config{
		Session:  session,
		DB:       db,
		Wait:     &wg,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	// set up mail

	// listen for signals SIGTERM and SIGINT
	go app.ListenForShutdown()

	// listen to web connection
	app.serve()
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

// Run in the background and listen for a shutdown
func (app *Config) ListenForShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// When you get the interrupt or terminate signal

	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {
	// Perform any clean-up tasks
	app.InfoLog.Println("Would run clean up tasks...")

	// block until waitgroup is empty
	app.Wait.Wait()

	app.InfoLog.Print("Closing channels and shutting down application...")
}
