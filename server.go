package main

import (
	"github.com/gorilla/mux"
	"github.com/howeyc/fsnotify"
	"html/template"
	"log"
	"net/http"
)

//reload templates on modify
func listenForChanges() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)
				//recompile templates
				templates = template.Must(template.ParseGlob("templates/*.html"))
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch("templates/")
	if err != nil {
		log.Fatal(err)
	}

	<-done

	watcher.Close()
}

func main() {
	go listenForChanges()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	r := mux.NewRouter()
	r.HandleFunc("/api/events", handleApiAllEvents).Methods("GET")
	r.HandleFunc("/api/events", handleApiEventCreate).Methods("PUT")
	r.HandleFunc("/api/events/{id}", handleApiEventDestroy).Methods("DELETE")
	r.HandleFunc("/api/events/{id}", handleApiEvent).Methods("GET")

	r.HandleFunc("/", handleHome)
	r.HandleFunc("/events/{name}", handleEvent)
	r.HandleFunc("/api", handleApi)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
