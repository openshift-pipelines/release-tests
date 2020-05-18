package triggers

import (
	"net/http"
	"time"
)

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 5
)

// CreateHTTPClient for connection re-use
func CreateHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}

	return client
}
