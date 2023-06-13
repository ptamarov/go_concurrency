package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/ptamarov/go_basics/concurrency/final-project/data"
)

// type to share the application configuration
type Config struct {
	Session  *scs.SessionManager
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
	Models   data.Models
}
