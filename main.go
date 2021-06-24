package main

import (
	"Event-Management-System-Go-PSQL/handler"
	"Event-Management-System-Go-PSQL/storage/postgres"
	"log"
	"net/http"

	"time"

	"github.com/gorilla/schema"
)

func main() {
	dbString := newDBFromConfig()
	store, err := postgres.NewStorage(dbString)

	if err != nil {
		log.Fatal(err)
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	r, err := handler.NewServer(store, decoder)
	if err != nil {
		log.Fatal("Handler not Found")
	}

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

}

func newDBFromConfig() string {
	dbParams := " " + "user=postgres"
	dbParams += " " + "host=localhost"
	dbParams += " " + "port=5432"
	dbParams += " " + "dbname=eventDB"
	dbParams += " " + "password=0"
	dbParams += " " + "sslmode=disable"

	return dbParams
}
