package k8s

import (
	"log"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Verify ServiceAccount <sa> does not exist", func(sa string) {
	k8s.VerifyNoServiceAccount(store.Clients().Ctx, store.Clients().KubeClient, sa, store.Namespace())
})

var _ = gauge.Step("Verify ServiceAccount <sa> exist", func(sa string) {
	k8s.VerifyServiceAccountExists(store.Clients().Ctx, store.Clients().KubeClient, sa, store.Namespace())
})

var _ = gauge.Step("Verify namespace <ns> exists", func(ns string) {
	k8s.VerifyNamespaceExists(store.Clients().Ctx, store.Clients().KubeClient, ns)
})

var _ = gauge.Step("Create cron job with schedule <schedule>", func(schedule string) {
	args := []string{"curl", "-X", "POST", "--data", "{}", store.GetScenarioData("route")}
	k8s.CreateCronJob(store.Clients(), args, schedule, store.Namespace())
})

var _ = gauge.Step("Delete cron job", func() {
	if err := k8s.DeleteCronJob(store.Clients(), store.GetScenarioData("cronjob"), store.Namespace()); err != nil {
		log.Printf("Delete cron job failed\n %v", err)
	}
})

var _ = gauge.Step("Validate default auto prune cronjob in target namespace", func() {
	namespace := store.TargetNamespace()
	k8s.AssertIfDefaultCronjobExists(store.Clients(), namespace)
})

var _ = gauge.Step("Store name of the cronjob in target namespace with schedule <schedule> to variable <variableName>", func(schedule, variable string) {
	namespace := store.TargetNamespace()
	cronJobName := k8s.GetCronjobNameWithSchedule(store.Clients(), namespace, schedule)
	store.PutScenarioData(variable, cronJobName)
})

var _ = gauge.Step("Assert pruner cronjob(s) in namespace <namespace> contains <num> number of container(s)", func(namespace, num string) {
	if namespace == "target namespace" {
		namespace = store.TargetNamespace()
	}
	k8s.AssertPrunerCronjobWithContainer(store.Clients(), namespace, num)
})

var _ = gauge.Step("Assert if cronjob with prefix <cronJobName> is <status> in target namespace", func(cronJobName, status string) {
	namespace := store.TargetNamespace()
	log.Printf("Verifying if the cronjob %v is %v in namespace %v", cronJobName, status, namespace)
	if status == "present" {
		k8s.AssertCronjobPresent(store.Clients(), cronJobName, namespace)
	} else {
		k8s.AssertCronjobNotPresent(store.Clients(), cronJobName, namespace)
	}
})
