package main

import (
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func handleHome(w http.ResponseWriter, r *http.Request) {
	events, err := getAllEvents(5)
	if err != nil {
		//TODO yeah make this a warning or something
		log.Println("Could not get events")
		http.Error(w, "500 Internal Server Error", 500)
	}
	err = templates.ExecuteTemplate(w, "home.html", events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleApi(w http.ResponseWriter, r *http.Request) {
	md, err := ioutil.ReadFile("static/api.md")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	p := struct {
		Body template.HTML
	}{
		Body: template.HTML(blackfriday.MarkdownCommon(md)),
	}

	err = templates.ExecuteTemplate(w, "base.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//TODO idk if we even need this
func handleEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n := vars["name"]
	_, err := getEventByName(n)

	if err != nil {
		log.Printf("%v", err)
	}

}
