package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"sync"
	"time"
)

var db *sqlx.DB
var onceDb sync.Once

func LoadDBConnection() *sqlx.DB {
	onceDb.Do(func() {
		db = initializeDb()
	})

	return db
}

func initializeDb() *sqlx.DB {
	DriverName := os.Getenv("DRIVER_NAME")
	DataSource := os.Getenv("DATASOURCE_URL")
	db, dbError := sqlx.Open(DriverName, DataSource)
	if dbError != nil {
		panic(dbError)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if dbError = db.Ping(); dbError != nil {
		panic(dbError)
	}

	log.Println("Database connection is successful.")

	return db
}
