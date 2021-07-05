package handler

import (
	"net/http"
)

/*---------------------------------------------------------Home/Index Template------------------ ------------------------------------*/

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	tmpl := s.templates.Lookup("events.html")
	UnableToFindHtmlTemplate(tmpl)
	et, err := s.store.GetEvent()
	UnableToGetData(err)
	tempData := events{
		Events: et,
	}
	err = tmpl.Execute(w, tempData)
	ExcutionTemplateError(err)

}
