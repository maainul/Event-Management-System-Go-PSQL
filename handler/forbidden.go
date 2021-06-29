package handler

import (
	"log"
	"net/http"
)

func (s *Server) getForbiddenPage(w http.ResponseWriter, r *http.Request) {
	tmp := s.templates.Lookup("forbidden.html")
	if tmp == nil {
		log.Println("Unable to look feedback_list.html")
		return
	}
	err := tmp.Execute(w, nil)
	if err != nil {
		log.Println("Error executing tempalte:", err)
		return
	}

}
