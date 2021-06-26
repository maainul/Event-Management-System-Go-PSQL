package handler

import (
	"Event-Management-System-Go-PSQL/storage"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/schema"
)

type (
	events struct {
		Events []storage.Events
	}
)

type EventFormData struct {
	CSRFField  template.HTML
	Form       storage.Events
	FormErrors map[string]string
}

func (s *Server) getEvents(w http.ResponseWriter, r *http.Request) {

	tmp := s.templates.Lookup("events.html")

	if tmp == nil {
		log.Println("Unable to look event ")
		return
	}
	et, err := s.store.GetEvent()

	fmt.Printf("%+v", et)

	if err != nil {
		log.Println("Unable to get event type.  ", err)
	}

	tempData := events{
		Events: et,
	}

	err = tmp.Execute(w, tempData)
	if err != nil {
		log.Println("Error executing tempalte:", err)
		return
	}
}

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("Method : Create Event ")

	data := EventFormData{
		CSRFField: csrf.TemplateField(r),
	}
	s.loadCreateEventTemplate(w, r, data)

}

var timeConverter = func(value string) reflect.Value {
	const shortForm =  "2006-01-02"
    
	if v, err := time.Parse(shortForm,value); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.Value{} // this is the same as the private const invalidType
}

func (s *Server) saveEvent(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Fatalln("Parsing error")
	}

	var form storage.Events
	/* 	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
	   		log.Fatalln("Decoding Error in line 73")
	   	}
	*/
	decoder := schema.NewDecoder()
	decoder.RegisterConverter(time.Time{}, timeConverter)

	if err := decoder.Decode(&form, r.Form); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v", form)
	if form.EventName == "" || form.NumberOfGuest == 0 || form.PerPersonPrice == 0 || form.EventDate.IsZero() || form.StartTime.IsZero() || form.EndTime.IsZero() {
		data := EventFormData{
			CSRFField: csrf.TemplateField(r),
			Form:      form,
			FormErrors: map[string]string{
				"EventName": "Event name cannot be null",
			},
		}

		s.loadCreateEventTemplate(w, r, data)
	}

	id, err := s.store.CreateEvent(form)
	if err != nil {
		log.Fatalln("Unable to save data:", err)
	}
	fmt.Printf("Save Event Data = %#v", id)

	http.Redirect(w, r, "/event", http.StatusSeeOther)

}
func (s *Server) loadCreateEventTemplate(w http.ResponseWriter, r *http.Request, form EventFormData) {
	tmpl := s.templates.Lookup("event-form.html")
	if tmpl == nil {
		log.Println("Unable to find form")
		return
	}
	err := tmpl.Execute(w, form)
	if err != nil {
		log.Println("Error executing template", err)
		return
	}
}
