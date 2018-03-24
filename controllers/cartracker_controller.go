package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"cartracker.api/common"
	"cartracker.api/data"
	"cartracker.api/services"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
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

		var prefix = common.ServerCfg.SubLocation
		if len(prefix) > 0 {
			prefix = "/" + prefix
		}

		fullData := map[string]interface{}{
			"userEmail": userEmail,
			"prefix":    prefix,
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
	handlePostRequest(w, r, true)
}

func handlePostRequest(w http.ResponseWriter, r *http.Request, singleObj bool) {
	//userInfo, isAuthorized, err := services.Authorized(w, r)
	//if err != nil {
	//	fmt.Fprintln(w, "aborted")
	//	return
	//}

	if true {

		var userEmail = "" //userInfo.Email

		params := mux.Vars(r)

		body, err1 := ioutil.ReadAll(r.Body)
		if err1 != nil {
			panic(err1)
		}

		rawBody := (*json.RawMessage)(&body)
		bytesBody, err := rawBody.MarshalJSON()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v\n", string(bytesBody[:]))

		var trackerInfo common.CarTrackInfo
		var trackerInfoList []common.CarTrackInfo

		if singleObj {
			errUnmarshal := json.Unmarshal(*rawBody, &trackerInfo)
			if errUnmarshal != nil {
				panic(errUnmarshal)
			}

			trackerInfoList = append(trackerInfoList, trackerInfo)
		} else {
			errUnmarshal := json.Unmarshal(*rawBody, &trackerInfoList)
			if errUnmarshal != nil {
				panic(errUnmarshal)
			}
		}

		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var upsertedIds []bson.ObjectId

		for _, tracker := range trackerInfoList {
			var existing = common.StringInSlice(params["trackingType"], common.RequiredInfoTypes)

			if !existing || !strings.EqualFold(params["trackingType"], tracker.InfoType) {
				w.WriteHeader(422) // unprocessable entity
				var msg = "Tracking type not matching! "
				msg += " type in URL:" + params["trackingType"]
				msg += " type in Data: " + tracker.InfoType
				json.NewEncoder(w).Encode(msg)
				return
			}

			addedInfo := storeData(tracker, userEmail)

			upsertedIds = append(upsertedIds, addedInfo.UpsertedId.(bson.ObjectId))
		}

		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(upsertedIds); err != nil {
			panic(err)
		}
	}
}

func storeData(trackerInfo common.CarTrackInfo, userEmail string) *mgo.ChangeInfo {

	//fmt.Println(postedData)

	sec, dec := math.Modf(trackerInfo.TrackDate)

	var trackerEntiy = &common.CarTrackEntity{
		ActualValue:  trackerInfo.ActualValue + " " + userEmail,
		InfoType:     trackerInfo.InfoType,
		NumericValue: trackerInfo.NumericValue,
		StringValue:  trackerInfo.StringValue,
		TrackDate:    time.Unix(int64(sec), int64(dec*(1e9))),
	}

	trackerBytes, err := json.Marshal(trackerEntiy)
	if err != nil {
		panic(err)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	println(string(trackerBytes[:]))

	addedInfo, addedErr := data.AddData(trackerEntiy, &waitGroup, common.MongoSession)
	if addedErr != nil {
		panic(addedErr)
	}

	waitGroup.Wait()

	return addedInfo
}

// CreateCarTrackerInfoByBatch handles bulkly uploaded data
func CreateCarTrackerInfoByBatch(w http.ResponseWriter, r *http.Request) {
	handlePostRequest(w, r, false)
}
