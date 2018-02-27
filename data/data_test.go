package data

import (
	"strings"
	"sync"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func TestGetConfig(t *testing.T) {
	var c DbConfig
	c.getConfig()

	if strings.Contains(c.AuthDatabase, "cartrack") {
		t.Logf("LOG: MongoDB Hosts %s", c.MongoDBHosts)
	}
}

func TestGetData(t *testing.T) {
	mongoSession := connect()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	trackinfoList := GetData(bson.M{"infotype": "SPEED"}, &waitGroup, mongoSession)

	if len(*trackinfoList) > 0 {
		t.Logf("LOG: track information count: %d\n", len(*trackinfoList))
	}

	waitGroup.Wait()

	defer mongoSession.Close()
}

func TestAddData(t *testing.T) {
	mongoSession := connect()

	var waitGroup sync.WaitGroup
	waitGroup.Add(3)

	trackInfo := CarTrackInfo{
		TrackDate:    time.Now(),
		InfoType:     "TEST",
		StringValue:  "TEST",
		NumericValue: 100,
		ActualValue:  "TEST",
	}

	addInfo, err := AddData(trackInfo, &waitGroup, mongoSession)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("LOG: Added one test data. info: %v %v\n", trackInfo, addInfo)

	var id = addInfo.UpsertedId

	trackinfoList := GetData(bson.M{"infotype": "TEST"}, &waitGroup, mongoSession)
	if len(*trackinfoList) > 0 {
		t.Logf("LOG: track information count: %d\n", len(*trackinfoList))
	}

	info, err := RemoveAllData(bson.M{"_id": id}, &waitGroup, mongoSession)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("LOG: Removed all test data. info: %v\n", info)

	waitGroup.Wait()

	defer mongoSession.Close()
}
