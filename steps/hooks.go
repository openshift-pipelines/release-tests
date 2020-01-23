package steps

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/getgauge-contrib/gauge-go/gauge"
	. "github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"github.com/openshift-pipelines/release-tests/pkg/olm"
)

func storeToScenarioDataStore(key string, value interface{}) {
	gauge.GetScenarioStore()[key] = value
}

func storeToSpecDataStore(key string, value interface{}) {
	gauge.GetSpecStore()[key] = value
}

func storeToSuiteDataStore(key string, value interface{}) {
	gauge.GetSuiteStore()[key] = value
}

func fetchFromScenarioDataStore(key string) interface{} {
	return gauge.GetScenarioStore()[key].(interface{})
}

func fetchFromSpecDataStore(key string) interface{} {
	return gauge.GetSpecStore()[key].(interface{})
}

func fetchFromSuiteDataStore(key string) interface{} {
	return gauge.GetSuiteStore()[key].(interface{})
}

// Hooks for gauge framework
var _ = gauge.BeforeSuite(func() {

	//Creates subscription yaml with configured details from env/test/tes.properties
	helper.CreateSubscriptionYaml(config.Flags.Channel, config.Flags.InstallPlan, config.Flags.CSV)

	// subcribe to operator which we have created
	opclient, _, suitecleanup := olm.Subscribe(config.Flags.OperatorVersion)

	storeToSuiteDataStore("opclient", opclient)
	storeToSuiteDataStore("suitecleanup", suitecleanup)

	if config.Flags.TknVersion == "" {
		log.Println("env \"TKN_VERSION\" is not set cannot proceed to run tests")
		os.Exit(0)
	}

	tknBinaryPath := helper.NewTknRunner(filepath.Join(helper.RootDir(), fmt.Sprintf("../build/tkn/v%s/tkn", config.Flags.TknVersion)))
	storeToSuiteDataStore("tknPath", tknBinaryPath)
}, []string{}, AND)

var _ = gauge.BeforeScenario(func() {
	client, namespace, cleanup := helper.NewClientSet()
	storeToScenarioDataStore("client", client)
	storeToScenarioDataStore("namespace", namespace)
	storeToScenarioDataStore("cleanup", cleanup)
}, []string{}, AND)

var _ = gauge.AfterScenario(func() {
	switch c := fetchFromScenarioDataStore("cleanup").(type) {
	case func():
		defer c()
	default:
		T.Errorf("Error: return type is not of type func()")
	}
}, []string{}, AND)

var _ = gauge.AfterSuite(func() {
	switch c := fetchFromSuiteDataStore("suitecleanup").(type) {
	case func():
		defer c()
	default:
		T.Errorf("Error: return type is not of type func()")
	}
}, []string{}, AND)
