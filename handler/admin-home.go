package handler

import (
	"net/http"
)

func (s *Server) getAdminHomePage(w http.ResponseWriter, r *http.Request) {

	tmp, result := s.loadTemplate("admin-home.html")
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
