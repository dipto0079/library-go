package handler

import (
	// "fmt"
	"net/http"
)

func (h *Handler) Home(rw http.ResponseWriter, r *http.Request) {

	session, _ := cookie.Get(r, "Golang-session")
	var authenticated interface{} = session.Values["authenticated"]
	var authUsermail  = session.Values["username"]
	//fmt.Println(authUsermail)
	if authenticated != nil {
		isAuthenticated := session.Values["authenticated"].(bool)
		if !isAuthenticated {
			http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)

			//fmt.Printf("not")
			//http.Redirect(rw, r, "/login", http.StatusFound)
			return
		}
		//	fmt.Printf("Ok")

		queryFilter := r.URL.Query().Get("query")

		books := []BookData{}
		//h.db.Select(&book, "SELECT * from books INNER JOIN category on books.cat_id = category.id")

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
			UserEmail:   authUsermail,
			Category: categorya,
		}
		// fmt.Println(lt)

		if err := h.templates.ExecuteTemplate(rw, "index.html", lt); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
