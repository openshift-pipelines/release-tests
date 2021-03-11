package steps

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
)

// Runs Before every Secenario
var _ = gauge.BeforeScenario(func() {
	cs, namespace, cleanup := k8s.NewClientSet()
	crNames := config.ResourceNames{
		TektonPipeline:  "pipeline",
		TektonTrigger:   "trigger",
		TektonAddon:     "addon",
		TektonConfig:    config.TektonConfigName,
		Namespace:       "",
		TargetNamespace: config.TargetNamespace,
	}

	store := gauge.GetScenarioStore()
	store["crnames"] = crNames
	store["clients"] = cs
	store["namespace"] = namespace
	store["scenario.cleanup"] = cleanup
}, []string{}, testsuit.AND)

// Runs After every Secenario
var _ = gauge.AfterScenario(func() {
	//switch c := gauge.GetScenarioStore()["scenario.cleanup"].(type) {
	//case func():
	//	c()
	//default:
	//	testsuit.T.Errorf("Error: return type is not of type func()")
	//}
}, []string{}, testsuit.AND)
