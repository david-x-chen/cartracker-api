package routers

import (
	"net/http"

	"github.com/david-x-chen/cartracker.api/common"

	"github.com/david-x-chen/cartracker.api/logger"
	"github.com/gorilla/mux"
)

// NewRouter setting up the router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		var loc = common.ServerCfg.SubLocation
		var pattern = route.Pattern

		if len(loc) > 0 {
			pattern = "/" + loc + pattern
		}

		router.
			Methods(route.Method).
			Path(pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	return router
}
