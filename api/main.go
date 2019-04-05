package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Server represents the API server
type Server struct {
	Port int `json:"port"`

	Version string
	Commit  string
	Branch  string
}

// Router returns the router
func (s *Server) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/version", s.VersionHandler)

	return router
}

// VersionHandler handles the /version route
func (s *Server) VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("{\"version\": \"%s\", \"commit\": \"%s\", \"branch\": \"%s\"}", s.Version, s.Commit, s.Branch)))
}
