package handler

import (
	"net/http"
	"strconv"
)

type FormData struct {
	ID     int    `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Errors map[string]string
}

type ListData struct {
	Category []FormData
}

// Show
func (h *Handler) CategoryList(rw http.ResponseWriter, r *http.Request) {

	category := []FormData{}
	h.db.Select(&category, "SELECT * FROM category")
	lt := ListData{
		Category: category,
	}
	if err := h.templates.ExecuteTemplate(rw, "list-category.html", lt); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Add
func (h *Handler) CategoryCreate(rw http.ResponseWriter, r *http.Request) {

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
func (h *Handler) CategoryStore(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	name := r.FormValue("Name")
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
	const insertTodo = `INSERT INTO category(name) VALUES ($1)`
	res := h.db.MustExec(insertTodo, name)
	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/Category/List", http.StatusTemporaryRedirect)
}

//Edit
func (h *Handler) CategoryEdit(rw http.ResponseWriter, r *http.Request) {
	Id := r.URL.Path[len("/Category/edit/"):]

	if Id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const getCat = `SELECT * FROM category WHERE id=$1`
	var category FormData
	h.db.Get(&category, getCat, Id)

	errs := map[string]string{"name": "This field is required"}
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
func (h *Handler) CategoryUpdate(rw http.ResponseWriter, r *http.Request) {
	Id := r.URL.Path[len("/Category/update/"):]

	if Id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	nam := r.FormValue("Name")
	id, err := strconv.Atoi(Id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if nam == "" {
		errs := map[string]string{"name": "This field is required"}
		h.editFormData(rw, id, nam, errs)
		return
	}

	if len(nam) < 3 {
		errs := map[string]string{"name": "This field must be greater than or equals 3"}
		h.editFormData(rw, id, nam, errs)
		return
	}
	const updateStatusCategory = `UPDATE category SET name=$1 WHERE id=$2`
	res := h.db.MustExec(updateStatusCategory, nam, id)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/Category/List", http.StatusTemporaryRedirect)
}

//Delete
func (h *Handler) CategoryDelete(rw http.ResponseWriter, r *http.Request) {
	Id := r.URL.Path[len("/Category/delete/"):]

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
