package handler

import (
	"Event-Management-System-Go-PSQL/storage"
	"html/template"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
)

type (
	speaker struct {
		Speakers []storage.Speakers
	}
)

type SpeakerFormData struct {
	CSRFField  template.HTML
	Form       storage.Speakers
	FormErrors map[string]string
}

func (s *Server) getSpeaker(w http.ResponseWriter, r *http.Request) {
	tmp := s.templates.Lookup("speakers_list.html")
	UnableToFindHtmlTemplate(tmp)
	spk, err := s.store.GetSpeakers()
	UnableToGetData(err)
	tempData := speaker{
		Speakers: spk,
	}
	err = tmp.Execute(w, tempData)
	ExcutionTemplateError(err)
}

func (s *Server) createSpeaker(w http.ResponseWriter, r *http.Request) {
	log.Println("Method : Create speaker called.")
	data := SpeakerFormData{
		CSRFField: csrf.TemplateField(r),
	}
	s.loadSpeakerTemplate(w, r, data)
}

func (s *Server) saveSpeaker(w http.ResponseWriter, r *http.Request) {
	ParseFormData(r)
	// decode form data
	var form storage.Speakers
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("Decoding error")
	}
	// validation
	if err := form.Validate(); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}
		data := SpeakerFormData{
			CSRFField:  csrf.TemplateField(r),
			Form:       form,
			FormErrors: vErrs,
		}
		s.loadSpeakerTemplate(w, r, data)
		return
	}
	// call database query
	_, err := s.store.CreateSpeaker(form)
	UnableToInsertData(err)
	http.Redirect(w, r, "/speaker", http.StatusSeeOther)
}

func (s *Server) loadSpeakerTemplate(w http.ResponseWriter, r *http.Request, form SpeakerFormData) {
	tmpl := s.templates.Lookup("speaker-form.html")
	UnableToFindHtmlTemplate(tmpl)
	err := tmpl.Execute(w, form)
	ExcutionTemplateError(err)
}
