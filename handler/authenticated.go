package handler

import (
	// "crypto/tls"
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	
	"net/http"
	"net/smtp"
	"text/template"


	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

type RegistrationData struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Email_verified string `db:"email_verified"`
	Status bool `db:"status,"`
	Errors   map[string]string
}
type loginData struct {
	Email    string
	Password string
	Errors   map[string]string
}

type Forgot struct {
	ID       int
	Name     string
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

func (f *Forgot) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Email, validation.Required.Error("This Filed cannot be blank")),
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
	
	var Status= false
	

	const insertUser = `INSERT INTO users(name,email,password,email_verified,status) VALUES ($1,$2,$3,$4,$5)`

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	s := user.Email
  
    se := base64.StdEncoding.EncodeToString([]byte(s))
	

	email_verified := se
	//fmt.Println(email_verified)


	res := h.db.MustExec(insertUser, user.Name, user.Email, string(pass),email_verified,Status)

	// Sender data.
	from := "sudipto397@gmail.com"
	password := "Su_Dipto@397"

	// Receiver email address.
	to := []string{
		user.Email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("templates/template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))


	var linkadd = fmt.Sprintf("http://localhost:3000/verified?token=%s", email_verified)

	t.Execute(&body, struct {
		Name    string
		Link string
	}{
		Name:    user.Name,
		Link:  linkadd,
	})

	// Sending email.
	qerr := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if qerr != nil {
		fmt.Println(qerr)
		return
	}
	fmt.Println("Email Sent!")
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
	var active = true
	//	fmt.Println(hashedPassword)
	const getuser = `SELECT * FROM users WHERE email=$1 WHERE status = $2 `
	var loginuser RegistrationData
	aerr := h.db.Get(&loginuser, getuser, usermail,active)
	if loginuser.ID == 0 {
		vErrs := map[string]string{"email": "", "password": ""}

		email := ""
		password := ""
		h.loginFormData(rw, email, password, vErrs)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(loginuser.Password), []byte(usera.Password)); err != nil {
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
	if err != nil {
		log.Fatal(err)
	}
	session.Values["authenticated"] = true
	session.Values["username"] = loginuser.ID

	if err := session.Save(r, rw); err != nil {
		log.Fatal(err)
	}
	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) userLogout(rw http.ResponseWriter, r *http.Request) {

	session, _ := h.sess.Get(r, sessionsName)
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	session.Save(r, rw)

	http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
}

func (h *Handler) userForgot(rw http.ResponseWriter, r *http.Request) {

	vErrs := map[string]string{"name": "", "email": "", "password": ""}

	email := ""

	h.forgotFormData(rw, email, vErrs)
	return
}

func (h *Handler) forgotFormData(rw http.ResponseWriter, email string, errs map[string]string) {
	form := Forgot{
		Email:  email,
		Errors: errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "auth-recoverpw.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) userForgotCheck(rw http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	var user Forgot
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
			h.forgotFormData(rw, user.Email, vErr)
			return
		}
		http.Error(rw, aErr.Error(), http.StatusInternalServerError)
		return
	}
	var usermail = user.Email
	//	fmt.Println(usermail)

	//const getuser = ` `
	var loginuser RegistrationData
	h.db.Get(&loginuser, "SELECT * FROM users WHERE email=$1", usermail)

	///	fmt.Println(loginuser)

	// Sender data.
	from := "sudipto397@gmail.com"
	password := "Su_Dipto@397"

	// Receiver email address.
	to := []string{
		usermail,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("templates/template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "All Ok",
		Message: "This is a test message in a HTML template",
	})

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")

	http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
}
func (h *Handler) userVerified(rw http.ResponseWriter, r *http.Request) {

	token := r.URL.Query().Get("token")
	
	fmt.Println(token)


	const getuser = `SELECT * FROM users WHERE email_verified=$1`
	var loginuser RegistrationData
	h.db.Get(&loginuser, getuser, token)



	fmt.Println(loginuser.ID)
	
	const updateStatusUser = `UPDATE users SET status = true WHERE id=$1`
	res := h.db.MustExec(updateStatusUser, loginuser.ID)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	

	http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
}