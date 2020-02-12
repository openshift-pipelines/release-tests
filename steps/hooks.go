package steps

import (
	"log"
	"os"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/olm"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
)

// Hooks for gauge framework
var _ = gauge.BeforeSuite(func() {

	// TODO:
	//err := flags.Parse()
	//err := tools.Validate()

	// TODO: validate all required tools are present

	// Creates subscription yaml with configured details from env/test/test.properties
	olm.CreateSubscriptionYaml(config.Flags.Channel, config.Flags.InstallPlan, config.Flags.CSV)

	// subcribe to operator which we have created
	olm.Subscribe()

	if config.Flags.TknVersion == "" {
		log.Println("env \"TKN_VERSION\" is not set cannot proceed to run tests")
		os.Exit(1)
	}

	// TODO add tkn to store
	//tknPath := config.File("..", "build", "tkn", "v"+config.Flags.TknVersion, "tkn")
	//if _, err := os.Stat(tknPath); os.IsNotExist(err) {
	//log.Printf("tkn cli not found in at %q ", tknPath)
	//os.Exit(1)
	//}
	//tknCmd := tkn.New(tknPath)
	//store["tkn"] = tknCmd

	// TODO: fix how store is used
	store := gauge.GetSuiteStore()
	store["suite.cleanup"] = func() { olm.Unsubscribe() }

}, []string{}, testsuit.AND)

var _ = gauge.AfterSuite(func() {
	switch c := gauge.GetSuiteStore()["suite.cleanup"].(type) {
	case func():
		c()
	default:
		testsuit.T.Errorf("Error: return type is not of type func()")
	}
	// TODO: fix flags magic
	//cleanup operator Traces
	operator.Cleanup(config.Flags.OperatorVersion)

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
