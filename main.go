package main

import (
	"example.com/m/handler"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	var schema = `
	CREATE TABLE IF NOT EXISTS category (
		id serial,
		name text,
		
		primary key(id)
	);
`
	db, err := sqlx.Connect("postgres", "user=postgres password=dipto dbname=library sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)

	h := handler.New(db)

	http.HandleFunc("/", h.Home)
	http.HandleFunc("/Category/List", h.CategoryList)
	http.HandleFunc("/Category/create", h.CategoryCreate)
	http.HandleFunc("/Category/store", h.CategoryStore)
	http.HandleFunc("/Category/edit/", h.CategoryEdit)
	http.HandleFunc("/Category/update/", h.CategoryUpdate)
	http.HandleFunc("/Category/delete/", h.CategoryDelete)

	log.Println("Server Starting....")
	if err := http.ListenAndServe("127.0.0.1:3000", nil); err != nil {
		log.Fatal(err)
	}
}
