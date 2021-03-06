package handler

import (
	"fmt"
	"math"
	"net/http"

	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type FormData struct {
	ID     int    `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Errors map[string]string
}

func (c *FormData) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required.Error("This Filed cannot be blank"), validation.Length(3, 0)),
	)
}

type ListData struct {
	Category []FormData
	QueryFilter string
	Pagination  []Pagination
	CurrentPage int
	PrePageURL  string
	NextPageURL string
}
type Pagination struct{
	URL string
	PageNo int
}

// Show
func (h *Handler) categoryList(rw http.ResponseWriter, r *http.Request) {

	var p int
	var err error
	var nextPageURL string
	var prePageURL string

	page := r.URL.Query().Get("page")
	
	if page != "" {
		p, err = strconv.Atoi(page)
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	offset := 0
	limit := 3

	if p > 0 {
		offset = limit * p - limit
	}
	
	total := 0

	//h.db.Get(&total, "SELECT count(*) FROM categories")
		// queryFilter := r.URL.Query().Get("query")
	category := []FormData{}

	h.db.Get(&total, "SELECT count(*) FROM category")
    h.db.Select(&category, "SELECT * FROM category OFFSET $1 LIMIT $2", offset, limit)

	// nameQuery := `SELECT * FROM category WHERE name ILIKE '%%' || $1 || '%%' order by id desc`
	// if err := h.db.Select(&category, nameQuery, queryFilter,); err != nil {
	// 	http.Error(rw, err.Error(), http.StatusInternalServerError)
	// 	return
	// }


	totalPage := int(math.Ceil(float64(total)/float64(limit)))
	pagination := make([]Pagination,totalPage)
	for i:=0; i<totalPage; i++ {
		pagination[i] = Pagination{
			URL: fmt.Sprintf("http://localhost:3000/Category/List?page=%d", i +1),
			PageNo: i +1,
		}
		if i + 1 == p {
			if i != 0 {
				prePageURL = fmt.Sprintf("http://localhost:3000/Category/List?page=%d", i)
			}else{
				nextPageURL = fmt.Sprintf("http://localhost:3000/Category/List?page=%d", i +2)
			}
		}
	}

	lt := ListData{
		Category: category,
		Pagination: pagination,
		CurrentPage: p,
		PrePageURL: prePageURL,
		NextPageURL: nextPageURL,
	}
	if err := h.templates.ExecuteTemplate(rw, "list-category.html", lt); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Add
func (h *Handler) categoryCreate(rw http.ResponseWriter, r *http.Request) {

	vErrs := map[string]string{"name": ""}
	name := ""
	h.createFormData(rw, name, vErrs)
	return

}

func (h *Handler) createFormData(rw http.ResponseWriter, name string, errs map[string]string) {
	form := FormData{
		Name:   name,
		Errors: errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "create-category.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

//Store
func (h *Handler) categoryStore(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	var catfild FormData
	if err := h.decoder.Decode(&catfild, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if aErr := catfild.Validate(); aErr != nil {
		//fmt.Printf("%T", aErr)
		vErrors, ok := aErr.(validation.Errors)
		if ok {
			vErr := make(map[string]string)
			for key, value := range vErrors {
				vErr[key] = value.Error()
			}
			h.createFormData(rw, catfild.Name, vErr)
			return
		}

		http.Error(rw, aErr.Error(), http.StatusInternalServerError)
		return
	}

	const insertcategory = `INSERT INTO category(name) VALUES ($1)`
	res := h.db.MustExec(insertcategory, catfild.Name)
	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/Category/List", http.StatusTemporaryRedirect)
}

//Edit
func (h *Handler) categoryEdit(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]

	if Id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const getCat = `SELECT * FROM category WHERE id=$1`
	var category FormData
	h.db.Get(&category, getCat, Id)

	errs := map[string]string{}
	h.editFormData(rw, category.ID, category.Name, errs)
	return
}

func (h *Handler) editFormData(rw http.ResponseWriter, id int, name string, errs map[string]string) {
	form := FormData{
		ID:     id,
		Name:   name,
		Errors: errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "edit-category.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

//Update
func (h *Handler) categoryUpdate(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]

	var catfild FormData

	if err := h.decoder.Decode(&catfild, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(Id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if aErr := catfild.Validate(); aErr != nil {
		vErrors, ok := aErr.(validation.Errors)
		if ok {
			vErr := make(map[string]string)
			for key, value := range vErrors {
				vErr[key] = value.Error()
			}
			h.createFormData(rw, catfild.Name, vErr)
			return
		}

		http.Error(rw, aErr.Error(), http.StatusInternalServerError)
		return
	}

	const getCategory = `SELECT * FROM category WHERE id=$1`
	var cat FormData
	h.db.Get(&cat, getCategory, Id)

	if id == 0 {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}
	const updateStatusCategory = `UPDATE category SET name=$1 WHERE id=$2`
	res := h.db.MustExec(updateStatusCategory, catfild.Name, Id)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/Category/List", http.StatusTemporaryRedirect)
}

//Delete
func (h *Handler) categoryDelete(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]

	if Id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const deleteCategory = `DELETE FROM category WHERE id=$1`
	res := h.db.MustExec(deleteCategory, Id)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/Category/List", http.StatusTemporaryRedirect)
}