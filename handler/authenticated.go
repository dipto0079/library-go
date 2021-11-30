package handler

import (
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"

	"golang.org/x/crypto/bcrypt"
)

type RegistrationData struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Errors   map[string]string
}
type loginData struct {
	Email    string
	Password string
	Errors   map[string]string
}

func (L *loginData) Validate() error {
	return validation.ValidateStruct(L,
		validation.Field(&L.Email, validation.Required.Error("This Filed cannot be blank")),
		validation.Field(&L.Password, validation.Required.Error("This Filed cannot be blank")),
	)
}

func (R *RegistrationData) Validate() error {
	return validation.ValidateStruct(R,
		validation.Field(&R.Name, validation.Required.Error("This Filed cannot be blank"), validation.Length(3, 0)),
		validation.Field(&R.Email, validation.Required.Error("This Filed cannot be blank")),
		validation.Field(&R.Password, validation.Required.Error("This Filed cannot be blank")),
	)
}

// Add
func (h *Handler) registrationCreate(rw http.ResponseWriter, r *http.Request) {

	vErrs := map[string]string{"name": "", "email": "", "password": ""}
	name := ""
	email := ""
	password := ""
	h.registrationFormData(rw, name, email, password, vErrs)
	return

}

func (h *Handler) registrationFormData(rw http.ResponseWriter, name string, email string, password string, errs map[string]string) {
	form := RegistrationData{
		Name:     name,
		Email:    email,
		Password: password,
		Errors:   errs,
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
			h.registrationFormData(rw, user.Name, user.Email, user.Password, vErr)
			return
		}
		http.Error(rw, aErr.Error(), http.StatusInternalServerError)
		return
	}

	const insertUser = `INSERT INTO users(name,email,password) VALUES ($1,$2,$3)`

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	res := h.db.MustExec(insertUser, user.Name, user.Email, string(pass))
	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
}

func (h *Handler) login(rw http.ResponseWriter, r *http.Request) {

	vErrs := map[string]string{"email": "", "password": ""}

	email := ""
	password := ""
	h.loginFormData(rw, email, password, vErrs)
	return

}

func (h *Handler) loginFormData(rw http.ResponseWriter, email string, password string, errs map[string]string) {
	form := loginData{
		Email:    email,
		Password: password,
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
	var usera loginData
	if err := h.decoder.Decode(&usera, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if aErr := usera.Validate(); aErr != nil {
		vErrors, ok := aErr.(validation.Errors)
		if ok {
			vErr := make(map[string]string)
			for key, value := range vErrors {
				vErr[key] = value.Error()
			}
			h.loginFormData(rw, usera.Email, usera.Password, vErr)
			return
		}
		http.Error(rw, aErr.Error(), http.StatusInternalServerError)
		return
	}
	var usermail = usera.Email

	//	fmt.Println(hashedPassword)
	const getuser = `SELECT * FROM users WHERE email=$1 `
	var loginuser RegistrationData
	aerr := h.db.Get(&loginuser, getuser, usermail)
	if loginuser.ID == 0 {
		vErrs := map[string]string{"email": "", "password": ""}

		email := ""
		password := ""
		h.loginFormData(rw, email, password, vErrs)
	}

	if err:= bcrypt.CompareHashAndPassword([]byte(loginuser.Password),[]byte(usera.Password));err !=nil{
		vErrs := map[string]string{"email": "", "password": ""}

		email := ""
		password := ""
		h.loginFormData(rw, email, password, vErrs)
	}

	if aerr != nil {
		http.Error(rw, aerr.Error(), http.StatusInternalServerError)
		return
	}
	 session, err := h.sess.Get(r, sessionsName)
	 if err!=nil{
		 log.Fatal(err)
	 }
		 session.Values["authenticated"] = true
		 session.Values["username"] = loginuser.ID
		
		if err:= session.Save(r, rw);err!=nil{
		log.Fatal(err)
		}
	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) userLogout(rw http.ResponseWriter, r *http.Request) {

	session, _ := h.sess.Get(r, sessionsName)
	session.Values["authenticated"] = false
	session.Save(r, rw)
	
	http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
}
