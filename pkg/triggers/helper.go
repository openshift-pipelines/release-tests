package triggers

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	resource "github.com/openshift-pipelines/release-tests/pkg/config"
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

// CreateHTTPSClient for connection re-use
func CreateHTTPSClient() *http.Client {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(resource.Path("testdata/triggers/certs/server.crt"), resource.Path("testdata/triggers/certs/server.key"))
	if err != nil {
		log.Fatal(err)
	}
	caCert, err := os.ReadFile(resource.Path("testdata/triggers/certs/ca.crt"))
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      caCertPool,
			},
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
	if err != nil {
		testsuit.T.Errorf("could not generate signature \n %v", err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func buildHeaders(req *http.Request, interceptor, eventType string) *http.Request {
	switch strings.ToLower(interceptor) {
	case "github":
		log.Printf("Building headers for github interceptor..")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Hub-Signature", "sha1="+GetSignature(store.GetPayload(), config.TriggersSecretToken))
		req.Header.Add("X-GitHub-Event", eventType)
	case "gitlab":
		log.Printf("Building headers for gitlab interceptor..")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-GitLab-Token", config.TriggersSecretToken)
		req.Header.Add("X-Gitlab-Event", eventType)
	case "bitbucket":
		log.Printf("Building headers for bitbucket interceptor..")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Hub-Signature", "sha1="+GetSignature(store.GetPayload(), config.TriggersSecretToken))
		req.Header.Add("X-Event-Key", "repo:"+eventType)
	default:
		testsuit.T.Errorf("Error: %s ", "Please provide valid event_listener type eg: (github, gitlab, bitbucket)")
	}
	return req
}
