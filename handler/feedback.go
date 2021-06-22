package handler

import (
	"Event-Management-System-Go-PSQL/storage"
	"fmt"
	"log"
	"net/http"
)

type (
	feedback struct {
		Feedback []storage.Feedback
	}
)

func (s *Server) getFeedback(w http.ResponseWriter, r *http.Request) {

	tmp := s.templates.Lookup("feedback_list.html")

	if tmp == nil {
		log.Println("Unable to look feedback_list.html")
		return
	}
	fb, err := s.store.GetFeedback()

	fmt.Printf("%+v", fb)

	if err != nil {
		log.Println("Unable to get event type.  ", err)
	}

	tempData := feedback{
		Feedback: fb,
	}

	err = tmp.Execute(w, tempData)
	if err != nil {
		log.Println("Error executing tempalte:", err)
		return
	}

}

/*
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
*/
