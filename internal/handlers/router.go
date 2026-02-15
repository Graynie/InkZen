package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"html/template"

	"github.com/Graynie/InkZen/internal/models"
	"github.com/Graynie/InkZen/internal/repository"
	"github.com/Graynie/InkZen/internal/services"
	"github.com/go-chi/chi/v5"
)

func NewRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	userService := services.NewUsuarioService()

	r.Get("/", HomeHandler)
	r.Post("/usuarios", CreateUserHandler(db, userService))
	r.With(JWTMiddleware).Get("/usuarios", GetUsersHandler(db))
	r.Post("/lecturas", CreateLecturaHandler(db))
	r.Put("/lecturas", UpdateLecturaHandler(db))
	r.Get("/mis-mangas", GetLecturasHandler(db))
	r.Get("/mangas-web", WebMangasHandler(db))
	r.Get("/mangas/new", CreateMangaFormHandler())
	r.Post("/mangas-web", CreateMangaWebHandler(db))
	r.Get("/manga", ViewMangaHandler(db))
	r.Get("/capitulo", ViewCapituloHandler(db))
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	r.Get("/register", RegisterFormHandler())
	r.Post("/register", RegisterHandler(db))

	r.Get("/login", LoginFormHandler())
	r.Post("/login", LoginHandler(db))

	r.Get("/logout", LogoutHandler())

	return r
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "InkZen funcionando correctamente 🚀")
}
func CreateUserHandler(db *sql.DB, userService *services.UsuarioService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.Usuario

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if user.Nombre == "" || user.Email == "" || user.Password == "" {
			http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
			return
		}

		user, err = userService.PrepareUser(user)
		if err != nil {
			http.Error(w, "Error procesando contraseña", http.StatusInternalServerError)
			return
		}

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
func CreateLecturaHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var lectura models.Lectura

		err := json.NewDecoder(r.Body).Decode(&lectura)
		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		err = repository.CrearLectura(db, lectura)
		if err != nil {
			http.Error(w, "Error creando lectura", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Lectura creada"))
	}
}
func UpdateLecturaHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var lectura models.Lectura

		err := json.NewDecoder(r.Body).Decode(&lectura)
		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		err = repository.ActualizarCapitulo(db, lectura.UsuarioID, lectura.MangaID, lectura.CapituloActual)
		if err != nil {
			http.Error(w, "Error actualizando capítulo", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Capítulo actualizado"))
	}
}
func GetLecturasHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		usuarioIDStr := r.URL.Query().Get("usuario_id")
		if usuarioIDStr == "" {
			http.Error(w, "usuario_id requerido", http.StatusBadRequest)
			return
		}

		usuarioID, _ := strconv.Atoi(usuarioIDStr)

		lecturas, err := repository.ObtenerLecturasUsuario(db, usuarioID)
		if err != nil {
			http.Error(w, "Error obteniendo lecturas", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(lecturas)
	}
}

func WebMangasHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 🔹 Obtener usuario desde cookie
		var nombreUsuario string

		usuarioID, errUser := getUserIDFromRequest(r)
		if errUser == nil && usuarioID > 0 {

			err := db.QueryRow(
				"SELECT nombre FROM usuarios WHERE id = ?",
				usuarioID,
			).Scan(&nombreUsuario)

			if err != nil {
				nombreUsuario = ""
			}
		}

		// 🔹 Obtener búsqueda
		query := r.URL.Query().Get("q")

		rows, err := db.Query(`
			SELECT id, titulo, autor, genero, idioma, editorial, descripcion, capitulos_tot, disponible
			FROM mangas
			WHERE disponible = 1 AND titulo LIKE ?
		`, "%"+query+"%")

		if err != nil {
			http.Error(w, "Error obteniendo mangas", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// 🔹 Declarar slice correctamente
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
				continue
			}
			mangas = append(mangas, m)
		}

		// 🔹 Enviar datos al template
		data := map[string]interface{}{
			"Mangas":  mangas,
			"Usuario": nombreUsuario,
		}

		tmpl, err := template.ParseFiles("web/templates/mangas.html")
		if err != nil {
			http.Error(w, "Error cargando template", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data)
	}
}

func CreateMangaFormHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("web/templates/create_manga.html")
		if err != nil {
			http.Error(w, "Error cargando formulario", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, nil)
	}
}
func CreateMangaWebHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error procesando formulario", http.StatusBadRequest)
			return
		}

		disponible := false
		if r.FormValue("disponible") == "true" {
			disponible = true
		}

		capitulosTot, _ := strconv.Atoi(r.FormValue("capitulos_tot"))

		manga := models.Manga{
			Titulo:       r.FormValue("titulo"),
			Autor:        r.FormValue("autor"),
			Genero:       r.FormValue("genero"),
			Idioma:       r.FormValue("idioma"),
			Editorial:    r.FormValue("editorial"),
			Descripcion:  r.FormValue("descripcion"),
			CapitulosTot: capitulosTot,
			Disponible:   disponible,
		}

		err := repository.CreateManga(db, manga)
		if err != nil {
			http.Error(w, "Error guardando manga", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/mangas-web", http.StatusSeeOther)
	}
}
func ViewMangaHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		mangaIDStr := r.URL.Query().Get("id")
		mangaID, _ := strconv.Atoi(mangaIDStr)

		// Obtener manga
		var mangas []models.Manga
		var manga models.Manga
		err := db.QueryRow(`
			SELECT id, titulo, autor, genero, idioma, editorial, descripcion, capitulos_tot, disponible
			FROM mangas
			WHERE id = ?
		`, mangaID).Scan(
			&manga.ID,
			&manga.Titulo,
			&manga.Autor,
			&manga.Genero,
			&manga.Idioma,
			&manga.Editorial,
			&manga.Descripcion,
			&manga.CapitulosTot,
			&manga.Disponible,
		)

		if err != nil {
			http.Error(w, "Manga no encontrado", http.StatusNotFound)
			return
		}

		// Leer capítulos desde carpeta
		path := fmt.Sprintf("web/static/uploads/%d/capitulos/", manga.ID)
		files, _ := os.ReadDir(path)

		type CapituloView struct {
			Numero string
			Leido  bool
			Actual bool
		}

		var capitulos []CapituloView

		// Obtener progreso
		usuarioID, errUser := getUserIDFromRequest(r)

		var nombreUsuario string

		if errUser == nil {
			db.QueryRow("SELECT nombre FROM usuarios WHERE id = ?", usuarioID).Scan(&nombreUsuario)
		}

		var capActual int

		errLectura := db.QueryRow(`
			SELECT capitulo_actual
			FROM lecturas
			WHERE usuario_id = ? AND manga_id = ?
		`, usuarioID, manga.ID).Scan(&capActual)

		if errLectura != nil {
			capActual = 0
		}
		for _, file := range files {
			if file.IsDir() {

				num, _ := strconv.Atoi(file.Name())

				c := CapituloView{
					Numero: file.Name(),
				}

				if num < capActual {
					c.Leido = true
				}
				if num == capActual {
					c.Actual = true
				}

				capitulos = append(capitulos, c)
			}
		}

		sort.Slice(capitulos, func(i, j int) bool {
			ci, _ := strconv.Atoi(capitulos[i].Numero)
			cj, _ := strconv.Atoi(capitulos[j].Numero)
			return ci < cj
		})

		var porcentaje int
		if manga.CapitulosTot > 0 {
			porcentaje = (capActual * 100) / manga.CapitulosTot
		}

		data := map[string]interface{}{
			"Manga":      manga,
			"Capitulos":  capitulos,
			"Continue":   capActual,
			"Progreso":   capActual,
			"Porcentaje": porcentaje,
			"Mangas":     mangas,
			"Usuario":    nombreUsuario,
		}

		tmpl, _ := template.ParseFiles("web/templates/manga_detalle.html")
		tmpl.Execute(w, data)
	}
}
func ViewCapituloHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		manga := r.URL.Query().Get("manga")
		cap := r.URL.Query().Get("cap")

		basePath := fmt.Sprintf("web/static/uploads/%s/capitulos/", manga)
		path := basePath + cap + "/"

		files, err := os.ReadDir(path)
		if err != nil {
			http.Error(w, "Capítulo no encontrado", http.StatusNotFound)
			return
		}

		var imagenes []string
		for _, file := range files {
			if !file.IsDir() {
				imagenes = append(imagenes,
					"/static/uploads/"+manga+"/capitulos/"+cap+"/"+file.Name())
			}
		}

		sort.Strings(imagenes)

		// 🔹 Guardar progreso automático
		usuarioID, errUser := getUserIDFromRequest(r)
		if errUser != nil {
			return
		}

		capituloNum, _ := strconv.Atoi(cap)
		mangaID, _ := strconv.Atoi(manga)

		var capActual int
		errCheck := db.QueryRow(`
			SELECT capitulo_actual 
			FROM lecturas 
			WHERE usuario_id = ? AND manga_id = ?
		`, usuarioID, mangaID).Scan(&capActual)

		if errCheck != nil {
			db.Exec(`
				INSERT INTO lecturas (usuario_id, manga_id, capitulo_actual)
				VALUES (?, ?, ?)
			`, usuarioID, mangaID, capituloNum)
		} else {
			if capituloNum > capActual {
				db.Exec(`
					UPDATE lecturas 
					SET capitulo_actual = ?
					WHERE usuario_id = ? AND manga_id = ?
				`, capituloNum, usuarioID, mangaID)
			}
		}

		data := map[string]interface{}{
			"Imagenes": imagenes,
			"Manga":    manga,
			"Capitulo": cap,
		}

		tmpl, _ := template.ParseFiles("web/templates/view_capitulo.html")
		tmpl.Execute(w, data)
	}
}
func getUserIDFromRequest(r *http.Request) (int, error) {

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		fmt.Println("No cookie")
		return 0, err
	}

	fmt.Println("Token recibido:", cookie.Value)

	userID, err := services.GetUserIDFromToken(cookie.Value)
	fmt.Println("UserID decodificado:", userID, "Error:", err)

	return userID, err
}
func RegisterFormHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("web/templates/register.html")
		tmpl.Execute(w, nil)
	}
}
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		nombre := r.FormValue("nombre")
		email := r.FormValue("email")
		password := r.FormValue("password")

		hashed, _ := services.HashPassword(password)

		_, err := db.Exec(`
			INSERT INTO usuarios (nombre, email, password)
			VALUES (?, ?, ?)
		`, nombre, email, hashed)

		if err != nil {
			http.Error(w, "Error creando usuario", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
func LoginFormHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("web/templates/login.html")
		tmpl.Execute(w, nil)
	}
}
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		email := r.FormValue("email")
		password := r.FormValue("password")

		var user models.Usuario

		err := db.QueryRow(`
			SELECT id, nombre, password 
			FROM usuarios 
			WHERE email = ?
		`, email).Scan(&user.ID, &user.Nombre, &user.Password)

		if err != nil {
			http.Error(w, "Usuario no encontrado", http.StatusUnauthorized)
			return
		}

		err = services.CheckPassword(user.Password, password)
		if err != nil {
			http.Error(w, "Contraseña incorrecta", http.StatusUnauthorized)
			return
		}

		token, _ := services.GenerateJWT(user.ID)

		http.SetCookie(w, &http.Cookie{
			Name:  "auth_token",
			Value: token,
			Path:  "/",
		})

		http.Redirect(w, r, "/mangas-web", http.StatusSeeOther)
	}
}
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		http.SetCookie(w, &http.Cookie{
			Name:   "auth_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
