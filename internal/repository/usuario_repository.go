package repository

import (
	"database/sql"

	"github.com/Graynie/InkZen/internal/models"
)

func CreateUser(db *sql.DB, user models.Usuario) error {
	query := `
	INSERT INTO usuarios (nombre, email, password)
	VALUES (?, ?, ?);
	`

	_, err := db.Exec(query, user.Nombre, user.Email, user.Password)
	return err
}
func GetUserByEmail(db *sql.DB, email string) (models.Usuario, error) {
	var user models.Usuario

	query := "SELECT id, nombre, email, password FROM usuarios WHERE email = ?"

	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Nombre,
		&user.Email,
		&user.Password,
	)

	return user, err
}
