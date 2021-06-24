package handler

import (
	"Event-Management-System-Go-PSQL/storage"
	"fmt"
	"log"
	"net/http"
)

type (
	eventTypeData struct {
		EventType []storage.EventType
	}
)

func (s *Server) getEventType(w http.ResponseWriter, r *http.Request) {

	tmp := s.templates.Lookup("event_type_list.html")

	if tmp == nil {
		log.Println("Unable to look event type.html")
		return
	}
	et, err := s.store.GetEventType()

	fmt.Printf("%+v", et)

	if err != nil {
		log.Println("Unable to get event type.  ", err)
	}

	tempData := eventTypeData{
		EventType: et,
	}

	err = tmp.Execute(w, tempData)
	if err != nil {
		log.Println("Error executing tempalte:", err)
		return
	}
}

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	tmpl := s.templates.Lookup("event-form.html")

	err := tmpl.Execute(w, tmpl)
	if err != nil {
		log.Println("Error executing template", err)
		return
	}
}
