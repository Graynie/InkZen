package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "InkZen funcionando correctamente 🚀")
	})

	fmt.Println("Servidor corriendo en http://localhost:3000")
	http.ListenAndServe(":3000", r)
}
