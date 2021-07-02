package handler

import (
	"Event-Management-System-Go-PSQL/storage/postgres"
	"html/template"
	"net/http"

	"github.com/Masterminds/sprig"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

type (
	Server struct {
		templates *template.Template
		store     *postgres.Storage
		decoder   *schema.Decoder
		session   *sessions.CookieStore
	}
)

func NewServer(st *postgres.Storage, decoder *schema.Decoder, session *sessions.CookieStore) (*mux.Router, error) {

	s := &Server{
		store:   st,
		decoder: decoder,
		session: session,
	}

	if err := s.parseTemplates(); err != nil {
		return nil, err
	}

	r := mux.NewRouter()

	csrf.Protect([]byte("keep-it-secret-keep-it-safe-----"), csrf.Secure(false))(r)

	/* staic files Handler */
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/"))))

	/* Login Handler */
	r.HandleFunc("/login", s.getLogin).Methods("GET")
	r.HandleFunc("/login", s.postLogin).Methods("POST")
	r.HandleFunc("/logout", s.logout).Methods("GET")

	/* User Create Handler */
	r.HandleFunc("/user/create", s.createUser).Methods("GET")
	r.HandleFunc("/user/create", s.saveUser).Methods("POST")

	/* Speakers Handlers User*/
	r.HandleFunc("/speaker", s.getSpeaker).Methods("GET")
	r.HandleFunc("/forbidden", s.getForbiddenPage).Methods("GET")

	/* Event Handlers User */
	r.HandleFunc("/event", s.getEvents).Methods("GET")
	r.HandleFunc("/", s.getEvents).Methods("GET")
	r.HandleFunc("/event/show", s.eventDetails).Methods("GET")

	/* Event Type Handlers User */
	r.HandleFunc("/event-type", s.getEventType).Methods("GET")

	/*------------------------------------------------USER AUTHENTICATION----------------------------------*/

	ur := r.NewRoute().Subrouter()
	ur.Use(s.userAuthMiddleware)

	/* Feedback Handler User*/
	ur.HandleFunc("/feedback/create", s.createFeedback).Methods("GET")
	ur.HandleFunc("/feedback/create", s.saveFeedback).Methods("POST")

	/* Booking From Events Details*/
	ur.HandleFunc("/booking/show", s.createBookingByEventId).Methods("GET")
	ur.HandleFunc("/booking/show/create", s.saveBookingByEventId).Methods("POST")

	ur.HandleFunc("/booking/boucher", s.bookingBoucher).Methods("GET")

	/*------------------------------------------------ADMIN AUTHENTICATION----------------------------------*/
	ar := r.NewRoute().Subrouter()
	ar.Use(s.adminAuthMiddleware)

	/* Admin Home  Handler */
	ar.HandleFunc("/auth/admin-home", s.getAdminHomePage).Methods("GET")

	/* Auth Event Type Handlers */
	ar.HandleFunc("/auth/event-type", s.getEventType).Methods("GET")
	ar.HandleFunc("/auth/event-type/create", s.createEventType).Methods("GET")
	ar.HandleFunc("/auth/event-type/create", s.saveEventType).Methods("POST")

	/* Event ype Handlers */
	ar.HandleFunc("/auth/event", s.authGetEvents).Methods("GET")
	ar.HandleFunc("/auth/event/create", s.createEvent).Methods("GET")
	ar.HandleFunc("/auth/event/create", s.saveEvent).Methods("Post")
	ar.HandleFunc("/auth/event/show", s.eventDetails).Methods("GET")

	/* Feedback Handlers */
	ar.HandleFunc("/auth/feedback", s.getFeedback).Methods("GET")

	/* User Handlers */
	ar.HandleFunc("/auth/user", s.getUser).Methods("GET")

	/* Speaker Handler */
	ar.HandleFunc("/auth/speaker/create", s.createSpeaker).Methods("GET")
	ar.HandleFunc("/auth/speaker/create", s.saveSpeaker).Methods("POST")

	return r, nil
}

func (s *Server) parseTemplates() error {
	templates := template.New("templates").Funcs(sprig.FuncMap())
	tmpl, err := templates.ParseGlob("assets/templates/*.html")
	if err != nil {
		return err
	}
	s.templates = tmpl
	return nil
}

/*------------------------------------------------ADMIN AUTHENTICATION MIDDLEWARE-----------------------------------*/
func (s *Server) adminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SessionCheckAndRedirect(s, r, next, w, true)
	})
}

/*------------------------------------------------USER AUTHENTICATION MIDDLEWARE----------------------------------*/
func (s *Server) userAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SessionCheckAndRedirect(s, r, next, w, false)
	})
}

/*------------------------------------------------Session Information Checker ----------------------------------*/
func SessionCheckAndRedirect(s *Server, r *http.Request, next http.Handler, w http.ResponseWriter, user bool) {
	uid, user_type := GetSetSessionValue(s, r)
	if uid != "" && user_type == user {
		next.ServeHTTP(w, r)
	} else {
		http.Redirect(w, r, "/forbidden", http.StatusSeeOther)
	}
}

func GetSetSessionValue(s *Server, r *http.Request) (interface{}, interface{}) {
	session, _ := s.session.Get(r, "event_management_app")
	uid := session.Values["user_id"]
	user_type := session.Values["is_admin"]
	return uid, user_type
}
