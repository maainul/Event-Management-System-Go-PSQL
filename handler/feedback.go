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
	feedback struct {
		Feedback []storage.Feedback
	}
)

type FeedbackFormData struct {
	CSRFField  template.HTML
	Form       storage.Feedback
	FormErrors map[string]string
}

func (s *Server) getFeedback(w http.ResponseWriter, r *http.Request) {
	tmp := s.templates.Lookup("feedback_list.html")
	UnableToFindHtmlTemplate(tmp)
	fb, err := s.store.GetFeedback()
	UnableToGetData(err)
	tempData := feedback{
		Feedback: fb,
	}
	err = tmp.Execute(w, tempData)
	ExcutionTemplateError(err)

}

func (s *Server) createFeedback(w http.ResponseWriter, r *http.Request) {
	log.Println("Method : Create Feedback called.")
	data := FeedbackFormData{
		CSRFField: csrf.TemplateField(r),
	}
	s.loadFeedbackTemplate(w, r, data)
}

func (s *Server) saveFeedback(w http.ResponseWriter, r *http.Request) {
	ParseFormData(r)
	var form storage.Feedback
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("Decoding error")
	}
	intVar := SessionUserId(s, r)
	form.UserId = int32(intVar)
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
		data := FeedbackFormData{
			CSRFField:  csrf.TemplateField(r),
			Form:       form,
			FormErrors: vErrs,
		}
		s.loadFeedbackTemplate(w, r, data)
		return
	}
	_, err := s.store.CreateFeedback(form)
	UnableToInsertData(err)
	http.Redirect(w, r, "/event", http.StatusSeeOther)
}

func (s *Server) loadFeedbackTemplate(w http.ResponseWriter, r *http.Request, form FeedbackFormData) {
	tmpl := s.templates.Lookup("feedback-form.html")
	UnableToFindHtmlTemplate(tmpl)
	err := tmpl.Execute(w, form)
	ExcutionTemplateError(err)

}
