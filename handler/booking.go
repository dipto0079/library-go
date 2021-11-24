package handler

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Booking struct {
	ID         int       `db:"id" json:"id"`
	User_id    int       `db:"user_id" json:"user_id"`
	Book_id    int       `db:"book_id" json:"book_id"`
	Start_time time.Time `db:"start_time" json:"start_time"`
	End_time   time.Time `db:"end_time" json:"end_time"`
	Errors     map[string]string
}

type BookingListData struct {
	Bookings []Booking
}
type BookingFileData struct {
	Book_id    int
	User_id    int
	Start_time time.Time
	End_time   time.Time
	Errors     map[string]string
	ST         string
	ET         string
}

func (b *BookingFileData) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.User_id, validation.Required.Error("This Filed cannot be blank")),
		validation.Field(&b.Book_id, validation.Required.Error("This Filed cannot be blank")),
		validation.Field(&b.ST, validation.Required.Error("This Filed cannot be blank")),
		validation.Field(&b.ET, validation.Required.Error("This Filed cannot be blank")),
	)
}

func (h *Handler) bookingCreate(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	Id := vars["id"]

	if Id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const getBook = `SELECT * FROM books WHERE id=$1`
	var book BookData
	h.db.Get(&book, getBook, Id)

	Book_id := book.ID
	User_id := 1
	Start_time := ""
	End_time := ""
	errs := map[string]string{}
	h.bookingFormData(rw, Book_id, User_id, Start_time, End_time, errs)
	return

}

func (h *Handler) bookingFormData(rw http.ResponseWriter, book_id int, user_id int, start_time string, end_time string, errs map[string]string) {

	form := BookingFileData{
		Book_id:    book_id,
		User_id:    user_id,
		ST: start_time,
		ET:   end_time,
		Errors:     errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "create-booking.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) bookingStore(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	Id := vars["id"]

	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var booking BookingFileData
	if err := h.decoder.Decode(&booking, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := booking.Validate(); err != nil {
		valError, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range valError {
				vErrs[key] = value.Error()
			}
			h.bookingFormData(rw, booking.Book_id, booking.User_id, booking.ST, booking.ET, vErrs)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	const getBook = `SELECT * FROM books WHERE id=$1`
	var books BookData
	h.db.Get(&books, getBook, Id)

	if books.ID == 0 {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}
	const updateStatusTodo = `UPDATE books SET status = false WHERE id=$1`
	ress := h.db.MustExec(updateStatusTodo, Id)
	if ok, err := ress.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	const insertBook = `INSERT INTO bookings(book_id,user_id,start_time,end_time) VALUES ($1,$2,$3,$4)`
	res := h.db.MustExec(insertBook, booking.Book_id, booking.User_id, booking.ST, booking.ET)
	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) bookiSingleList(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]
	if Id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}
	var booking BookingFileData
	const getBooking = `SELECT * FROM bookings WHERE book_id=$1`
	h.db.Get(&booking, getBooking, Id)

	start_time := booking.Start_time.Format("Mon Jan _2 2006 15:04 AM")
	end_time := booking.End_time.Format("Mon Jan _2 2006 15:04 AM")

	booking.ST = start_time
	booking.ET = end_time

	if err := h.templates.ExecuteTemplate(rw, "single-booking.html", booking); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}