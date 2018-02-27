package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/david-x-chen/cartracker.api/common"
	"github.com/david-x-chen/cartracker.api/data"
	"github.com/david-x-chen/cartracker.api/services"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// Index - home page
func Index(w http.ResponseWriter, r *http.Request) {
	userInfo, isAuthorized, err := services.Authorized(w, r)
	if err != nil {
		fmt.Fprintln(w, "aborted")
		return
	}

	if isAuthorized {

		var userEmail = userInfo.Email

		//common.PushStaticResource(w, "/static/css/style.css")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		fullData := map[string]interface{}{
			"userEmail": userEmail,
			"prefix":    common.ServerCfg.SubLocation,
		}
		fmt.Print(userEmail)
		common.RenderTemplate(w, r, common.Tmpls["home.html"], "base", fullData)
	}
}

// GetCarTrackerInfo getting all the records
func GetCarTrackerInfo(w http.ResponseWriter, r *http.Request) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	carTrackInfoList := data.GetData(nil, &waitGroup, common.MongoSession)

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

	results := data.GetData(query, &waitGroup, common.MongoSession)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(results); err != nil {
		panic(err)
	}

	waitGroup.Wait()
}

// CreateCarTrackerInfo is the function to create record with posted data from client side
func CreateCarTrackerInfo(w http.ResponseWriter, r *http.Request) {
	userInfo, isAuthorized, err := services.Authorized(w, r)
	if err != nil {
		fmt.Fprintln(w, "aborted")
		return
	}

	if isAuthorized {

		var userEmail = userInfo.Email

		params := mux.Vars(r)

		var trackerInfo common.CarTrackInfo
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

		var required = common.StringInSlice(params["trackingType"], common.RequiredInfoTypes)

		if !required || !strings.EqualFold(params["trackingType"], trackerInfo.InfoType) {
			w.WriteHeader(422) // unprocessable entity
			var msg = "Tracking type not matching!"
			json.NewEncoder(w).Encode(msg)
			return
		}

		var waitGroup sync.WaitGroup
		waitGroup.Add(1)

		trackerInfo.ActualValue += ";" + userEmail

		addedInfo, addedErr := data.AddData(trackerInfo, &waitGroup, common.MongoSession)
		if addedErr != nil {
			panic(addedErr)
		}

		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(addedInfo.UpsertedId); err != nil {
			panic(err)
		}

		waitGroup.Wait()
	}
}
