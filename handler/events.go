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

type (
	events struct {
		Events        []storage.Events
		CountAllEvent int32
	}
)

type EventFormData struct {
	CSRFField  template.HTML
	Form       storage.Events
	FormErrors map[string]string
	EventType  []storage.EventType
	Speakers   []storage.Speakers
}

/*--------------------------------------------------GET EVENT ------------------------------------*/

func (s *Server) getEvents(w http.ResponseWriter, r *http.Request) {
	tmp := s.templates.Lookup("events.html")
	UnableToFindHtmlTemplate(tmp)
	et, err := s.store.GetEvent()
	UnableToGetData(err)
	tempData := events{
		Events: et,
	}
	err =  tmp.Execute(w, tempData)
}

/*--------------------------------------------------GET EVENT BY ADMIN ------------------------------------*/

func (s *Server) authGetEvents(w http.ResponseWriter, r *http.Request) {
	tmp := s.templates.Lookup("admin-home.html")
	et, err := s.store.GetEvent()
	UnableToGetData(err)
	ce := s.store.CountEvent()
	tempData := events{
		Events:        et,
		CountAllEvent: ce,
	}
	tempData.CountAllEvent = ce

	err = tmp.Execute(w, tempData)
	ExcutionTemplateError(err)

}

/* -----------------------------------------Create Event Handler------------------------------------------------------------*/
func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("Method : Create Event ")
	data := EventFormData{
		CSRFField: csrf.TemplateField(r),
	}
	s.loadCreateEventTemplate(w, r, data)
}

/* -----------------------------------------Save Event Handler------------------------------------------------------------*/

func (s *Server) saveEvent(w http.ResponseWriter, r *http.Request) {
	ParseFormData(r)
	var form storage.Events
	if err := s.decoder.Decode(&form, r.Form); err != nil {
		fmt.Println(err)
	}
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
		data := EventFormData{
			CSRFField:  csrf.TemplateField(r),
			Form:       form,
			FormErrors: vErrs,
		}
		s.loadCreateEventTemplate(w, r, data)
		return
	}
	_, err := s.store.CreateEvent(form)
	UnableToInsertData(err)
	http.Redirect(w, r, "/auth/event", http.StatusSeeOther)

}

/* -----------------------------------------Load Create Tempalte Handler------------------------------------------------------------*/

func (s *Server) loadCreateEventTemplate(w http.ResponseWriter, r *http.Request, form EventFormData) {
	tmpl := s.templates.Lookup("event-form.html")
	UnableToFindHtmlTemplate(tmpl)
	et, err := s.store.GetEventType()
	UnableToGetData(err)
	sp, err := s.store.GetSpeakers()
	UnableToGetData(err)
	tempData := EventFormData{
		Form:       storage.Events{},
		FormErrors: map[string]string{},
		EventType:  et,
		Speakers:   sp,
	}
	err = tmpl.Execute(w, tempData)
	ExcutionTemplateError(err)
}

/* ----------------Show Event Details By ID----------------------------------*/

func (s *Server) eventDetails(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		log.Println("Not found")
	}
	tmp := s.templates.Lookup("event_details.html")
	UnableToFindHtmlTemplate(tmp)
	et, err := s.store.GetDataById(id)
	UnableToGetData(err)
	err = tmp.Execute(w, et)
	ExcutionTemplateError(err)

}
