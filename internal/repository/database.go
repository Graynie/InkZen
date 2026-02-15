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

	// Tabla usuarios
	userTable := `
	CREATE TABLE IF NOT EXISTS usuarios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`

	_, err := db.Exec(userTable)
	if err != nil {
		log.Fatal(err)
	}

	// Tabla mangas (contenido digital)
	mangaTable := `
	CREATE TABLE IF NOT EXISTS mangas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		titulo TEXT NOT NULL,
		autor TEXT NOT NULL,
		genero TEXT,
		idioma TEXT,
		editorial TEXT,
		descripcion TEXT,
		capitulos_tot INTEGER DEFAULT 0,
		disponible BOOLEAN DEFAULT 1
	);
	`

	_, err = db.Exec(mangaTable)
	if err != nil {
		log.Fatal(err)
	}

	// Tabla lecturas (relación usuario - manga)
	lecturaTable := `
	CREATE TABLE IF NOT EXISTS lecturas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		usuario_id INTEGER NOT NULL,
		manga_id INTEGER NOT NULL,
		capitulo_actual INTEGER DEFAULT 0,
		FOREIGN KEY(usuario_id) REFERENCES usuarios(id),
		FOREIGN KEY(manga_id) REFERENCES mangas(id)
	);
	`

	_, err = db.Exec(lecturaTable)
	if err != nil {
		log.Fatal(err)
	}
}
