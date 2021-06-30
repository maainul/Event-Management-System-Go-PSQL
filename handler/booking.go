package handler

import (
	"Event-Management-System-Go-PSQL/storage"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
)

type BookingFormData struct {
	CSRFField  template.HTML
	Form       storage.Booking
	FormErrors map[string]string
	Event      []storage.Events
}

func (s *Server) createBooking(w http.ResponseWriter, r *http.Request) {
	log.Println("Booking : Create Method")
	data := BookingFormData{
		CSRFField: csrf.TemplateField(r),
	}
	s.loadBookingTemplate(w, r, data)
}

func (s *Server) saveBooking(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalln("parsing Error")
	}

	var form storage.Booking
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("Booking Page Decoding Error")
	}
	// GET EVENT ID
	// covert int32---->int---->string
	t := strconv.Itoa(int(form.EventId))

	// Search EVENT by id for matching price and seat amount
	et, err := s.store.GetDataById(t)
	if err != nil {
		log.Println("Unable to find data")
	}

	fmt.Println("FUll OBJECT - ", et)
	fmt.Println("Ticket Remaining = ", et.TicketRemaining)
	fmt.Println("id = ", t)
	fmt.Println("guest = ", et.NumberOfGuest)
	fmt.Println("price = ", et.PerPersonPrice)

	// total ticket price = ticket price *number of ticket
	form.TotalAmount = form.NumberOfTicket * et.PerPersonPrice
	form.UserId = 1

	// decrement value as user's input
	_, err = s.store.DecrementRemainingTicketById(form.EventId, form.NumberOfTicket)

	// validation
	if err := form.Validate(); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}

		data := BookingFormData{
			CSRFField:  csrf.TemplateField(r),
			Form:       form,
			FormErrors: vErrs,
		}
		s.loadBookingTemplate(w, r, data)
		return
	}
	fmt.Println("83 line pass")
	_, err = s.store.CreateBooking(form)
	fmt.Println("85 line pass")
	if err != nil {
		log.Fatalln("Unable to find Booking form ", err)
	}

	http.Redirect(w, r, "/booking/boucher", http.StatusSeeOther)
}

func (s *Server) bookingBoucher(w http.ResponseWriter, r *http.Request) {
	tmpl := s.templates.Lookup("booking-boucher.html")
	if tmpl == nil {
		log.Println("Unable to find form")
		return
	}

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

}

func (s *Server) loadBookingTemplate(w http.ResponseWriter, r *http.Request, form BookingFormData) {
	tmpl := s.templates.Lookup("booking-form.html")
	if tmpl == nil {
		log.Println("Unable to find form")
		return
	}

	ev, err := s.store.GetEvent()
	if err != nil {
		log.Println("Unable to find Any Event ")
	}
	tempData := BookingFormData{
		Form:       storage.Booking{},
		FormErrors: map[string]string{},
		Event:      ev,
	}

	if err := tmpl.Execute(w, tempData); err != nil {
		log.Println("Error executing tempalte : ", err)
		return
	}
}
