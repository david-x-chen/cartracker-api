package controllers

import (
	"log"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

// FileServer serving static files
func FileServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["cache_id"]

	// Logging for the example
	log.Println(id)

	// Check if the id is valid, else 404, 301 to the new URL, etc - goes here!
	// (this is where you'd look up the SHA-1 hash)

	// Assuming it's valid
	file := vars["filename"]

	// Logging for the example
	log.Println(file)

	// Super simple. Doesn't set any cache headers, check existence, avoid race conditions, etc.
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(file)))
	http.ServeFile(w, r, "./static/"+file)
}
