package handler

import (
	"net/http"
)

func (s *Server) getAdminHomePage(w http.ResponseWriter, r *http.Request) {
	s.authGetEvents(w, r)
}
