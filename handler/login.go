package handler

import (
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

func (s *Server) getLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Method: getLogin")
	formData := LoginTempData{
		CSRFField: csrf.TemplateField(r),
	}
	s.loadLoginTemplate(w, r, formData)
}

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
	// Compare the stored hashed password, with the hashed version of the password that was received
	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(form.Password)); err != nil {
		log.Println("Password does not match.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	Session_User_ID := result.ID
	Session_Is_Admin := result.IsAdmin

	session, _ := s.session.Get(r, "event_management_app")
	session.Values["user_id"] = Session_User_ID
	session.Values["is_admin"] = Session_Is_Admin
	if err := session.Save(r, w); err != nil {
		log.Fatalln("error while saving user id into session")
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.session.Get(r, "event_management_app")
	session.Values["user_id"] = 0
	session.Values["is_admin"] = false

	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) loadLoginTemplate(w http.ResponseWriter, r *http.Request, form LoginTempData) {
	tmp := s.templates.Lookup("login.html")
	err := tmp.Execute(w, form)
	ExcutionTemplateError(err)
}

/*

func (s *Server) postLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Method: postLogin")
	ParseFormData(r)
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
/*
		   emailandPasswordStruct := s.store.GetUserEmailAndPass(form.Email, form.Password)
		   Session_User_ID := emailandPasswordStruct.ID //user id
		   Session_Is_Admin := emailandPasswordStruct.IsAdmin

		   if emailandPasswordStruct.Email == "" && emailandPasswordStruct.Password == "" {
			   http.Redirect(w, r, "/login", http.StatusSeeOther)
		   }

		   if emailandPasswordStruct.IsAdmin == true {
			   session, _ := s.session.Get(r, "event_management_app")
			   session.Values["user_id"] = Session_User_ID
			   session.Values["is_admin"] = Session_Is_Admin
			   err := session.Save(r, w)
			   if err != nil {
				   http.Error(w, err.Error(), http.StatusInternalServerError)
				   return
			   }
			   fmt.Println("Your are Admin")
			   http.Redirect(w, r, "/auth/admin-home", http.StatusSeeOther) // admin
		   }

		   if emailandPasswordStruct.IsAdmin == false {
			   session, _ := s.session.Get(r, "event_management_app")
			   session.Values["user_id"] = Session_User_ID
			   session.Values["is_admin"] = Session_Is_Admin
			   if err := session.Save(r, w); err != nil {
				   log.Fatalln("error while saving user id into session")
			   }
			   fmt.Println("You are user")
			   http.Redirect(w, r, "/event", http.StatusSeeOther) // user index
			   fmt.Println("This is user : hendler/login.go")
		   }

	   }
*/
