package cli

import (
	"fmt"
	"log"
	"strings"

	"github.com/getgauge-contrib/gauge-go/gauge"
	m "github.com/getgauge-contrib/gauge-go/models"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/oc"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Create <table>", func(table *m.Table) {
	for _, row := range table.Rows {
		resource := row.Cells[1]
		oc.Create(resource, store.Namespace())
	}
})

var _ = gauge.Step("Enable TLS config for eventlisteners", func() {
	oc.EnableTLSConfigForEventlisteners(store.Namespace())

})

var _ = gauge.Step("Verify kubernetes events for eventlistener", func() {
	oc.VerifyKubernetesEventsForEventListener(store.Namespace())
})

var _ = gauge.Step("Delete <table>", func(table *m.Table) {
	for _, row := range table.Rows {
		resource := row.Cells[1]
		oc.Delete(resource, store.Namespace())
	}
})

var _ = gauge.Step("Create & Link secret <secret> to service account <sa>", func(secret, sa string) {
	oc.CreateSecretWithSecretToken(secret, store.Namespace())
	oc.LinkSecretToSA(secret, sa, store.Namespace())
})

var _ = gauge.Step("Update pruner config <keepPresence> keep <keep> schedule <schedule> resouces <resources> and <keepSincePresence> keep-since <keepSince>", func(keepPresence, keep, schedule, resouces, keepSincePresence, keepSince string) {
	resoucesSplit := strings.Split(resouces, ",")
	resourcesList := strings.Join(resoucesSplit, "\",\"")
	patch_data := ""
	if keepPresence == "with" && keepSincePresence == "without" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep\":%s,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", keep, schedule, resourcesList)
	} else if keepPresence == "without" && keepSincePresence == "with" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep-since\":%s,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", keepSince, schedule, resourcesList)
	} else if keepPresence == "with" && keepSincePresence == "with" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep\":%s,\"keep-since\":%s,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", keep, keepSince, schedule, resourcesList)
	} else if keepPresence == "without" && keepSincePresence == "without" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", schedule, resourcesList)
	}
	oc.UpdateTektonConfig(patch_data)
})

var _ = gauge.Step("Update pruner config with invalid data <keepPresence> keep <keep> schedule <schedule> resouces <resources> and <keepSincePresence> keep-since <keepSince> and expect error message <errorMessage>", func(keepPresence, keep, schedule, resouces, keepSincePresence, keepSince, errorMessage string) {
	resoucesSplit := strings.Split(resouces, ",")
	resourcesList := strings.Join(resoucesSplit, "\",\"")
	patch_data := ""
	if keepPresence == "with" && keepSincePresence == "without" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep\":%s,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", keep, schedule, resourcesList)
	} else if keepPresence == "without" && keepSincePresence == "with" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep-since\":%s,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", keepSince, schedule, resourcesList)
	} else if keepPresence == "with" && keepSincePresence == "with" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep\":%s,\"keep-since\":%s,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", keep, keepSince, schedule, resourcesList)
	} else if keepPresence == "without" && keepSincePresence == "without" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", schedule, resourcesList)
	}
	oc.UpdateTektonConfigwithInvalidData(patch_data, errorMessage)
})

var _ = gauge.Step("Remove auto pruner configuration from config CR", func() {
	log.Print("Removing pruner configuration from config CR")
	oc.RemovePrunerConfig()
})

var _ = gauge.Step("Assert if cronjob with prefix <cronJobName> is <status> in target namespace", func(cronJobName, status string) {
	namespace := store.TargetNamespace()
	log.Printf("Verifying if the cronjob %v is %v in namespace %v", cronJobName, status, namespace)
	oc.VerifyCronjobStatus(cronJobName, status, namespace)
})

var _ = gauge.Step("Annotate namespace with <annotation>", func(annotation string) {
	log.Printf("Annotating namespace %v with %v", store.Namespace(), annotation)
	oc.AnnotateNamespace(store.Namespace(), annotation)
})

var _ = gauge.Step("Remove annotation <annotation> from namespace", func(annotation string) {
	log.Printf("Removing annotation %v from namespace %v", store.Namespace(), annotation)
	oc.AnnotateNamespace(store.Namespace(), annotation+"-")
})

var _ = gauge.Step("Add label <label> to namespace", func(label string) {
	log.Printf("Labelling namespace %v with %v", store.Namespace(), label)
	oc.LabelNamespace(store.Namespace(), label)
})

var _ = gauge.Step("Remove label <label> from the namespace", func(label string) {
	log.Printf("Removing annotation %v from namespace %v", store.Namespace(), label)
	oc.AnnotateNamespace(store.Namespace(), label+"-")
})

var _ = gauge.Step("<cts> clustertasks are <status>", func(cts, status string) {
	if cts == "community" {
		cts = config.CommunityClustertasks
	}
	log.Printf("Checking if clustertasks %v is %v", cts, status)
	ctsList := strings.Split(cts, ",")
	for _, c := range ctsList {
		pipelines.GetClusterTask(store.Clients(), c, status)
	}
})