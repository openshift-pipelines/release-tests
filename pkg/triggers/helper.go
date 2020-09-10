package triggers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/store"
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
	return hex.EncodeToString(h.Sum(nil))
}

func buildHeaders(req *http.Request, payload string) *http.Request {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(payload), &result)
	assert.FailOnError(err)
	for key, value := range result {
		if key == "X-GitLab-Token" {
			req.Header.Add(key, os.Getenv("SECRET_TOKEN"))
		} else if key == "X-Hub-Signature" {
			req.Header.Add(key, "sha1="+GetSignature(store.GetPayload(), os.Getenv("SECRET_TOKEN")))
		} else {
			req.Header.Add(key, value.(string))
		}
	}
	return req
}
