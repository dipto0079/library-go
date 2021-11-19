package handler

import (
	"net/http"
)

func (h *Handler) Home(rw http.ResponseWriter, r *http.Request) {
	books := []BookData{}
	//h.db.Select(&book, "SELECT * from books INNER JOIN category on books.cat_id = category.id")
	h.db.Select(&books, "SELECT * FROM books order by id desc")

	for key, value := range books {
		const getCat = `SELECT name FROM category WHERE id=$1`
		var category FormData
		h.db.Get(&category, getCat, value.Cat_id)
		books[key].Cat_Name = category.Name
	}

	lt := BookListData{
		Book: books,
	}

	if err := h.templates.ExecuteTemplate(rw, "index.html", lt); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (h *Handler) homeSearching(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	ser := r.FormValue("Searching")
	if ser == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const getSrc = `SELECT * FROM books WHERE name ILIKE '%%' || $1 || '%%'`
	var books []BookData
	h.db.Select(&books, getSrc,ser)

	for key, value := range books {
		const getCat = `SELECT name FROM category WHERE id=$1`
		var category FormData
		h.db.Get(&category, getCat, value.Cat_id)
		books[key].Cat_Name = category.Name
	}

	lt := BookListData{
		Book: books,
	}

	if err := h.templates.ExecuteTemplate(rw, "index.html", lt); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
