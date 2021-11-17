package handler

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	"net/http"
	"text/template"
)

type Handler struct {
	templates *template.Template
	db        *sqlx.DB
	decoder   *schema.Decoder
}

func New(db *sqlx.DB, decoder *schema.Decoder) *mux.Router {
	h := &Handler{
		db:      db,
		decoder: decoder,
	}

	h.parseTemplates()
	r := mux.NewRouter()
	r.HandleFunc("/", h.Home)
	//Category
	r.HandleFunc("/Category/List", h.categoryList)
	r.HandleFunc("/Category/create", h.categoryCreate)
	r.HandleFunc("/Category/store", h.categoryStore)
	r.HandleFunc("/Category/{id:[0-9]+}/edit", h.categoryEdit)
	r.HandleFunc("/Category/{id:[0-9]+}/update", h.categoryUpdate)
	r.HandleFunc("/Category/{id:[0-9]+}/delete", h.categoryDelete)
	//Book
	r.HandleFunc("/Book/List", h.bookList)
	r.HandleFunc("/Book/Create", h.bookCreate)
	r.HandleFunc("/Book/store", h.bookStore)
	r.HandleFunc("/Book/{id:[0-9]+}/edit", h.bookEdit)
	r.HandleFunc("/Book/{id:[0-9]+}/update", h.bookUpdate)
	r.HandleFunc("/Book/{id:[0-9]+}/delete", h.bookdelete)
	r.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := h.templates.ExecuteTemplate(rw, "404.html", nil); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	return r
}

func (h *Handler) parseTemplates() {
	h.templates = template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/create-category.html",
		"templates/list-category.html",
		"templates/edit-category.html",
		"templates/create-book.html",
		"templates/list-book.html",
		"templates/edit-book.html",
		"templates/404.html",
	))
}
