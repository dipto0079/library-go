package main

import (
	"example.com/m/handler"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	var createTable = `
CREATE TABLE IF NOT EXISTS category (
		id serial,
		name text,
		
		primary key(id)
	);
CREATE TABLE IF NOT EXISTS books (
		id serial,
		name text,
		cat_id integer,
		status boolean,
		
		primary key(id)
	);
CREATE TABLE IF NOT EXISTS bookings (
		id serial,
		book_id integer,
		user_id integer,
		start_time  timestamp,
		end_time timestamp,
		primary key(id)
	);
CREATE TABLE IF NOT EXISTS users (
		id serial,
		name text,
		email text,
		password  text,

		primary key(id)
	);
`
	db, err := sqlx.Connect("postgres", "user=postgres password=dipto dbname=library sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(createTable)

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	r := handler.New(db, decoder)

	log.Println("Server Starting....")
	if err := http.ListenAndServe("127.0.0.1:3000", r); err != nil {
		log.Fatal(err)
	}
}
