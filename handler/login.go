package handler

import (
	"Event-Management-System-Go-PSQL/storage"
	"html/template"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
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

/*------------------------------------------------------------GET Login Info ------------------------------------------------*/

func (s *Server) getLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Method: getLogin")
	formData := LoginTempData{
		CSRFField: csrf.TemplateField(r),
	}
	s.loadLoginTemplate(w, r, formData)
}

/*----------------------------------------------------------------POST: Login Save -----------------------------------------*/

func (s *Server) postLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Method: postLogin")
	ParseFormData(r)
	var form Login
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("decoding error")
	}
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
	email := form.Email
	result := s.store.GetUserInfo(email)
	ComparePassword(result, form, w, r)
	sessionUID := result.ID
	isAdmin := result.IsAdmin
	session, _ := s.session.Get(r, "event_management_app")
	session.Values["user_id"] = IntToStringConversion(sessionUID)
	session.Values["is_admin"] = isAdmin
	if err := session.Save(r, w); err != nil {
		log.Fatalln("error while saving user id into session")
	}
	LoginRedirect(isAdmin, w, r)
}

/*--------------------------------------------------------------Logout ------------------------------------------------------*/
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.session.Get(r, "event_management_app")
	session.Values["user_id"] = ""
	session.Values["is_admin"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*---------------------------------------------------------Login Template------------------ ------------------------------------*/
func (s *Server) loadLoginTemplate(w http.ResponseWriter, r *http.Request, form LoginTempData) {
	tmp := s.templates.Lookup("login.html")
	err := tmp.Execute(w, form)
	ExcutionTemplateError(err)
}

/*---------------------------------------------------------Compare Two Password------------------ ------------------------------------*/
func ComparePassword(result *storage.User, form Login, w http.ResponseWriter, r *http.Request) {
	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(form.Password)); err != nil {
		log.Println("Password does not match.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

/*---------------------------------------------------------Login Redirect Page------------------ ------------------------------------*/
func LoginRedirect(isAdmin bool, w http.ResponseWriter, r *http.Request) {
	if isAdmin == true {
		http.Redirect(w, r, "/auth/admin-home", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}
