package server

import (
	"fmt"
	"net/http"

	configuration "github.com/khezen/espipe/configuration"
	"github.com/khezen/espipe/indexer"
)

const endpoint = ":5000"

// Server - Contains data required for serving web REST requests
type Server struct {
	indexer indexer.Indexer
	quit    chan error
}

// New - Create new service for serving web REST requests
func New(config configuration.Configuration, quit chan error) (*Server, error) {
	i, err := indexer.New(config)
	if err != nil {
		return nil, err
	}
	return &Server{
		*i,
		quit,
	}, nil
}

// ListenAndServe - Blocks the current goroutine, opens an HTTP port and serves the web REST requests
func (s *Server) ListenAndServe() {
	http.HandleFunc("/espipe/health/", s.handleHealthCheck)
	http.HandleFunc("/espipe/", s.handleIndexDoc)
	fmt.Printf("opening espipe at %v\n", endpoint)
	s.quit <- http.ListenAndServe(endpoint, nil)
}

// GET /espipe/health
func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
