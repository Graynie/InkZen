package repository

import (
	"database/sql"

	"github.com/Graynie/InkZen/internal/models"
)

func CrearLectura(db *sql.DB, lectura models.Lectura) error {
	query := `
	INSERT INTO lecturas (usuario_id, manga_id, capitulo_actual)
	VALUES (?, ?, ?)
	`
	_, err := db.Exec(query, lectura.UsuarioID, lectura.MangaID, lectura.CapituloActual)
	return err
}

func ActualizarCapitulo(db *sql.DB, usuarioID int, mangaID int, capitulo int) error {
	query := `
	UPDATE lecturas
	SET capitulo_actual = ?
	WHERE usuario_id = ? AND manga_id = ?
	`
	_, err := db.Exec(query, capitulo, usuarioID, mangaID)
	return err
}

func ObtenerLecturasUsuario(db *sql.DB, usuarioID int) ([]models.Lectura, error) {
	rows, err := db.Query(`
		SELECT id, usuario_id, manga_id, capitulo_actual
		FROM lecturas
		WHERE usuario_id = ?
	`, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lecturas []models.Lectura

	for rows.Next() {
		var l models.Lectura
		err := rows.Scan(&l.ID, &l.UsuarioID, &l.MangaID, &l.CapituloActual)
		if err != nil {
			return nil, err
		}
		lecturas = append(lecturas, l)
	}

	return lecturas, nil
}
