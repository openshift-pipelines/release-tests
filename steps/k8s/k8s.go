package k8s

import (
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

var _ = gauge.Step("Create cron job with schedule <schedule>", func(schedule string) {
	args := []string{"curl", "-X", "POST", "--data", "{}", store.GetScenarioData("route")}
	k8s.CreateCronJob(store.Clients(), args, schedule, store.Namespace())
})

var _ = gauge.Step("Wait for cron job to be active", func() {
	k8s.WaitForCronJobToBeSceduled(store.Clients(), 1, store.GetScenarioData("cronjob"), store.Namespace())
})

var _ = gauge.Step("Delete cron job", func() {
	k8s.DeleteCronJob(store.Clients(), store.GetScenarioData("cronjob"), store.Namespace())
})

var _ = gauge.Step("Validate default auto prune cronjob in target namespace", func() {
	namespace := store.TargetNamespace()
	k8s.AssertIfDefaultCronjobExists(store.Clients(), namespace)
})
