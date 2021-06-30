package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/csrf"
)

type Login struct {
	Email    string
	Password string
}

type LoginTempData struct {
	CSRFField  template.HTML
	Form       Login
	FormErrors map[string]string
}

func (l Login) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email, validation.Required, is.Email),
		validation.Field(&l.Password, validation.Required, validation.Length(6, 12)),
	)
}

func (s *Server) getLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Method: getLogin")
	formData := LoginTempData{
		CSRFField: csrf.TemplateField(r),
	}
	s.loadLoginTemplate(w, r, formData)
}

func (s *Server) postLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Method: postLogin")
	// pase login information
	if err := r.ParseForm(); err != nil {
		log.Fatalln("parsing error")
	}

	var form Login
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("decoding error")
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
		data := LoginTempData{
			CSRFField:  csrf.TemplateField(r),
			Form:       form,
			FormErrors: vErrs,
		}
		s.loadLoginTemplate(w, r, data)
		return
	}

	/* 	// call database for match email and password
	   	if err = bcrypt.CompareHashAndPassword([]byte(form.Password), []byte(creds.Password)); err != nil {
	   		// If the two passwords don't match, return a 401 status
	   		w.WriteHeader(http.StatusUnauthorized)
	   	} */
	emailandPasswordStruct := s.store.GetUserEmailAndPass(form.Email, form.Password)
	sValue := emailandPasswordStruct.ID //user id

	if emailandPasswordStruct.Email == "" && emailandPasswordStruct.Password == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	if emailandPasswordStruct.IsAdmin == true {
		session, _ := s.session.Get(r, "event_management_app")
		session.Values["user_id"] = sValue // user id
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/auth/admin-home", http.StatusSeeOther)
		fmt.Println("Admin is true 93")
	}

	if emailandPasswordStruct.IsAdmin == false {
		session, _ := s.session.Get(r, "event_management_app")
		session.Values["user_id"] = sValue // customized value
		if err := session.Save(r, w); err != nil {
			log.Fatalln("error while saving user id into session")
		}
		http.Redirect(w, r, "/event", http.StatusSeeOther)
		fmt.Println("Admin is true 101")
	}

}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.session.Get(r, "event_management_app")
	session.Values["user_id"] = 0
	session.Save(r, w)
	http.Redirect(w, r, "/event", http.StatusSeeOther)
}

func (s *Server) loadLoginTemplate(w http.ResponseWriter, r *http.Request, form LoginTempData) {
	tmp := s.templates.Lookup("login.html")
	if err := tmp.Execute(w, form); err != nil {
		log.Println("Error executing template :", err)
		return
	}
}
