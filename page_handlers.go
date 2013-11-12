package main

import (
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func handleHome(rw http.ResponseWriter, req *http.Request) {
	events, err := getAllEvents(100)
	if err != nil {
		//TODO yeah make this a warning or something
		log.Println("Could not get events")
		http.Error(rw, "500 Internal Server Error", 500)
	}
	err = templates.ExecuteTemplate(rw, "home.html", events)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
