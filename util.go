package main

import (
	"encoding/json"
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

func auth(key string) bool {
	//TODO make this env variable
	if key == "partytimeexcellent" {
		return true
	}
	return false
}
