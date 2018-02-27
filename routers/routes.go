package routers

import (
	"net/http"

	"github.com/david-x-chen/cartracker.api/controllers"
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
		controllers.Index,
	},
	Route{
		"FileServer",
		"GET",
		"/static/{cache_id}/{filename}",
		controllers.FileServer,
	},
	Route{
		"Authorize",
		"GET",
		"/authorize",
		controllers.Authorize,
	},
	Route{
		"OAuth2Callback",
		"GET",
		"/oauth2callback",
		controllers.OAuth2Callback,
	},
	Route{
		"GetCarTracker",
		"GET",
		"/cartracker",
		controllers.GetCarTrackerInfo,
	},
	Route{
		"GetCarTrackerInfoByType",
		"GET",
		"/cartracker/{trackingType}",
		controllers.GetCarTrackerInfoByType,
	},
	Route{
		"CreateCarTrackerInfo",
		"POST",
		"/cartracker/{trackingType}",
		controllers.CreateCarTrackerInfo,
	},
}
