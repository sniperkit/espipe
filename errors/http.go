package errors

import (
	"errors"
)

var (
	// ErrPathNotFound - 404
	ErrPathNotFound = errors.New("ErrPathNotFound - The request path is not supported")
	// ErrWrongMethod - 405
	ErrWrongMethod = errors.New("ErrWrongMethod - The request http method does not match expectation")
)

// HTTPStatusCode -
func HTTPStatusCode(err error) int {
	switch err {
	case ErrPathNotFound:
		return 404
	case ErrWrongMethod:
		return 405
	default:
		return 500
	}
}
