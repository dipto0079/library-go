package handler

import (
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type BookData struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Cat_id   int    `db:"cat_id" json:"cat_id"`
	Status   bool   `db:"status"`
	Cat_Name string
	Category []FormData
	Errors   map[string]string
}

type BookListData struct {
	Book        []BookData
	QueryFilter string
	Category []FormData
}

func (b *BookData) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.Name, validation.Required.Error("This Filed cannot be blank"), validation.Length(3, 0)),
		validation.Field(&b.Cat_id, validation.Required.Error("This Filed cannot be blank")),
	)
}

// Show
func (h *Handler) bookList(rw http.ResponseWriter, r *http.Request) {
	queryFilter := r.URL.Query().Get("query")

	books := []BookData{}

	nameQuery := `SELECT * FROM books WHERE name ILIKE '%%' || $1 || '%%' order by id desc`
	if err := h.db.Select(&books, nameQuery, queryFilter); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	for key, value := range books {
		const getCat = `SELECT name FROM category WHERE id=$1`
		var category FormData
		h.db.Get(&category, getCat, value.Cat_id)
		books[key].Cat_Name = category.Name
	}

	categorya := []FormData{}

		namezQuery := `SELECT * FROM category  order by id desc`

		if err := h.db.Select(&categorya,namezQuery ); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
		}

	lt := BookListData{
		Book:        books,
		QueryFilter: queryFilter,
		Category: categorya,
	}

	if err := h.templates.ExecuteTemplate(rw, "list-book.html", lt); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Add
func (h *Handler) bookCreate(rw http.ResponseWriter, r *http.Request) {

	vErrs := map[string]string{"name": "", "cat_id": ""}
	name := ""
	cat_id := 0
	h.createBookFormData(rw, name, cat_id, vErrs)
	return

}

func (h *Handler) createBookFormData(rw http.ResponseWriter, name string, cat_id int, errs map[string]string) {

	category := []FormData{}
	h.db.Select(&category, "SELECT * FROM category")
	form := BookData{
		Name:     name,
		Cat_id:   cat_id,
		Errors:   errs,
		Category: category,
	}
	if err := h.templates.ExecuteTemplate(rw, "create-book.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

//Store
func (h *Handler) bookStore(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	var book BookData
	if err := h.decoder.Decode(&book, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if aErr := book.Validate(); aErr != nil {
		//fmt.Printf("%T", aErr)
		vErrors, ok := aErr.(validation.Errors)
		if ok {
			vErr := make(map[string]string)
			for key, value := range vErrors {
				vErr[key] = value.Error()
			}
			h.createFormData(rw, book.Name, vErr)
			return
		}

		http.Error(rw, aErr.Error(), http.StatusInternalServerError)
		return
	}

	const insertBook = `INSERT INTO books(name,cat_id,status) VALUES ($1,$2,$3)`
	res := h.db.MustExec(insertBook, book.Name, book.Cat_id, true)
	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/Book/List", http.StatusTemporaryRedirect)
}

//Edit
func (h *Handler) bookEdit(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]

	if Id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const getBook = `SELECT * FROM books WHERE id=$1`
	var book BookData
	h.db.Get(&book, getBook, Id)

	errs := map[string]string{}
	h.editBookData(rw, book.ID, book.Name, book.Cat_id, errs)
	return
}

func (h *Handler) editBookData(rw http.ResponseWriter, id int, name string, cat_id int, errs map[string]string) {
	category := []FormData{}
	h.db.Select(&category, "SELECT * FROM category")
	form := BookData{
		ID:       id,
		Name:     name,
		Cat_id:   cat_id,
		Errors:   errs,
		Category: category,
	}
	if err := h.templates.ExecuteTemplate(rw, "edit-book.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

//Update
func (h *Handler) bookUpdate(rw http.ResponseWriter, r *http.Request) {
	categories := []FormData{}
	h.db.Select(&categories, "SELECT * FROM categories")

	vars := mux.Vars(r)
	Id := vars["id"]

	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var book BookData
	if err := h.decoder.Decode(&book, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(Id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := book.Validate(); err != nil {
		valError, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range valError {
				vErrs[key] = value.Error()
			}
			h.editBookData(rw, id, book.Name, book.Cat_id, vErrs)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	const getBook = `SELECT * FROM books WHERE id=$1`
	var books BookData
	h.db.Get(&books, getBook, Id)

	//fmt.Println(books)

	if books.ID == 0 {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const updateStatusCategory = `UPDATE books SET name=$1, cat_id=$2,status=$3 WHERE id=$4`
	res := h.db.MustExec(updateStatusCategory, book.Name, book.Cat_id, true, Id)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/Book/List", http.StatusTemporaryRedirect)
}

//Delete
func (h *Handler) bookdelete(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]

	const deleteCategory = `DELETE FROM books WHERE id=$1`
	res := h.db.MustExec(deleteCategory, Id)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/Book/List", http.StatusTemporaryRedirect)
}

func (h *Handler) bookActive(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]

	const updateStatusTodo = `UPDATE books SET status = true WHERE id=$1`
	res := h.db.MustExec(updateStatusTodo, Id)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/Book/List", http.StatusTemporaryRedirect)
}

func (h *Handler) bookDeactivate(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]

	const updateStatusTodo = `UPDATE books SET status = false WHERE id=$1`
	res := h.db.MustExec(updateStatusTodo, Id)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/Book/List", http.StatusTemporaryRedirect)
}