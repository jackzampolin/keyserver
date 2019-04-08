package api

import (
	"encoding/json"
	"net/http"
)

// VersionHandler handles the /version route
func (s *Server) VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write(s.newVersion().marshal())
}

type version struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Branch  string `json:"branch"`
}

func (s *Server) newVersion() version {
	return version{s.Version, s.Commit, s.Branch}
}

func (v version) marshal() []byte {
	out, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return out
}
