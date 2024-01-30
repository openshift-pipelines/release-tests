package operator

import (
	"fmt"
	"strings"
	"time"

	"encoding/base64"
	"encoding/json"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
)

func CreateSecretsForTektonResults() {
	var password string = cmd.MustSucceed("openssl", "rand", "-base64", "20").Stdout()
	password = strings.Replace(password, "\n", "", -1)
	cmd.MustSucceed("oc", "create", "secret", "-n", "openshift-pipelines", "generic", "tekton-results-postgres", "--from-literal=POSTGRES_USER=result", "--from-literal=POSTGRES_PASSWORD="+password)
	//generating tls certifiacte
	cmd.MustSucceed("openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", "key.pem", "-out", "cert.pem", "-days", "365", "-nodes", "-subj", "/CN=tekton-results-api-service.openshift-pipelines.svc.cluster.local", "-addext", "subjectAltName=DNS:tekton-results-api-service.openshift-pipelines.svc.cluster.local")
	//creating secret with generated certificate
	cmd.MustSucceed("oc", "create", "secret", "tls", "-n", "openshift-pipelines", "tekton-results-tls", "--cert=cert.pem", "--key=key.pem")
}

func EnsureResutsReady() {
	cmd.MustSuccedIncreasedTimeout(time.Minute*5, "oc", "wait", "--for=condition=Ready", "tektoninstallerset", "-l", "operator.tekton.dev/type=result")
}

func CreateResultsRoute() {
	cmd.Run("oc", "create", "route", "-n", "openshift-pipelines", "passthrough", "tekton-results-api-service", "--service=tekton-results-api-service", "--port=8080")
}

func GetResultsApi() string {
	var results_api string = cmd.MustSucceed("oc", "get", "route", "tekton-results-api-service", "-n", "openshift-pipelines", "--no-headers", "-o", "custom-columns=:spec.host").Stdout() + ":443"
	results_api = strings.ReplaceAll(results_api, "\n", "")
	return results_api
}

func GetResultsAnnotations(resourceType string) (string, string) {
	var log_uuid string = cmd.MustSucceed("tkn", resourceType, "describe", "--last", "-o", "jsonpath='{.metadata.annotations.results\\.tekton\\.dev/log}'").Stdout()
	var record_uuid string = cmd.MustSucceed("tkn", resourceType, "describe", "--last", "-o", "jsonpath='{.metadata.annotations.results\\.tekton\\.dev/record}'").Stdout()
	record_uuid = strings.ReplaceAll(record_uuid, "'", "")
	log_uuid = strings.ReplaceAll(log_uuid, "'", "")
	return log_uuid, record_uuid
}

func VerifyResultsLogs(resourceType string) {
	var log_uuid string
	var results_api string
	log_uuid, _ = GetResultsAnnotations(resourceType)
	results_api = GetResultsApi()
	var results_log string = cmd.MustSucceed("opc", "results", "logs", "get", "--insecure", "--addr", results_api, log_uuid).Stdout()
	if strings.Contains(results_log, "record not found") {
		testsuit.T.Errorf("Results log not found")
	} else {
		type ResultLogs struct {
			Name string `json:"name"`
			Data string `json:"data"`
		}
		resultsJsonData := cmd.MustSucceed("opc", "results", "logs", "get", "--insecure", "--addr", results_api, log_uuid).Stdout()
		var resultLogs ResultLogs
		err := json.Unmarshal([]byte(resultsJsonData), &resultLogs)
		if err != nil {
			testsuit.T.Errorf("Error parsing JSON")
		}
		decodedResultsLogs, err := base64.StdEncoding.Strict().DecodeString(resultLogs.Data)
		if err != nil {
			testsuit.T.Errorf("Error decoding base64 data")
		}
		if !strings.Contains(string(decodedResultsLogs), "Hello, Results!") || !strings.Contains(string(decodedResultsLogs), "Goodbye, Results!") {
			testsuit.T.Errorf("Logs are incorrect")
		}
	}
}

func VerifyResultsRecords(resourceType string) {
	var record_uuid string
	var results_api string
	_, record_uuid = GetResultsAnnotations(resourceType)
	results_api = GetResultsApi()
	var results_record string = cmd.MustSucceed("opc", "results", "records", "get", "--insecure", "--addr", results_api, record_uuid).Stdout()
	if strings.Contains(results_record, "record not found") {
		testsuit.T.Errorf("Results record not found")
	} else {
		type ResultRecords struct {
			Data struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"data"`
		}
		resultsJsonData := cmd.MustSucceed("opc", "results", "records", "get", "--insecure", "--addr", results_api, record_uuid, "-o", "json").Stdout()
		var resultRecords ResultRecords
		err := json.Unmarshal([]byte(resultsJsonData), &resultRecords)
		if err != nil {
			testsuit.T.Errorf("Error parsing JSON: %v", err)
		}
		decodedResultsLogs, err := base64.StdEncoding.Strict().DecodeString(resultRecords.Data.Value)
		if err != nil {
			testsuit.T.Errorf("Error decoding base64 data: %v", err)
		}
		fmt.Printf(string(decodedResultsLogs))
		if !strings.Contains(string(decodedResultsLogs), "Hello, Results!") || !strings.Contains(string(decodedResultsLogs), "Goodbye, Results!") {
			testsuit.T.Errorf("Records are incorrect")
		}
	}
}
