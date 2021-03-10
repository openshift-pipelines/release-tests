package triggers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

const (
	//MaxIdleConnections specifies max connection to the http client
	MaxIdleConnections int = 30
	//RequestTimeout specifies request timeout with http client
	RequestTimeout int = 15
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

func buildHeaders(req *http.Request, interceptor, eventType string) *http.Request {
	switch strings.ToLower(interceptor) {
	case "github":
		log.Printf("Building headers for github interceptor..")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Hub-Signature", "sha1="+GetSignature(store.GetPayload(), os.Getenv("SECRET_TOKEN")))
		req.Header.Add("X-GitHub-Event", eventType)
	case "gitlab":
		log.Printf("Building headers for gitlab interceptor..")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-GitLab-Token", os.Getenv("SECRET_TOKEN"))
		req.Header.Add("X-Gitlab-Event", eventType)
	case "bitbucket":
		log.Printf("Building headers for bitbucket interceptor..")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Hub-Signature", "sha1="+GetSignature(store.GetPayload(), os.Getenv("SECRET_TOKEN")))
		req.Header.Add("X-Event-Key", "repo:"+eventType)
	default:
		testsuit.T.Errorf("Error: %s ", "Please provide valid event_listener type eg: (github, gitlab, bitbucket)")
	}
	return req
}
