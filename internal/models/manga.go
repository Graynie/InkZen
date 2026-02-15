package models

type Manga struct {
	ID           int
	Titulo       string
	Autor        string
	Genero       string
	Idioma       string
	Editorial    string
	Descripcion  string
	CapitulosTot int
	Disponible   bool
}
