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

	/* Admin Home  Handler */
	r.HandleFunc("/auth/admin-home", s.getAdminHomePage).Methods("GET")

	/* Event Handlers User */
	r.HandleFunc("/event", s.getEvents).Methods("GET")
	r.HandleFunc("/", s.getEvents).Methods("GET")
	r.HandleFunc("/event/show", s.eventDetails).Methods("GET")

	/* Event Type Handlers User */
	r.HandleFunc("/event-type", s.getEventType).Methods("GET")

	/* Feedback Handler User*/
	r.HandleFunc("/feedback/create", s.createFeedback).Methods("GET")
	r.HandleFunc("/feedback/create", s.saveFeedback).Methods("POST")

	/* User Create Handler */
	r.HandleFunc("/user/create", s.createUser).Methods("GET")
	r.HandleFunc("/user/create", s.saveUser).Methods("POST")

	/* Speakers Handlers User*/
	r.HandleFunc("/speaker", s.getSpeaker).Methods("GET")
	r.HandleFunc("/forbidden", s.getForbiddenPage).Methods("GET")

	/* Booking */
	r.HandleFunc("/booking/create", s.createBooking).Methods("GET")
	r.HandleFunc("/booking/create", s.saveBooking).Methods("POST")
	r.HandleFunc("/booking/boucher", s.bookingBoucher).Methods("GET")
	r.HandleFunc("/booking/create/show", s.createBookingByEventId).Methods("GET")
	// r.HandleFunc("/booking/create/show", s.s).Methods("POST")

	/*------------------------------------------------AUTHENTICATION----------------------------------*/
	/* Auth Event Type Handlers */
	r.HandleFunc("/auth/event-type", s.getEventType).Methods("GET")
	r.HandleFunc("/auth/event-type/create", s.createEventType).Methods("GET")
	r.HandleFunc("/auth/event-type/create", s.saveEventType).Methods("POST")

	/* Event ype Handlers */
	r.HandleFunc("/auth/event", s.authGetEvents).Methods("GET")
	r.HandleFunc("/auth/event/create", s.createEvent).Methods("GET")
	r.HandleFunc("/auth/event/create", s.saveEvent).Methods("Post")
	r.HandleFunc("/auth/event/show", s.eventDetails).Methods("GET")

	/* Feedback Handlers */
	r.HandleFunc("/auth/feedback", s.getFeedback).Methods("GET")

	/* User Handlers */
	r.HandleFunc("/auth/user", s.getUser).Methods("GET")

	/* Speaker Handler */
	r.HandleFunc("/auth/speaker", s.getSpeaker).Methods("GET")
	r.HandleFunc("/speaker/create", s.createSpeaker).Methods("GET")
	r.HandleFunc("/speaker/create", s.saveSpeaker).Methods("POST")

	/* Middleware */
	ar := r.NewRoute().Subrouter()
	ar.Use(s.authMiddleware)

	return r, nil
}

func (s *Server) parseTemplates() error {
	templates := template.New("templates").Funcs(template.FuncMap{
		"strrev": func(str string) string {
			n := len(str)
			runes := make([]rune, n)
			for _, rune := range str {
				n--
				runes[n] = rune
			}
			return string(runes[n:])
		},
	}).Funcs(sprig.FuncMap())

	tmpl, err := templates.ParseGlob("assets/templates/*.html")
	if err != nil {
		return err
	}
	s.templates = tmpl
	return nil
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.session.Get(r, "event_management_system")
		value := session.Values["user_id"]
		if _, ok := value.(string); ok {
			next.ServeHTTP(w, r)
		} else {
			//http.Error(w, "Forbidden", http.StatusForbidden)
			http.Redirect(w, r, "/forbidden", http.StatusSeeOther)
		}
	})
}
