package triggers

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"io/ioutil"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

const (
	//MaxIdleConnections specifies max connection to the http client
	MaxIdleConnections int = 20
	//RequestTimeout specifies request timeout with http client
	RequestTimeout int = 10
)

// CreateHTTPClient for connection re-use
func CreateHTTPClient() *http.Client {
	//e1 := os.Setenv("GODEBUG", "x509ignoreCN=0")
	//fmt.Println("set env error", e1)
	//confFile, e1 := ioutil.ReadFile("/etc/pki/tls/openssl.cnf")
	//fmt.Println("openssl configuration", string(confFile), "***********e1", e1)
	//fmt.Println("gogdebug valuesrae", os.Getenv("GODEBUG"))
	gopath := os.Getenv("GOPATH")
	// Load client cert
	cert, err := tls.LoadX509KeyPair(gopath+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs/tls.crt",
		gopath+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs/tls.key")
	if err != nil {
		log.Fatal(err)
	}
	d, e := ioutil.ReadFile(gopath+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs/rootCA.crt")
	fmt.Println("the read file is", string(d), "**********************", e)
	caCert, err := ioutil.ReadFile(gopath+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs/rootCA.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
			TLSClientConfig: &tls.Config{
				//InsecureSkipVerify:          true,
				//ServerName: "tls.test.apps.savita47new.tekton.codereadyqe.com",
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
