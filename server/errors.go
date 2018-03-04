package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/khezen/espipe/errors"
)

func (s *Server) serveError(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Printf("%v", err.Error())
	w.Header().Set("Connection", "close")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	statusCode := errors.HTTPStatusCode(err)
	w.WriteHeader(statusCode)
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, err.Error())
}
