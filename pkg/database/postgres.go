package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgres(databaseURL string) *sql.DB {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("database connected")
	return db
}
