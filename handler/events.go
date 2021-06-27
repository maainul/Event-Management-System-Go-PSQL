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
	EventType  []storage.EventType
	Speakers   []storage.Speakers
}

/*--------------------------------------------------GET EVENT ------------------------------------*/

func (s *Server) getEvents(w http.ResponseWriter, r *http.Request) {

	tmp := s.templates.Lookup("events.html")

	if tmp == nil {
		log.Println("Unable to look event ")
		return
	}
	et, err := s.store.GetEvent()

	//	fmt.Printf("%+v", et)

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
	//	fmt.Printf("Save Event Data = %#v", id)

	http.Redirect(w, r, "/event", http.StatusSeeOther)

}

/* -----------------------------------------Load Create Tempalte Handler------------------------------------------------------------*/

func (s *Server) loadCreateEventTemplate(w http.ResponseWriter, r *http.Request, form EventFormData) {
	tmpl := s.templates.Lookup("event-form.html")
	if tmpl == nil {
		log.Println("Unable to find form")
		return
	}

	et, err := s.store.GetEventType()
	sp, err := s.store.GetSpeakers()
	//fmt.Printf("%+v", et)
	if err != nil {
		log.Println("Unable to get event type.  ", err)
	}
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
		/* log.Println("Error executing template", err)
		return */
	}

}

/* ----------------Show Event Details By ID----------------------------------*/

func (s *Server) eventDetails(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		log.Println("Not found")
	}
	tmp := s.templates.Lookup("event_details.html")
	if tmp == nil {
		log.Println("Unable to load event details page.")
		return
	}
	et, err := s.store.GetDataById(id)
	fmt.Println("Execte from line 157.", et)
	if err != nil {
		log.Println("Unable to get event type.  ", err)
	}
	/* if et == true {

	}
	et.Status =  */
	err = tmp.Execute(w, et)
	if err != nil {
		log.Println("Error executing template:", err)
		return
	}
}

/* -----------------------------------------Time Converter------------------------------------------------------------*/

var timeConverter = func(value string) reflect.Value {
	const shortForm = "2006-01-02"

	if v, err := time.Parse(shortForm, value); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.Value{} // this is the same as the private const invalidType
}
