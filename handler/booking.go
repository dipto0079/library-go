package handler

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
	"net/http"
)


type Booking struct {
	ID       int    `db:"id" json:"id"`
	User_id   int    `db:"user_id" json:"user_id"`
	Book_id   int    `db:"book_id" json:"book_id"`
	Start_time   string    `db:"start_time" json:"start_time"`
	End_time   string    `db:"end_time" json:"end_time"`
	Errors   map[string]string
}

type BookingListData struct {
	Bookings []Booking
}
type BookingFileData struct {
	Book_id int
	User_id int
	Start_time string
	End_time string
	Errors   map[string]string
}

func (b *BookingFileData) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.User_id, validation.Required.Error("This Filed cannot be blank")),
		validation.Field(&b.Book_id, validation.Required.Error("This Filed cannot be blank")),
		validation.Field(&b.Start_time, validation.Required.Error("This Filed cannot be blank")),
		validation.Field(&b.End_time, validation.Required.Error("This Filed cannot be blank")),
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

	Book_id:= book.ID
	User_id:=1
	Start_time:=""
	End_time:=""
	errs := map[string]string{}
	h.bookingFormData(rw,Book_id,User_id, Start_time,End_time, errs)
	return

}

func (h *Handler) bookingFormData(rw http.ResponseWriter,book_id int, user_id int,start_time string,end_time string, errs map[string]string) {

	form := BookingFileData{
		Book_id:   book_id,
		User_id:     user_id,
		Start_time:     start_time,
		End_time:     end_time,
		Errors:   errs,
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
			h.bookingFormData(rw, booking.Book_id,booking.User_id,booking.Start_time,booking.End_time, vErrs)
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
	res := h.db.MustExec(insertBook, booking.Book_id, booking.User_id, booking.Start_time,booking.End_time)
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
	var booking Booking
	const getBooking = `SELECT * FROM bookings WHERE book_id=$1`
	h.db.Get(&booking, getBooking, Id)

	fmt.Println(booking)

	if err := h.templates.ExecuteTemplate(rw, "single-book.html", booking); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	//for key, value := range booking {
	//	const getCat = `SELECT name FROM books WHERE id=$1`
	//	var books BookListData
	//	h.db.Get(&books, getCat, value.book_id)
	//	booking[key].Cat_Name = category.Name
	//}
	//
	//lt := BookListData{
	//	Book: books,
	//}
	//
	//if err := h.templates.ExecuteTemplate(rw, "list-book.html", lt); err != nil {
	//	http.Error(rw, err.Error(), http.StatusInternalServerError)
	//	return
	//}
}

