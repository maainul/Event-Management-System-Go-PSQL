package handler

import (
	"Event-Management-System-Go-PSQL/storage/postgres"
	"html/template"
	"net/http"

	"github.com/Masterminds/sprig"
	"github.com/gorilla/mux"
)

type (
	Server struct {
		templates *template.Template
		store     *postgres.Storage
	}
)

func NewServer(st *postgres.Storage) (*mux.Router, error) {

	s := &Server{
		store: st,
	}

	if err := s.parseTemplates(); err != nil {
		return nil, err
	}

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/"))))
	r.HandleFunc("/", s.getHome).Methods("GET")
	r.HandleFunc("/event-types", s.getEventType).Methods("GET")
	r.HandleFunc("/speakers", s.getSpeakers).Methods("GET")
	r.HandleFunc("/events", s.getEvents).Methods("GET")
	r.HandleFunc("/speaker-form", s.speakerForm).Methods("GET")
	r.HandleFunc("/speaker-create-process", s.speakerCreateProcesss).Methods("POST")

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
