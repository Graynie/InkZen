package repository

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func NewDatabase() *sql.DB {
	db, err := sql.Open("sqlite", "./db/inkzen.db")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func InitSchema(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS usuarios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
