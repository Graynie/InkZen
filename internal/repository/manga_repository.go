package repository

import (
	"database/sql"

	"github.com/Graynie/InkZen/internal/models"
)

func CreateManga(db *sql.DB, manga models.Manga) error {
	query := `
	INSERT INTO mangas 
	(titulo, autor, genero, idioma, editorial, descripcion, capitulos_tot, disponible)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(
		query,
		manga.Titulo,
		manga.Autor,
		manga.Genero,
		manga.Idioma,
		manga.Editorial,
		manga.Descripcion,
		manga.CapitulosTot,
		manga.Disponible,
	)

	return err
}

func GetAllMangas(db *sql.DB) ([]models.Manga, error) {
	rows, err := db.Query(`
		SELECT id, titulo, autor, genero, idioma, editorial, descripcion, capitulos_tot, disponible
		FROM mangas
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mangas []models.Manga

	for rows.Next() {
		var m models.Manga
		err := rows.Scan(
			&m.ID,
			&m.Titulo,
			&m.Autor,
			&m.Genero,
			&m.Idioma,
			&m.Editorial,
			&m.Descripcion,
			&m.CapitulosTot,
			&m.Disponible,
		)
		if err != nil {
			return nil, err
		}

		mangas = append(mangas, m)
	}

	return mangas, nil
}
