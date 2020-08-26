package triggers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/openshift-pipelines/release-tests/pkg/assert"
)

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 10
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

// GetSignature is a HMAC sha1 generator
func GetSignature(input []byte, key string) string {
	keyForSign := []byte(key)
	h := hmac.New(sha1.New, keyForSign)
	_, err := h.Write(input)
	assert.NoError(err, "Couldn't generate signature")
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
