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

func (s *Server) speakerForm(w http.ResponseWriter, r *http.Request) {
	tmp := s.templates.Lookup("speaker_form.html")

	if tmp == nil {
		log.Println("Unable to find form")
		return
	}

	err := tmp.Execute(w, tmp)
	if err != nil {
		log.Println("Error executing template", err)
		return
	}

}

func (s *Server) speakerCreateProcesss(w http.ResponseWriter, r *http.Request) {
	createQuery := `INSERT INTO speakers (first_name, last_name, phone, address,username,email,created_at,updated_at) VALUES (:first_name, :last_name, :phone, :address,:username,:email,:created_at,:updated_at)`

	var bk storage.Speakers
	bk.FirstName = r.FormValue("first_name")
	bk.LastName = r.FormValue("last_name")
	bk.Phone = r.FormValue("phone")

}

