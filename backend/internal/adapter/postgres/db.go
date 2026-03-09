package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewDB(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping database:", err)
	}

	log.Println("connected to database")
	return db
}
