package handler

import (
	"Event-Management-System-Go-PSQL/storage"
	"fmt"
	"html/template"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
)

type BookingFormData struct {
	CSRFField   template.HTML
	Form        storage.Booking
	FormErrors  map[string]string
	Event       []storage.Events
	SingleEvent storage.Events
}

// Booking form with Some Event information
func (s *Server) createBooking(w http.ResponseWriter, r *http.Request) {
	log.Println("Booking : Create Method")
	data := BookingFormData{
		CSRFField: csrf.TemplateField(r),
	}
	s.loadBookingTemplate(w, r, data)
}

// Save booking with Event Id
func (s *Server) saveBooking(w http.ResponseWriter, r *http.Request) {
	ParseFormData(r)
	var form storage.Booking

	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("Decoding error")
	}
	// GET EVENT ID
	// covert int32---->int---->string
	id := form.EventId
	t := IntToStringConversion(id)
	// Search EVENT by id for matching price and seat amount
	et, err := s.store.GetDataById(t)
	// log message if data not found
	UnableToGetData(err)
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
	UnableToInsertData(err)
	http.Redirect(w, r, "/booking/boucher", http.StatusSeeOther)
}

func (s *Server) bookingBoucher(w http.ResponseWriter, r *http.Request) {
	tmpl := s.templates.Lookup("boucher.html")
	UnableToFindHtmlTemplate(tmpl)
	err := tmpl.Execute(w, nil)
	ExcutionTemplateError(err)

}

func (s *Server) loadBookingTemplate(w http.ResponseWriter, r *http.Request, form BookingFormData) {
	tmpl := s.templates.Lookup("booking-form.html")
	UnableToFindHtmlTemplate(tmpl)
	ev, err := s.store.GetEvent()
	UnableToGetData(err)
	tempData := BookingFormData{
		Form:       storage.Booking{},
		FormErrors: map[string]string{},
		Event:      ev,
	}

	err = tmpl.Execute(w, tempData)
	ExcutionTemplateError(err)
}

// Booking form With Event Id
func (s *Server) createBookingByEventId(w http.ResponseWriter, r *http.Request) {
	log.Println("Booking : Create Method")
	data := BookingFormData{
		CSRFField: csrf.TemplateField(r),
	}
	id := r.FormValue("id")
	s.loadBookingTemplateByEventId(w, r, data, id)
}

// Load form with Event id
func (s *Server) loadBookingTemplateByEventId(w http.ResponseWriter, r *http.Request, form BookingFormData, id string) {
	tmpl := s.templates.Lookup("booking-by-id.html")
	UnableToFindHtmlTemplate(tmpl)
	ev, err := s.store.GetDataById(id)
	UnableToGetData(err)
	tempData := BookingFormData{
		Form:        storage.Booking{},
		FormErrors:  map[string]string{},
		SingleEvent: ev,
	}
	fmt.Println("Booking Form data  Single event = ", ev)
	err = tmpl.Execute(w, tempData)
	ExcutionTemplateError(err)
}

/* // Booking form
func (s *Server) bookingEventByEventId(w http.ResponseWriter, r *http.Request) {
	log.Println("Method : booking event called")
	tempData := BookingFormData{
		CSRFField: csrf.TemplateField(r),
	}
	tmpl := s.templates.Lookup("booking-by-id.html")
	if tmpl == nil {
		log.Println("Unable to find Booking form")
	}

	id := r.FormValue("id")
	if id == "" {
		log.Println("Not found.")
	}

	et, err := s.store.GetDataById(id)

	tempData = BookingFormData{
		Form:        storage.Booking{},
		FormErrors:  map[string]string{},
		SingleEvent: et,
	}
	err = tmpl.Execute(w, tempData)

	if err != nil {
		log.Println("Error executing template:", err)
	}
} */

/* // Booking save
func (s *Server) saveBookingByEventId(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalln("parsing Error")
	}

	var form storage.Booking
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("Booking Page Decoding Error")
	}

	id := r.FormValue("id")
	if id == "" {
		log.Println("Not found.")
	}

	fmt.Println("133 number line = ", id)
	et, err := s.store.GetDataById(id)
	fmt.Println(et)

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
} */
