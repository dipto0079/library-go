package handler

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/sessions"
	"net/http"

)
type RegistrationData struct {
	ID     int    `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Email   string `db:"email" json:"email"`
	Password   string `db:"password" json:"password"`
	Errors map[string]string
}

func (R *RegistrationData) Validate() error {
	return validation.ValidateStruct(R,
		validation.Field(&R.Name, validation.Required.Error("This Filed cannot be blank"), validation.Length(3, 0)),
		validation.Field(&R.Email, validation.Required.Error("This Filed cannot be blank")),
		validation.Field(&R.Password, validation.Required.Error("This Filed cannot be blank")),
	)
}
var cookie *sessions.CookieStore

func init() {
	cookie = sessions.NewCookieStore([]byte("Golang-Blogs"))
}
// Add
func (h *Handler) registrationCreate(rw http.ResponseWriter, r *http.Request) {

	vErrs := map[string]string{"name": "","email":"","password":""}
	name := ""
	email := ""
	password := ""
	h.registrationFormData(rw, name,email,password, vErrs)
	return

}

func (h *Handler) registrationFormData(rw http.ResponseWriter, name string,email string,password string, errs map[string]string) {
	form := RegistrationData{
		Name:   name,
		Email: email,
		Password: password,
		Errors: errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "registration.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UserStore(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	var user RegistrationData
	if err := h.decoder.Decode(&user, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if aErr := user.Validate(); aErr != nil {
		vErrors, ok := aErr.(validation.Errors)
		if ok {
			vErr := make(map[string]string)
			for key, value := range vErrors {
				vErr[key] = value.Error()
			}
			h.registrationFormData(rw, user.Name,user.Email,user.Password, vErr)
			return
		}
		http.Error(rw, aErr.Error(), http.StatusInternalServerError)
		return
	}


	const insertBook = `INSERT INTO users(name,email,password) VALUES ($1,$2,$3)`
	res := h.db.MustExec(insertBook, user.Name, user.Email,user.Password)
	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) login(rw http.ResponseWriter, r *http.Request) {

	vErrs := map[string]string{"name": "","email":"","password":""}
	name := ""
	email := ""
	password := ""
	h.loginFormData(rw, name,email,password, vErrs)
	return

}

func (h *Handler) loginFormData(rw http.ResponseWriter, name string,email string,password string, errs map[string]string) {
	form := RegistrationData{
		Name:   name,
		Email: email,
		Password: password,
		Errors: errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "login.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) userLogin(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	var user RegistrationData
	if err := h.decoder.Decode(&user, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if aErr := user.Validate(); aErr != nil {
		vErrors, ok := aErr.(validation.Errors)
		if ok {
			vErr := make(map[string]string)
			for key, value := range vErrors {
				vErr[key] = value.Error()
			}
			h.registrationFormData(rw, user.Name,user.Email,user.Password, vErr)
			return
		}
		http.Error(rw, aErr.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := cookie.Get(r, "Golang-session")
	session.Values["authenticated"] = true
	//session.Save(r, w)

	var usermail = user.Email

	const getuser = `SELECT * FROM users WHERE email=$1`
	var loginuser RegistrationData
	h.db.Get(&loginuser, getuser, usermail)

	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}