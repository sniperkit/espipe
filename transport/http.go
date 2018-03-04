package transport

import (
	"net/http"
	"time"
)

var (
	httpTransport = &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	httpClient = &http.Client{
		Transport: httpTransport,
	}
)

// HTTPClient return the http client
func HTTPClient() *http.Client {
	return httpClient
}
