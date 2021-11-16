package handler

import (
	"github.com/jmoiron/sqlx"
	"text/template"
)

type Handler struct {
	templates *template.Template
	db        *sqlx.DB
}

func New(db *sqlx.DB) *Handler {
	h := &Handler{
		db: db,
	}

	h.parseTemplates()

	return h
}

func (h *Handler) parseTemplates() {
	h.templates = template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/create-category.html",
		"templates/list-category.html",
		"templates/edit-category.html",
	))
}
