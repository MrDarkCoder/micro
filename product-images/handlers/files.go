package handlers

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/MrDarkCoder/productimages/files"
	"github.com/gorilla/mux"
)

// Files is a handler for reading and writing files
type Files struct {
	log   *log.Logger
	store files.Storage
}

// NewFiles creates a new File handler
func NewFiles(s files.Storage, l *log.Logger) *Files {
	return &Files{store: s, log: l}
}

// ServeHTTP implements the http.Handler interface
func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Println("Handle POST", "id", id, "filename", fn)
	log.Println("Handle POST", "id", id, "filename", fn)

	// no need to check for invalid id or filename as the mux router will not send requests
	// here unless they have the correct parameters

	f.saveFile(id, fn, rw, r)
}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.log.Println("Invalid path", "path", uri)
	http.Error(rw, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}

// saveFile saves the contents of the request to a file
func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r *http.Request) {
	f.log.Println("Save file for product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r.Body)
	if err != nil {
		f.log.Println("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
