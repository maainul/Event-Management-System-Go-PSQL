package handler

import (
	"Event-Management-System-Go-PSQL/storage"
	"fmt"
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

func (s *Server) createFeedback(w http.ResponseWriter, r *http.Request) {
	log.Println("Method : Create Feedback called.")

	data := FeedbackFormData{
		CSRFField: csrf.TemplateField(r),
	}

	s.loadFeedbackTemplate(w, r, data)

}

func (s *Server) saveFeedback(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalln("Parsing error")
	}

	var form storage.Feedback
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("Decoding error")
	}
	fmt.Printf("%#v", form)
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
	id, err := s.store.CreateFeedback(form)
	if err != nil {
		log.Fatalln("Unable to save data :", err)

	}
	fmt.Printf("%#v", id)

	http.Redirect(w, r, "/feedback", http.StatusSeeOther)
}

func (s *Server) loadFeedbackTemplate(w http.ResponseWriter, r *http.Request, form FeedbackFormData) {
	tmpl := s.templates.Lookup("feedback-form.html")
	if tmpl == nil {
		log.Println("Unable to find form")
		return
	}
	if err := tmpl.Execute(w, form); err != nil {
		log.Println("Error executing tempalte : ", err)
		return
	}
}
