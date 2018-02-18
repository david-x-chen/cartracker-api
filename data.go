package main

import (
	"io/ioutil"
	"log"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
)

type (
	config struct {
		// MongoDBHosts host server
		MongoDBHosts string `yaml:"mongodbhosts"` //"localhost"
		// AuthDatabase database
		AuthDatabase string `yaml:"authdatabase"` //"cartracker"
		// AuthUserName user name of mongodb
		AuthUserName string `yaml:"authusername"`
		// AuthPassword password of mongodb user
		AuthPassword string `yaml:"authpassword"`
		// TestDatabase test db to connect
		TestDatabase string `yaml:"testdatabase"` //"cartracker"
	}

	// CarTrackInfo Car track information
	CarTrackInfo struct {
		TrackDate    time.Time `json:"trackdate,omitempty"`
		InfoType     string    `json:"infotype,omitempty"`
		StringValue  string    `json:"stringvalue,omitempty"`
		NumericValue float32   `json:"numericvalue,omitempty"`
		ActualValue  string    `json:"actualvalue,omitempty"`
	}
)

var requiredInfoTypes = []string{
	"RPM", "SPEED", "STATUS", "ENGINE_LOAD", "SHORT_FUEL_TRIM_1", "LONG_FUEL_TRIM_1",
	"THROTTLE_POS", "COMMANDED_EQUIV_RATIO", "MAF", "INTAKE_TEMP", "COOLANT_TEMP",
	"CONTROL_MODULE_VOLTAGE", "TIMING_ADVANCE", "RUN_TIME"}

func (c *config) getConfig() *config {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

var mongoConfig config

//
func connect() *mgo.Session {
	mongoConfig.getConfig()

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{mongoConfig.MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: mongoConfig.AuthDatabase,
		//Username: mongoConfig.AuthUserName,
		//Password: mongoConfig.AuthPassword,
	}

	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic(err)
	}

	mongoSession.SetMode(mgo.Monotonic, true)

	//defer mongoSession.Close()

	return mongoSession
}

// GetData is a function that is launched as a goroutine to perform the MongoDB work.
func GetData(query bson.M, waitGroup *sync.WaitGroup, mongoSession *mgo.Session) *[]CarTrackInfo {
	// Decrement the wait group count so the program knows this
	// has been completed once the goroutine exits.
	defer waitGroup.Done()

	// Request a socket connection from the session to process our query.
	// Close the session when the goroutine exits and put the connection back
	// into the pool.
	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	// Get a collection to execute the query against.
	collection := sessionCopy.DB(mongoConfig.TestDatabase).C("obd2info")

	// Index
	index := mgo.Index{
		Key:        []string{"infotype", "trackdate"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	carTrackInfoList := new([]CarTrackInfo)
	// Retrieve the list of track information.
	err = collection.Find(query).All(carTrackInfoList)
	if err != nil {
		log.Printf("RunQuery: ERROR: %s\n", err)
		return nil
	}

	log.Printf("RunQuery: Count[%d]\n", len(*carTrackInfoList))

	return carTrackInfoList
}

// AddData is the function to insert record to db
func AddData(info CarTrackInfo, waitGroup *sync.WaitGroup, mongoSession *mgo.Session) (*mgo.ChangeInfo, error) {
	defer waitGroup.Done()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	// Get a collection to execute the query against.
	collection := sessionCopy.DB(mongoConfig.TestDatabase).C("obd2info")

	// Index
	index := mgo.Index{
		Key:        []string{"infotype", "trackdate"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	var ID = bson.NewObjectId()

	info.TrackDate = time.Now()

	upsertInfo, err := collection.Upsert(bson.M{"_id": ID}, info)
	if err != nil {
		panic(err)
	}

	return upsertInfo, err
}

// RemoveAllData is the function to remove all records from db
func RemoveAllData(query bson.M, waitGroup *sync.WaitGroup, mongoSession *mgo.Session) (*mgo.ChangeInfo, error) {
	defer waitGroup.Done()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	// Get a collection to execute the query against.
	collection := sessionCopy.DB(mongoConfig.TestDatabase).C("obd2info")

	// Index
	index := mgo.Index{
		Key:        []string{"infotype", "trackdate"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	info, errRemove := collection.RemoveAll(query)
	if errRemove != nil {
		panic(errRemove)
	}

	print(info)

	return info, errRemove
}
