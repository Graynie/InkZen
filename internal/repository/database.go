package repository

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./db/inkzen.db")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
