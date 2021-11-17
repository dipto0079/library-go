package handler

import (
	"net/http"
)

type BookData struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Cat_id   int    `db:"cat_id" json:"cat_id"`
	Category []FormData
	Errors   map[string]string
}

type BookListData struct {
	Book []BookData
}

// Show
func (h *Handler) bookList(rw http.ResponseWriter, r *http.Request) {

	book := []BookData{}
	//h.db.Select(&book, "SELECT * from books INNER JOIN category on books.cat_id = category.id")
	h.db.Select(&book, "SELECT * FROM books")

	lt := BookListData{
		Book: book,
	}
	//fmt.Println(lt)
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

func (h *Handler) bookStore(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	name := r.FormValue("Name")
	cat_id := r.FormValue("Cat_id")
	if name == "" {
		errs := map[string]string{"name": "This field is required"}
		h.createFormData(rw, name, errs)
		return
	}
	if len(name) < 3 {
		errs := map[string]string{"name": "This field must be greater than or equals 3"}
		h.createFormData(rw, name, errs)
		return
	}

	const insertBook = `INSERT INTO books(name,cat_id,status) VALUES ($1,$2,$3)`
	res := h.db.MustExec(insertBook, name, cat_id, true)
	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}
