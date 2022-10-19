package steps

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/gauge_messages"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/oc"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Runs Before every Secenario
var _ = gauge.BeforeScenario(func(exInfo *gauge_messages.ExecutionInfo) {
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
	store["targetNamespace"] = config.TargetNamespace

	autoPruneTagPresent := true
	for _, tag := range exInfo.CurrentScenario.Tags {
		if tag == "auto-prune" {
			autoPruneTagPresent = false
		}
	}
	if !autoPruneTagPresent {
		oc.AnnotateNamespace(namespace, "operator.tekton.dev/prune.skip=true")
	}
	oc.Create("testdata/pvc.yaml", namespace)
}, []string{}, testsuit.AND)

// Runs After every Secenario
var _ = gauge.AfterScenario(func(exInfo *gauge_messages.ExecutionInfo) {
	switch c := gauge.GetScenarioStore()["scenario.cleanup"].(type) {
	case func():
		if exInfo.CurrentSpec.IsFailed {
			log.Printf("Skipping deletion of namespace %s as the test got failed", store.Namespace())
		} else {
			c()
		}
	default:
		testsuit.T.Errorf("Error: return type is not of type func()")
	}
}, []string{}, testsuit.AND)

// Store default pruner config
var _ = gauge.BeforeSpec(func(exInfo *gauge_messages.ExecutionInfo) {
	cs, _, _ := k8s.NewClientSet()

	tc, err := cs.TektonConfig().Get(context.TODO(), config.TektonConfigName, metav1.GetOptions{})
	if err != nil {
		testsuit.T.Errorf("Error: could not get TektonConfig")
	}

	store := gauge.GetSpecStore()
	if tc.Spec.Pruner.Keep != nil {
		store["keep"] = strconv.FormatUint(uint64(*(tc.Spec.Pruner.Keep)), 10)
	} else {
		store["keep"] = "null"
	}
	if tc.Spec.Pruner.KeepSince != nil {
		store["keepSince"] = strconv.FormatUint(uint64(*(tc.Spec.Pruner.KeepSince)), 10)
	} else {
		store["keepSince"] = "null"
	}
	store["resources"] = tc.Spec.Pruner.Resources
	store["schedule"] = tc.Spec.Pruner.Schedule
}, []string{"auto-prune"}, testsuit.AND)

// Revert changes made by pruner tests
var _ = gauge.AfterSpec(func(exInfo *gauge_messages.ExecutionInfo) {
	keep := gauge.GetSpecStore()["keep"]
	keepSince := gauge.GetSpecStore()["keepSince"]
	resources := gauge.GetSpecStore()["resources"].([]string)
	schedule := gauge.GetSpecStore()["schedule"]

	resourcesList := strings.Join(resources, "\",\"")
	patch_data := fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep\":%s,\"keep-since\":%s,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", keep, keepSince, schedule, resourcesList)

	oc.UpdateTektonConfig(patch_data)
}, []string{"auto-prune"}, testsuit.AND)
