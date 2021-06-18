package handler

import (
	"Event-Management-System-Go-PSQL/storage"
	"fmt"
	"log"
	"net/http"
)

type (
	speakers struct {
		Speakers []storage.Speakers
	}
)

func (s *Server) getSpeakers(w http.ResponseWriter, r *http.Request) {

	tmp := s.templates.Lookup("speakers_list.html")

	if tmp == nil {
		log.Println("Unable to look speakers_list.html")
		return
	}
	et, err := s.store.GetSpeakers()

	fmt.Printf("%+v", et)

	if err != nil {
		log.Println("Unable to get event type.  ", err)
	}

	tempData := speakers{
		Speakers: et,
	}

	err = tmp.Execute(w, tempData)
	if err != nil {
		log.Println("Error executing tempalte:", err)
		return
	}

}
