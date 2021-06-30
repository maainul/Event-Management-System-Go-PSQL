package handler

import (
	"Event-Management-System-Go-PSQL/storage"
	"html/template"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
)

type eventTypeData struct {
	EventType []storage.EventType
}
type EventTypeFormData struct {
	CSRFField  template.HTML
	Form       storage.EventType
	FormErrors map[string]string
}

/*------------------------------------------------ Get all event Type----------------------------------*/

func (s *Server) getEventType(w http.ResponseWriter, r *http.Request) {
	tmp := s.templates.Lookup("event_type_list.html")
	UnableToFindHtmlTemplate(tmp)
	et, err := s.store.GetEventType()
	UnableToGetData(err)
	tempData := eventTypeData{
		EventType: et,
	}
	err = tmp.Execute(w, tempData)
	ExcutionTemplateError(err)
}

/*------------------------------------------------Create event Type Form----------------------------------*/

func (s *Server) createEventType(w http.ResponseWriter, r *http.Request) {
	log.Println("Method : createEventType")
	data := EventTypeFormData{
		CSRFField: csrf.TemplateField(r),
	}
	s.loadCreateEventTypeTemplate(w, r, data)

}
func (s *Server) saveEventType(w http.ResponseWriter, r *http.Request) {
	log.Println("Method : Save event Type Called ")
	ParseFormData(r)
	var form storage.EventType
	// Decoding Data
	err := s.decoder.Decode(&form, r.PostForm)
	if err != nil {
		log.Fatalln("Form Decoding Error")
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

		data := EventTypeFormData{
			CSRFField:  csrf.TemplateField(r),
			Form:       form,
			FormErrors: vErrs,
		}
		s.loadCreateEventTypeTemplate(w, r, data)
		return
	}
	_, err = s.store.CreateEventType(form)
	UnableToInsertData(err)
	http.Redirect(w, r, "/auth/event-type", http.StatusSeeOther)
}

func (s *Server) loadCreateEventTypeTemplate(w http.ResponseWriter, r *http.Request, form EventTypeFormData) {
	tmpl := s.templates.Lookup("event-type-form.html")
	UnableToFindHtmlTemplate(tmpl)
	err := tmpl.Execute(w, form)
	ExcutionTemplateError(err)

}
