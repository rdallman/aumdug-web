package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/howeyc/fsnotify"
	"html/template"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
	"time"
)

type (
	Events []Event
	Event  struct {
		Id          bson.ObjectId `bson:"_id,omitempty" json:"-"`
		Name        string        `bson:"Name"          json:"name"`
		Description string        `bson:"Description"   json:"description"`
		Date        time.Time     `bson:"Date"          json:"date"`
		Created     time.Time     `bson:"Created"       json:"created"`
	}
	//TODO make an interface out of this
	//TODO figure out sessions in order to make above happen
	//eventRepo struct {
	//Collection *mgo.Collection
	//}
)

var (
	mgoSession *mgo.Session
	dbName     = "aumdug"
	templates  = template.Must(template.ParseGlob("templates/*.html"))
)

func getSession() *mgo.Session {
	if mgoSession == nil {
		session, err := mgo.Dial("localhost")
		if err != nil {
			fmt.Println(err)
		}
		//for real go?
		mgoSession = session
	}
	return mgoSession.Clone()
}

func withCollection(c string) *mgo.Collection {
	session := getSession()
	return session.DB(dbName).C(c)
}

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

func handleAllEvents(rw http.ResponseWriter, req *http.Request) {
	events, err := getAllEvents(100)

	if err != nil {
		log.Println(err)
		http.Error(rw, "500 Internal Server Error", 500)
		return
	}

	writeJson(rw, events)
}

func getAllEvents(limit int) (events Events, err error) {
	iter := withCollection("events").Find(nil).Limit(100).Iter()
	err = iter.All(&events)
	return
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

func createEvent(e *Event) (err error) {
	if e.Id.Hex() == "" {
		e.Id = bson.NewObjectId()
	}
	if e.Created.IsZero() {
		e.Created = time.Now()
	}
	_, err = withCollection("events").UpsertId(e.Id, e)
	return
}

func handleHome(rw http.ResponseWriter, req *http.Request) {
	events, err := getAllEvents(100)
	if err != nil {
		//TODO yeah make this a warning or something
		fmt.Fprintf(rw, "Could not get events")
	}
	err = templates.ExecuteTemplate(rw, "home.html", events)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
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
	r.HandleFunc("/events", handleAllEvents).Methods("GET")
	r.HandleFunc("/events", handleEventCreate).Methods("PUT")
	//r.HandleFunc("/events/{id}", EventHandler).Methods("GET")

	r.HandleFunc("/", handleHome)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
