package main

import (
	"log"
	"net/http"

	"example.com/m/handler"
	"github.com/gorilla/schema"
	
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
		email_verified  text,
		status boolean,
		forgot text,

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

	store := sessions.NewCookieStore([]byte("19C451AEDFB24C338A9A2A5C31D66051"))
	r := handler.New(db, decoder, store)

	log.Println("Server Starting....")
	if err := http.ListenAndServe("127.0.0.1:3000", r); err != nil {
		log.Fatal(err)
	}
}
