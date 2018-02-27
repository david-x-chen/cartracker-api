package routers

import (
	"fmt"
	"net/http"

	"github.com/david-x-chen/cartracker.api/common"

	"github.com/david-x-chen/cartracker.api/logger"
	"github.com/gorilla/mux"
)

// NewRouter setting up the router
func NewRouter() *mux.Router {
	var loc = common.ServerCfg.SubLocation
	router := mux.NewRouter().StrictSlash(true)
	if len(loc) > 0 {
		router = router.PathPrefix(fmt.Sprintf("/%s/", loc)).Subrouter()
	}
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	return router
}
