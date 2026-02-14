package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/Graynie/InkZen/internal/models"
	"github.com/Graynie/InkZen/internal/repository"
	"github.com/Graynie/InkZen/internal/services"
)

func NewRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Get("/", HomeHandler)
	r.Post("/usuarios", CreateUserHandler(db))
	r.With(JWTMiddleware).Get("/usuarios", GetUsersHandler(db))
	r.Post("/login", LoginHandler(db))

	return r
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "InkZen funcionando correctamente 🚀")
}

func CreateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.Usuario

		// Decodificar JSON
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		// Validación básica
		if user.Nombre == "" || user.Email == "" || user.Password == "" {
			http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
			return
		}

		// Hashear contraseña
		user, err = services.PrepareUser(user)
		if err != nil {
			http.Error(w, "Error procesando contraseña", http.StatusInternalServerError)
			return
		}

		// Guardar en base
		err = repository.CreateUser(db, user)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				http.Error(w, "El email ya está registrado", http.StatusBadRequest)
				return
			}

			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Usuario creado correctamente"))
	}
}

func GetUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT id, nombre, email FROM usuarios")
		if err != nil {
			http.Error(w, "Error consultando usuarios", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var usuarios []models.Usuario

		for rows.Next() {
			var u models.Usuario
			err := rows.Scan(&u.ID, &u.Nombre, &u.Email)
			if err != nil {
				http.Error(w, "Error leyendo usuarios", http.StatusInternalServerError)
				return
			}
			usuarios = append(usuarios, u)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(usuarios)
	}
}
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var credentials models.Usuario

		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		user, err := repository.GetUserByEmail(db, credentials.Email)
		if err != nil {
			http.Error(w, "Usuario no encontrado", http.StatusUnauthorized)
			return
		}

		err = services.CheckPassword(user.Password, credentials.Password)
		if err != nil {
			http.Error(w, "Contraseña incorrecta", http.StatusUnauthorized)
			return
		}

		token, err := services.GenerateJWT(user.ID)
		if err != nil {
			http.Error(w, "Error generando token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})

	}
}
