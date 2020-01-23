package steps

import (
	"sync"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/client"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
)

var once sync.Once
var _ = gauge.Step("Operator should be installed", func() {
	once.Do(func() {
		operator.ValidateOperatorInstall(GetOperatorClient())
	})
})

func GetNameSpace() string {
	return fetchFromScenarioDataStore("namespace").(string)
}

func GetClient() *client.Clients {
	switch c := fetchFromScenarioDataStore("client").(type) {
	case *client.Clients:
		return c
	default:
		return nil
	}
}

func GetOperatorClient() *client.Clients {
	switch c := fetchFromSuiteDataStore("opclient").(type) {
	case *client.Clients:
		return c
	default:
		return nil
	}
}

func GetTknBinaryPath() helper.TknRunner {
	switch n := fetchFromSuiteDataStore("tknPath").(type) {
	case helper.TknRunner:
		return n
	default:
		panic("Error")
	}
}
