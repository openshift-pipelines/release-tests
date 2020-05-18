package steps

import (
	"log"
	"os"
	"strings"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
)

// Runs Before every Secenario
var _ = gauge.BeforeScenario(func() {
	cs, namespace, cleanup := k8s.NewClientSet()
	store := gauge.GetScenarioStore()
	store["clients"] = cs
	store["namespace"] = namespace
	store["scenario.cleanup"] = cleanup
	if strings.Contains(strings.ToLower(os.Getenv("OPERATOR_ENV")), "stag") {
		log.Printf("Running on (stage) environment : %s", os.Getenv("OPERATOR_ENV"))
		cmd.MustSucceed("oc", "adm", "policy", "add-cluster-role-to-user", "system:image-puller", "-z", "pipeline", "-n", namespace)
	}
}, []string{}, testsuit.AND)

// Runs After every Secenario
var _ = gauge.AfterScenario(func() {
	switch c := gauge.GetScenarioStore()["scenario.cleanup"].(type) {
	case func():
		c()
	default:
		testsuit.T.Errorf("Error: return type is not of type func()")
	}
}, []string{}, testsuit.AND)
