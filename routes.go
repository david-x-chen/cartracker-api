package main

import (
	"net/http"
)

// Route type
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is list
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"GetCarTracker",
		"GET",
		"/cartracker",
		GetCarTrackerInfo,
	},
	Route{
		"GetCarTrackerInfoByType",
		"GET",
		"/cartracker/{trackingType}",
		GetCarTrackerInfoByType,
	},
	Route{
		"CreateCarTrackerInfo",
		"POST",
		"/cartracker/{trackingType}",
		CreateCarTrackerInfo,
	},
}
