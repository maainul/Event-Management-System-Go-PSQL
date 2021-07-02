package handler

import (
	"Event-Management-System-Go-PSQL/storage"
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

/*---------------------------------------------------------------------Booking form With Event Id From Event Details Page-------------------------------*/
func (s *Server) createBookingByEventId(w http.ResponseWriter, r *http.Request) {
	log.Println("Booking : Create Method")
	data := BookingFormData{
		CSRFField: csrf.TemplateField(r),
	}
	id := r.FormValue("id")
	s.loadBookingTemplateByEventId(w, r, data, id)
}

/*-----------------------------------------------------------------------------Load Booking form With Event Id-----------------------------------------*/
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

/*---------------------------------------------------------------------------Lode Booking form With Event Id------------------------------------------*/

func (s *Server) saveBookingByEventId(w http.ResponseWriter, r *http.Request) {
	ParseFormData(r)
	var form storage.Booking
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("Booking Page Decoding Error")
	}
	id := form.EventId
	t := IntToStringConversion(id)
	result, err := s.store.GetDataById(t)
	intVar := SessionUserId(s, r)
	form.UserId = int32(intVar)
	if form.NumberOfTicket > result.TicketRemaining {
		log.Println("Input Ticket is More than Remaining Ticket")
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
	}
	ticket_price := result.PerPersonPrice * form.NumberOfTicket
	form.TotalAmount = ticket_price

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
		s.loadBookingTemplateByEventId(w, r, data, t)
		return
	}
	_, err = s.store.CreateBooking(form)
	UnableToGetData(err)
	http.Redirect(w, r, "/booking/boucher", http.StatusSeeOther)
}

/*--------------------------------------------------------Booking Information/ Boucher/Invoice of Booking---------------------------------------------*/
func (s *Server) bookingBoucher(w http.ResponseWriter, r *http.Request) {
	tmpl := s.templates.Lookup("boucher.html")
	UnableToFindHtmlTemplate(tmpl)
	err := tmpl.Execute(w, nil)
	ExcutionTemplateError(err)

}
