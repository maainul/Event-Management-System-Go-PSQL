package handler

import (
	"Event-Management-System-Go-PSQL/storage"
	"fmt"
	"html/template"
	"log"
	"net/http"

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
	tmp, result := s.loadTemplate("events.html")
	if result {
		return
	}
	et, err := s.store.GetEvent()
	unableToGetData(err, "Unable to Get Event Data")
	tempData := events{
		Events: et,
	}
	res := execute(tmp, w, tempData)
	if res {
		return
	}
}

/*--------------------------------------------------GET EVENT BY ADMIN ------------------------------------*/

func (s *Server) authGetEvents(w http.ResponseWriter, r *http.Request) {
	tmp, result := s.loadTemplate("admin-home.html")
	if result {
		return
	}
	et, err := s.store.GetEvent()
	unableToGetData(err, "Unable to Get Event")
	ce := s.store.CountEvent()
	tempData := events{
		Events:        et,
		CountAllEvent: ce,
	}
	tempData.CountAllEvent = ce
	println(tempData.CountAllEvent)
	res := execute(tmp, w, tempData)
	if res {
		return
	}
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
	if err := r.ParseForm(); err != nil {
		log.Fatalln("Parsing error")
	}
	var form storage.Events
	if err := s.decoder.Decode(&form, r.Form); err != nil {
		fmt.Println(err)
	}
	if form.EventName == "" || form.NumberOfGuest == 0 || form.PerPersonPrice == 0 || form.EventDate == "" || form.EventStartTime == "" || form.EventEndTime == "" {
		data := EventFormData{
			CSRFField: csrf.TemplateField(r),
			Form:      form,
			FormErrors: map[string]string{
				"EventName": "Event name cannot be null",
			},
		}
		s.loadCreateEventTemplate(w, r, data)
	}
	_, err := s.store.CreateEvent(form)
	if err != nil {
		log.Fatalln("Unable to save data:", err)
	}
	http.Redirect(w, r, "/event", http.StatusSeeOther)

}

/* -----------------------------------------Load Create Tempalte Handler------------------------------------------------------------*/

func (s *Server) loadCreateEventTemplate(w http.ResponseWriter, r *http.Request, form EventFormData) {
	tmpl, result := s.loadTemplate("event-form.html")
	if result {
		return
	}
	et, err := s.store.GetEventType()
	unableToGetData(err, "Unable to get Event Type")
	sp, err := s.store.GetSpeakers()
	unableToGetData(err, "Unable to get Speakers")
	tempData := EventFormData{
		CSRFField:  "",
		Form:       storage.Events{},
		FormErrors: map[string]string{},
		EventType:  et,
		Speakers:   sp,
	}
	err = tmpl.Execute(w, tempData)
	if err != nil {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
}

/* ----------------Show Event Details By ID----------------------------------*/

func (s *Server) eventDetails(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		log.Println("Not found")
	}
	tmp, result := s.loadTemplate("event_details.html")
	if result {
		return
	}
	et, err := s.store.GetDataById(id)
	unableToGetData(err, "Unable to get event type.")
	err = tmp.Execute(w, et)
	if err != nil {
		log.Println("Error executing template:", err)
		return
	}

}

func execute(tmp *template.Template, w http.ResponseWriter, data events) bool {
	err := tmp.Execute(w, data)
	if err != nil {
		log.Println("Error executing tempalte:", err)
		return true
	}
	return false
}

func (s *Server) loadTemplate(str string) (*template.Template, bool) {
	tmp := s.templates.Lookup(str)
	if tmp == nil {
		log.Println("Unable to Find Template ")
		return nil, true
	}
	return tmp, false
}

func unableToGetData(err error, message string) {
	if err != nil {
		log.Println(message, err)
	}
}
