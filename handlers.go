package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// Index - home page
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// GetCarTrackerInfo getting all the records
func GetCarTrackerInfo(w http.ResponseWriter, r *http.Request) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	carTrackInfoList := GetData(nil, &waitGroup, mongoSession)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(carTrackInfoList); err != nil {
		panic(err)
	}

	waitGroup.Wait()
}

// GetCarTrackerInfoByType getting the records related to certain info type
func GetCarTrackerInfoByType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	query := bson.M{"infotype": strings.ToUpper(params["trackingType"])}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	results := GetData(query, &waitGroup, mongoSession)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(results); err != nil {
		panic(err)
	}

	waitGroup.Wait()
}

// CreateCarTrackerInfo is the function to create record with posted data from client side
func CreateCarTrackerInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var trackerInfo CarTrackInfo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.Unmarshal(body, &trackerInfo); err != nil {
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	var required = stringInSlice(params["trackingType"], requiredInfoTypes)

	if !required || !strings.EqualFold(params["trackingType"], trackerInfo.InfoType) {
		w.WriteHeader(422) // unprocessable entity
		var msg = "Tracking type not matching!"
		json.NewEncoder(w).Encode(msg)
		return
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	addedInfo, addedErr := AddData(trackerInfo, &waitGroup, mongoSession)
	if addedErr != nil {
		panic(addedErr)
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(addedInfo.UpsertedId); err != nil {
		panic(err)
	}

	waitGroup.Wait()
}
