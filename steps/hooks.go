package steps

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	"github.com/openshift-pipelines/release-tests/pkg/tkn"
)

// Hooks for gauge framework
var _ = gauge.BeforeSuite(func() {

	// Creates subscription yaml with configured details from env/test/tes.properties
	operator.CreateSubscriptionYaml(config.Flags.Channel, config.Flags.InstallPlan, config.Flags.CSV)

	// subcribe to operator which we have created
	opclient, _, suitecleanup := operator.Subscribe(config.Flags.OperatorVersion)

	if config.Flags.TknVersion == "" {
		log.Println("env \"TKN_VERSION\" is not set cannot proceed to run tests")
		os.Exit(0)
	}

	tknCmd := tkn.New(filepath.Join(
		helper.RootDir(),
		fmt.Sprintf("../build/tkn/v%s/tkn", config.Flags.TknVersion)))

	store := gauge.GetSuiteStore()
	store["opclient"] = opclient
	store["suitecleanup"] = suitecleanup
	store["tkn"] = tknCmd

}, []string{}, testsuit.AND)

var _ = gauge.BeforeScenario(func() {
	cs, namespace, cleanup := clients.NewClients()

	store := gauge.GetScenarioStore()
	store["clients"] = cs
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
