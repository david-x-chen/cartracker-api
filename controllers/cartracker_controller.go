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

		errUnmarshal := json.Unmarshal(*rawBody, &trackerInfo)
		if errUnmarshal != nil {
			panic(errUnmarshal)
		}

		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		postedData := fmt.Sprintf("%v", trackerInfo)
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

		var existing = common.StringInSlice(params["trackingType"], common.RequiredInfoTypes)

		if !existing || !strings.EqualFold(params["trackingType"], trackerEntiy.InfoType) {
			w.WriteHeader(422) // unprocessable entity
			var msg = "Tracking type not matching! "
			msg += " type in URL:" + params["trackingType"]
			msg += " type in Data: " + trackerEntiy.InfoType
			msg += " posted data: " + postedData
			msg += " converted data: " + string(trackerBytes[:])
			json.NewEncoder(w).Encode(msg)
			return
		}

		var waitGroup sync.WaitGroup
		waitGroup.Add(1)

		println(string(trackerBytes[:]))

		addedInfo, addedErr := data.AddData(trackerEntiy, &waitGroup, common.MongoSession)
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
