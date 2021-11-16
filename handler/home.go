package handler

import (
	"net/http"
)

func (h *Handler) Home(rw http.ResponseWriter, r *http.Request) {
	if err := h.templates.ExecuteTemplate(rw, "index.html", nil); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
