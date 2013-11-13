package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handleApiAllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := getAllEvents(100)

	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	writeJson(w, events)
}

func handleApiEventCreate(w http.ResponseWriter, r *http.Request) {
	//return success or empty event?
	data := struct {
		Success bool  `json:"success"`
		Event   Event `json:"event"`
	}{
		Success: false,
	}

	//TODO auth?
	//request := struct {
	//Key   string `json:"key"`
	//Event Event  `json:"event"`
	//}{}

	var e Event

	if readJson(r, &e) /*&& auth(request.Key) */ {
		if err := createEvent(&e); err != nil {
			log.Printf("%v\n", err)
		} else {
			data.Success = true
			data.Event = e
		}
	}

	writeJson(w, data)
}

func handleApiEventDestroy(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Success bool `json:"success"`
	}{
		Success: false,
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if err := destroyEvent(id); err != nil {
		log.Printf("%v", err)
	} else {
		data.Success = true
	}

	writeJson(w, data)
}

func handleApiEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	e, err := getEvent(id)

	if err != nil {
		log.Printf("%v", err)
	}

	writeJson(w, e)
}

func handleApiEventUpdate(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Success bool  `json:"success"`
		Event   Event `json:"event"`
	}{
		Success: false,
	}
	var e Event
	//TODO auth?
	if readJson(r, &e) {
		if err := updateEvent(&e); err != nil {
			log.Printf("%v", err)
		} else {
			data.Success = true
			data.Event = e
		}
	}

	writeJson(w, data)
}
