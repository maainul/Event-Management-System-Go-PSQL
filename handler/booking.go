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

/*-----------------------------------------------------------------------------GET : Booking form With Event Id/Show Booking Form----------------------------------------------*/
func (s *Server) createBooking(w http.ResponseWriter, r *http.Request) {
	log.Println("Booking : Create Method")
	data := BookingFormData{
		CSRFField: csrf.TemplateField(r),
	}
	s.loadBookingTemplate(w, r, data)
}

/*--------------------------------------------------------------------------------POST : Save Booking With DropDown of Event----------------------------------------------*/
func (s *Server) saveBooking(w http.ResponseWriter, r *http.Request) {
	ParseFormData(r)
	var form storage.Booking
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("Decoding error")
	}
	id := form.EventId
	t := IntToStringConversion(id)
	et, err := s.store.GetDataById(t)
	UnableToGetData(err)
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
	_, err = s.store.CreateBooking(form)
	fmt.Println("85 line pass")
	http.Redirect(w, r, "/booking/boucher", http.StatusSeeOther)
}

/*------------------------------------------------------------------------------Booking Information/ Boucher/Invoice of Booking---------------------------------------------*/
func (s *Server) bookingBoucher(w http.ResponseWriter, r *http.Request) {
	tmpl := s.templates.Lookup("boucher.html")
	UnableToFindHtmlTemplate(tmpl)
	err := tmpl.Execute(w, nil)
	ExcutionTemplateError(err)

}

/*--------------------------------------------------------------------------------Load booking template with Dropdown of Event----------------------------------------------*/

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

/*-----------------------------------------------------------------------Booking form With Event Id From Event Details Page---------------------------------------*/
func (s *Server) createBookingByEventId(w http.ResponseWriter, r *http.Request) {
	log.Println("Booking : Create Method")
	data := BookingFormData{
		CSRFField: csrf.TemplateField(r),
	}
	id := r.FormValue("id")
	s.loadBookingTemplateByEventId(w, r, data, id)
}

/*--------------------------------------------------------------------------------Load Booking form With Event Id----------------------------------------------*/
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
	err = tmpl.Execute(w, tempData)
	ExcutionTemplateError(err)
}

/*--------------------------------------------------------------------------------Lode Booking form With Event Id----------------------------------------------*/

func (s *Server) saveBookingByEventId(w http.ResponseWriter, r *http.Request) {
	ParseFormData(r)
	var form storage.Booking
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("Booking Page Decoding Error")
	}
	id := form.EventId
	t := IntToStringConversion(id)
	_, err := s.store.GetDataById(t)
	form.UserId = 1
	_, err = s.store.DecrementRemainingTicketById(form.EventId, form.NumberOfTicket) // decrement value as user's input
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
		s.loadBookingTemplateByEventId(w, r, data,t)
		return
	}
	_, err = s.store.CreateBooking(form)
	UnableToGetData(err)
	http.Redirect(w, r, "/booking/boucher", http.StatusSeeOther)
}
