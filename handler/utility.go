package handler

import (
	"Event-Management-System-Go-PSQL/storage"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Parsing Data
func ParseFormData(r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalln("Form Data Parsing Error")
	}
}

// Decoding Data
func DecodeFormData(s *Server, form storage.Booking, r *http.Request) {
	err := s.decoder.Decode(&form, r.PostForm)
	if err != nil {
		log.Fatalln("Form Decoding Error")
	}
}

// Unable to find data
func UnableToGetData(err error) {
	if err != nil {
		log.Println("Unable to find data")
	}
}

// if insert is not possible then give this message
func UnableToInsertData(err error) {
	if err != nil {
		log.Fatalln("Unable to Insert Data ", err)
	}
}

// Integer to string conversion
func IntToStringConversion(id int32) string {
	t := strconv.Itoa(int(id))
	return t
}

// Template error check
func UnableToFindHtmlTemplate(tmpl *template.Template) {
	if tmpl == nil {
		log.Println("Unable to find html template")
		return
	}
}

// Exectuion error
func ExcutionTemplateError(err error) {
	if err != nil {
		log.Println("Error executing tempalte : ", err)
		return
	}
}

/*--------------------------------------------------------Booking Information/ Boucher/Invoice of Booking---------------------------------------------*/

func InterfaceConversion(val interface{}) string {
	iAreaId := val.(string)
	iAreaId, _ = val.(string)
	return iAreaId
}

func SessionUserId(s *Server, r *http.Request) int32 {
	uid, _ := GetSetSessionValue(s, r)
	toInt := InterfaceConversion(uid)
	intVar, err := strconv.Atoi(toInt)
	if err != nil {
		log.Println("Interface conversion error")
	}
	return int32(intVar)
}
