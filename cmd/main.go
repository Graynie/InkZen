package main

import (
	"fmt"
	"net/http"

	"github.com/Graynie/InkZen/internal/handlers"
	"github.com/Graynie/InkZen/internal/repository"
)

func main() {
	db := repository.NewDatabase()
	defer db.Close()

	repository.InitSchema(db)

	router := handlers.NewRouter(db)

	fmt.Println("Servidor corriendo en http://localhost:3000")
	http.ListenAndServe(":3000", router)
}
