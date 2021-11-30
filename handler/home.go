package handler

import (
	"net/http"
)

func (h *Handler) Home(rw http.ResponseWriter, r *http.Request) {

		queryFilter := r.URL.Query().Get("query")

		books := []BookData{}
		

		nameQuery := `SELECT * FROM books WHERE name ILIKE '%%' || $1 || '%%' order by id desc`
		if err := h.db.Select(&books, nameQuery, queryFilter); err != nil {
			
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


		if err := h.templates.ExecuteTemplate(rw, "index.html", lt); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}


}
