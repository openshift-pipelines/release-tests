package store

import (
	"net/http"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/tkn"
)

func Namespace() string {
	return gauge.GetScenarioStore()["namespace"].(string)

}

func Clients() *clients.Clients {
	switch cs := gauge.GetScenarioStore()["clients"].(type) {
	case *clients.Clients:
		return cs
	default:
		return nil
	}
}

func HttpResponse() *http.Response {
	switch cs := gauge.GetScenarioStore()["response"].(type) {
	case *http.Response:
		return cs
	default:
		return nil
	}
}

func GetPayload() []byte {
	switch cs := gauge.GetScenarioStore()["payload"].(type) {
	case []byte:
		return cs
	default:
		return nil
	}
}

func Tkn() tkn.Cmd {
	switch n := gauge.GetSuiteStore()["tkn"].(type) {
	case tkn.Cmd:
		return n
	default:
		panic("Error: type for tkn is not as expected")
	}
}

func PutScenarioData(key, value string) {
	gauge.GetScenarioStore()[key] = value
}

func GetScenarioData(key string) string {
	return gauge.GetScenarioStore()[key].(string)
}
