package server

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/khezen/espipe/document"
	"github.com/khezen/espipe/errors"
	"github.com/khezen/espipe/template"
)

// POST /espipe/{template}/{type}
func (s *Server) handleIndexDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.serveError(w, r, errors.ErrWrongMethod)
	}
	urlSplit := strings.Split(strings.Trim(strings.ToLower(r.URL.Path), "/"), "/")
	if len(urlSplit) != 3 {
		s.serveError(w, r, errors.ErrPathNotFound)
		return
	}
	docTemplate := template.Name(urlSplit[1])
	docType := document.Type(urlSplit[2])
	docBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.serveError(w, r, err)
		return
	}
	err = s.indexer.Index(docTemplate, docType, docBytes)
	if err != nil {
		s.serveError(w, r, err)
		return
	}
	w.Header().Set("Connection", "close")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
