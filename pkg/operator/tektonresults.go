package operator

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"encoding/base64"
	"encoding/json"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
)

func CreateSecretsForTektonResults() {
	var password = cmd.MustSucceed("openssl", "rand", "-base64", "20").Stdout()
	password = strings.ReplaceAll(password, "\n", "")
	cmd.MustSucceed("oc", "create", "secret", "-n", "openshift-pipelines", "generic", "tekton-results-postgres", "--from-literal=POSTGRES_USER=result", "--from-literal=POSTGRES_PASSWORD="+password)
	// generating tls certificate
	cmd.MustSucceed("openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", "key.pem", "-out", "cert.pem", "-days", "365", "-nodes", "-subj", "/CN=tekton-results-api-service.openshift-pipelines.svc.cluster.local", "-addext", "subjectAltName=DNS:tekton-results-api-service.openshift-pipelines.svc.cluster.local")
	// creating secret with generated certificate
	cmd.MustSucceed("oc", "create", "secret", "tls", "-n", "openshift-pipelines", "tekton-results-tls", "--cert=cert.pem", "--key=key.pem")
}

func EnsureResultsReady() {
	cmd.MustSuccedIncreasedTimeout(time.Minute*5, "oc", "wait", "--for=condition=Ready", "tektoninstallerset", "-l", "operator.tekton.dev/type=result", "--timeout=120s")
}

func CreateResultsRoute() {
	cmd.Run("oc", "create", "route", "-n", "openshift-pipelines", "passthrough", "tekton-results-api-service", "--service=tekton-results-api-service", "--port=8080")
}

func GetResultsApi() string {
	var results_api = cmd.MustSucceed("oc", "get", "route", "tekton-results-api-service", "-n", "openshift-pipelines", "--no-headers", "-o", "custom-columns=:spec.host").Stdout() + ":443"
	results_api = strings.ReplaceAll(results_api, "\n", "")
	return results_api
}

func GetResultsAnnotations(resourceType string) (string, string, string) {
	var result_uuid = cmd.MustSucceed("tkn", resourceType, "describe", "--last", "-o", "jsonpath='{.metadata.annotations.results\\.tekton\\.dev/result}'").Stdout()
	var record_uuid = cmd.MustSucceed("tkn", resourceType, "describe", "--last", "-o", "jsonpath='{.metadata.annotations.results\\.tekton\\.dev/record}'").Stdout()
	var stored = cmd.MustSucceed("tkn", resourceType, "describe", "--last", "-o", "jsonpath='{.metadata.annotations.results\\.tekton\\.dev/stored}'").Stdout()
	record_uuid = strings.ReplaceAll(record_uuid, "'", "")
	result_uuid = strings.ReplaceAll(result_uuid, "'", "")
	stored = strings.ReplaceAll(stored, "'", "")
	return result_uuid, record_uuid, stored
}

func VerifyResultsStored(resourceType string) {
	_, _, storedAnnotation := GetResultsAnnotations(resourceType)

	if storedAnnotation == "" {
		testsuit.T.Fail(fmt.Errorf("annotation results.tekton.dev/stored is not set"))
	}

	stored, err := strconv.ParseBool(storedAnnotation)

	if err != nil {
		testsuit.T.Fail(fmt.Errorf("annotation results.tekton.dev/record doesn't contain a boolean value"))
	}
	if !stored {
		testsuit.T.Fail(fmt.Errorf("annotation results.tekton.dev/record is set to false"))
	}
}

func VerifyResultsLogs(resourceType string) {
	var record_uuid string
	var results_api string
	_, record_uuid, _ = GetResultsAnnotations(resourceType)
	results_api = GetResultsApi()

	if record_uuid == "" {
		testsuit.T.Fail(fmt.Errorf("annotation results.tekton.dev/record is not set"))
	}

	var resultsJsonData = cmd.MustSucceed("opc", "results", "logs", "get", "--insecure", "--addr", results_api, record_uuid).Stdout()
	if strings.Contains(resultsJsonData, "record not found") {
		testsuit.T.Errorf("Results log not found")
	} else {
		type ResultLogs struct {
			Name string `json:"name"`
			Data string `json:"data"`
		}
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
	_, record_uuid, _ = GetResultsAnnotations(resourceType)
	results_api = GetResultsApi()
	var results_record = cmd.MustSucceed("opc", "results", "records", "get", "--insecure", "--addr", results_api, record_uuid).Stdout()
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
		if !strings.Contains(string(decodedResultsLogs), "Hello, Results!") || !strings.Contains(string(decodedResultsLogs), "Goodbye, Results!") {
			testsuit.T.Errorf("Records are incorrect")
		}
	}
}
