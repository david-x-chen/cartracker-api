package main

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

var mongoSession *mgo.Session

// our main function
func main() {
	mongoSession = connect()

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8000", router))
}
