package store

import (
	"net/http"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/opc"
	"github.com/tektoncd/operator/test/utils"
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

func GetCRNames() utils.ResourceNames {
	switch names := gauge.GetScenarioStore()["crnames"].(type) {
	case utils.ResourceNames:
		return names
	default:
		return utils.ResourceNames{}
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

func Opc() opc.Cmd {
	switch n := gauge.GetSuiteStore()["opc"].(type) {
	case opc.Cmd:
		return n
	default:
		panic("Error: type for opc is not as expected")
	}
}

func PutScenarioData(key, value string) {
	gauge.GetScenarioStore()[key] = value
}

func PutScenarioDataSlice(key string, value []string) {
	gauge.GetScenarioStore()[key] = value
}

func GetScenarioDataSlice(key string) []string {
	return gauge.GetScenarioStore()[key].([]string)
}
func GetScenarioData(key string) string {
	return gauge.GetScenarioStore()[key].(string)
}

func TargetNamespace() string {
	return gauge.GetScenarioStore()["targetNamespace"].(string)
}
