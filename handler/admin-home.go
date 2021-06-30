package handler

import (
	"net/http"
)

func (s *Server) getAdminHomePage(w http.ResponseWriter, r *http.Request) {

	tmp := s.templates.Lookup("admin-home.html")
	UnableToFindHtmlTemplate(tmp)
	et, err := s.store.GetEvent()
	UnableToGetData(err)
	tempData := events{
		Events: et,
	}
	err = tmp.Execute( w, tempData)
	ExcutionTemplateError(err)
}
