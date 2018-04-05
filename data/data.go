package data

import (
	"log"
	"sync"
	"time"

	"cartracker.api/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// InitMongoSession initialises session
func InitMongoSession() *mgo.Session {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{common.MongoConfig.MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: common.MongoConfig.AuthDatabase,
		//Username: common.MongoConfig.AuthUserName,
		//Password: common.MongoConfig.AuthPassword,
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
func GetData(query bson.M, waitGroup *sync.WaitGroup, mongoSession *mgo.Session) *[]common.CarTrackEntity {
	// Decrement the wait group count so the program knows this
	// has been completed once the goroutine exits.
	defer waitGroup.Done()

	// Request a socket connection from the session to process our query.
	// Close the session when the goroutine exits and put the connection back
	// into the pool.
	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	// Get a collection to execute the query against.
	collection := sessionCopy.DB(common.MongoConfig.TestDatabase).C("obd2info")

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

	carTrackEntityList := new([]common.CarTrackEntity)
	// Retrieve the list of track information.
	err = collection.Find(query).All(carTrackEntityList)
	if err != nil {
		log.Printf("RunQuery: ERROR: %s\n", err)
		return nil
	}

	log.Printf("RunQuery: Count[%d]\n", len(*carTrackEntityList))

	return carTrackEntityList
}

// AddData is the function to insert record to db
func AddData(info *common.CarTrackEntity, waitGroup *sync.WaitGroup, mongoSession *mgo.Session) (*mgo.ChangeInfo, error) {
	defer waitGroup.Done()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	// Get a collection to execute the query against.
	collection := sessionCopy.DB(common.MongoConfig.TestDatabase).C("obd2info")

	// Index
	index := mgo.Index{
		Key:        []string{"_id", "infotype", "trackdate"},
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
	collection := sessionCopy.DB(common.MongoConfig.TestDatabase).C("obd2info")

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
