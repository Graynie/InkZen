package main

import (
	"fmt"
	"net/http"

	"github.com/Graynie/InkZen/internal/handlers"
)

func main() {
	router := handlers.NewRouter()

	fmt.Println("Servidor corriendo en http://localhost:3000")
	http.ListenAndServe(":3000", router)
}
