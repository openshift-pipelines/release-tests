package steps

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"github.com/openshift-pipelines/release-tests/pkg/olm"
)

// Hooks for gauge framework
var _ = gauge.BeforeSuite(func() {

	//Creates subscription yaml with configured details from env/test/tes.properties
	helper.CreateSubscriptionYaml(config.Flags.Channel, config.Flags.InstallPlan, config.Flags.CSV)

	// subcribe to operator which we have created
	opclient, _, suitecleanup := olm.Subscribe(config.Flags.OperatorVersion)

	if config.Flags.TknVersion == "" {
		log.Println("env \"TKN_VERSION\" is not set cannot proceed to run tests")
		os.Exit(0)
	}

	tknBinaryPath := helper.NewTknRunner(filepath.Join(helper.RootDir(), fmt.Sprintf("../build/tkn/v%s/tkn", config.Flags.TknVersion)))
	store := gauge.GetSuiteStore()
	store["opclient"] = opclient
	store["suitecleanup"] = suitecleanup
	store["tknPath"] = tknBinaryPath

}, []string{}, testsuit.AND)

var _ = gauge.BeforeScenario(func() {
	client, namespace, cleanup := helper.NewClientSet()

	store := gauge.GetScenarioStore()
	store["client"] = client
	store["namespace"] = namespace
	store["cleanup"] = cleanup

}, []string{}, testsuit.AND)

var _ = gauge.AfterScenario(func() {

	switch c := gauge.GetScenarioStore()["cleanup"].(type) {
	case func():
		c()
	default:
		testsuit.T.Errorf("Error: return type is not of type func()")
	}
}, []string{}, testsuit.AND)

var _ = gauge.AfterSuite(func() {
	switch c := gauge.GetSuiteStore()["suitecleanup"].(type) {
	case func():
		c()
	default:
		testsuit.T.Errorf("Error: return type is not of type func()")
	}
}, []string{}, testsuit.AND)
