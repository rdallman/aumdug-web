package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type (
	Events []Event
	Event  struct {
		Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
		Name        string        `bson:"Name"          json:"name"`
		Description string        `bson:"Description"   json:"description"`
		Date        time.Time     `bson:"Date"          json:"date"`
		Created     time.Time     `bson:"Created"       json:"c"`
		//Tags        []string      `bson:"Tags"          json:"tags"`
	}
	//TODO make an interface out of this
	//TODO figure out sessions in order to make above happen
	//TODO eliminate s.DB.C in all queries
	//eventRepo struct {
	//Collection *mgo.Collection
	//}
)

var (
	mgoSession *mgo.Session
	dbName     = "aumdug"
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

// GET /api/events
func getAllEvents(limit int) (events Events, err error) {
	s := getSession()
	defer s.Close()
	err = s.DB(dbName).C("events").Find(nil).Limit(limit).All(&events)
	return
}

// PUT /api/events
func createEvent(e *Event) (err error) {
	if e.Id.Hex() == "" {
		e.Id = bson.NewObjectId()
	}
	if e.Created.IsZero() {
		e.Created = time.Now()
	}

	s := getSession()
	defer s.Close()
	_, err = s.DB(dbName).C("events").UpsertId(e.Id, e)
	return
}

// GET /api/events/{id}
func getEvent(id string) (e Event, err error) {
	bid := bson.ObjectIdHex(id)

	s := getSession()
	defer s.Close()
	err = s.DB(dbName).C("events").FindId(bid).One(&e)
	fmt.Printf("%v", e.Id)
	return
}

// for web
//TODO may be unnecessary
func getEventByName(name string) (e Event, err error) {

	s := getSession()
	defer s.Close()
	err = s.DB(dbName).C("events").Find(bson.M{"Name": name}).One(&e)
	return
}

// PUT /api/events/{id}
func updateEvent(e *Event) (err error) {
	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": bson.M{
				"Name":        e.Name,
				"Description": e.Description,
				"Date":        e.Date,
			}}}

	s := getSession()
	defer s.Close()
	_, err = s.DB(dbName).C("events").FindId(e.Id).Apply(change, e)

	return
}

// DELETE /api/events/{id}
func destroyEvent(id string) (err error) {
	bid := bson.ObjectIdHex(id)

	s := getSession()
	defer s.Close()
	err = s.DB(dbName).C("events").RemoveId(bid)
	return
}
