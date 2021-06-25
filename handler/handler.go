package handler

import (
	"Event-Management-System-Go-PSQL/storage/postgres"
	"html/template"
	"net/http"

	"github.com/Masterminds/sprig"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type (
	Server struct {
		templates *template.Template
		store     *postgres.Storage
		decoder   *schema.Decoder
	}
)

func NewServer(st *postgres.Storage, decoder *schema.Decoder) (*mux.Router, error) {

	s := &Server{
		store:   st,
		decoder: decoder,
	}

	if err := s.parseTemplates(); err != nil {
		return nil, err
	}

	r := mux.NewRouter()

	//	r.Use(csrf.Protect([]byte("Secure and safe token"), csrf.Secure(false)))

	csrf.Protect([]byte("keep-it-secret-keep-it-safe-----"), csrf.Secure(false))(r)

	/* staic files Handler */
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/"))))

	/* index Handler */
	r.HandleFunc("/", s.getHome).Methods("GET")

	/* Event Type Handlers */
	r.HandleFunc("/event-type", s.getEventType).Methods("GET")
	r.HandleFunc("/event-type/create", s.createEventType).Methods("GET")
	r.HandleFunc("/event-type/create", s.saveEventType).Methods("POST")

	/* Speakers Handlers */
	r.HandleFunc("/speaker", s.getSpeakers).Methods("GET")
	r.HandleFunc("/speaker/create", s.speakerForm).Methods("GET")

	/* Event Handlers */
	r.HandleFunc("/event", s.getEvents).Methods("GET")
	//	r.HandleFunc("/event/create", s.createEvent).Methods("GET")

	/* Feedback Handlers */
	r.HandleFunc("/feedback", s.getFeedback).Methods("GET")
	r.HandleFunc("/feedback/create", s.createFeedback).Methods("GET")
	r.HandleFunc("/feedback/create", s.saveFeedback).Methods("POST")

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
