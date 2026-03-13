package postgres

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewDB(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(3 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping database:", err)
	}

	log.Println("connected to database")
	return db
}
