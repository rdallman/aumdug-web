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

func getAllEvents(limit int) (events Events, err error) {
	//TODO clean this up
	s := getSession()
	defer s.Close()
	iter := s.DB(dbName).C("events").Find(nil).Limit(100).Iter()
	err = iter.All(&events)
	return
}

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
