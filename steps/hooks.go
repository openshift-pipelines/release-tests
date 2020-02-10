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
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/olm"
	"github.com/openshift-pipelines/release-tests/pkg/tkn"
)

// Hooks for gauge framework
var _ = gauge.BeforeSuite(func() {

	// TODO:
	//err := flags.Parse()
	//err := tools.Validate()

	// TODO: validate all required tools are present

	// Creates subscription yaml with configured details from env/test/tes.properties
	olm.CreateSubscriptionYaml(config.Flags.Channel, config.Flags.InstallPlan, config.Flags.CSV)

	// subcribe to operator which we have created
	// TODO: fix flags magic
	opVersion := config.Flags.OperatorVersion
	olm.Subscribe(opVersion)

	if config.Flags.TknVersion == "" {
		log.Println("env \"TKN_VERSION\" is not set cannot proceed to run tests")
		os.Exit(0)
	}

	tknCmd := tkn.New(filepath.Join(
		helper.RootDir(),
		fmt.Sprintf("../build/tkn/v%s/tkn", config.Flags.TknVersion)))

	/// TODO: fix how store is used
	store := gauge.GetSuiteStore()
	store["suite.cleanup"] = func() { olm.Unsubscribe(opVersion) }
	store["tkn"] = tknCmd

}, []string{}, testsuit.AND)

var _ = gauge.AfterSuite(func() {
	switch c := gauge.GetSuiteStore()["suite.cleanup"].(type) {
	case func():
		c()
	default:
		testsuit.T.Errorf("Error: return type is not of type func()")
	}
}, []string{}, testsuit.AND)

var _ = gauge.BeforeScenario(func() {
	cs, namespace, cleanup := k8s.NewClientSet()

	store := gauge.GetScenarioStore()
	store["clients"] = cs
	store["namespace"] = namespace
	store["scenario.cleanup"] = cleanup

}, []string{}, testsuit.AND)

var _ = gauge.AfterScenario(func() {

	switch c := gauge.GetScenarioStore()["scenario.cleanup"].(type) {
	case func():
		c()
	default:
		testsuit.T.Errorf("Error: return type is not of type func()")
	}
}, []string{}, testsuit.AND)
