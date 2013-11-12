package main

import (
	"log"
	"net/http"
)

func handleAllEvents(rw http.ResponseWriter, req *http.Request) {
	events, err := getAllEvents(100)

	if err != nil {
		log.Println(err)
		http.Error(rw, "500 Internal Server Error", 500)
		return
	}

	writeJson(rw, events)
}

func handleEventCreate(rw http.ResponseWriter, req *http.Request) {
	var e Event
	data := struct {
		Success bool  `json:"success"`
		Event   Event `json:"event"`
	}{
		Success: false,
	}
	if readJson(req, &e) {
		if err := createEvent(&e); err != nil {
			log.Printf("%v\n", err)
		} else {
			data.Success = true
			data.Event = e
		}
	}

	writeJson(rw, data)
}
