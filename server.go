package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/howeyc/fsnotify"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func writeJson(w http.ResponseWriter, v interface{}) {
	//to avoid json vulnerabilities, wrap v in object literal
	doc := map[string]interface{}{"d": v}

	if data, err := json.Marshal(doc); err != nil {
		log.Printf("Error marshalling json: %v\n", err)
	} else {
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func readJson(r *http.Request, v interface{}) bool {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Could not read request body %v\n", err)
		return false
	}

	if err = json.Unmarshal(body, v); err != nil {
		log.Printf("Coult not parse request body %v\n", err)
		return false
	}

	return true
}

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

	r := mux.NewRouter()
	r.HandleFunc("/api/events", handleAllEvents).Methods("GET")
	r.HandleFunc("/api/events", handleEventCreate).Methods("PUT")
	//r.HandleFunc("/api/events/{id}", EventHandler).Methods("GET")

	r.HandleFunc("/", handleHome)

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.ListenAndServe(":8080", nil)
}
