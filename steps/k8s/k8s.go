package k8s

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Create cron job with schedule <schedule>", func(schedule string) {
	args := []string{"wget", "--spider", store.GetScenarioData("route")}
	k8s.CreateCronJob(store.Clients(), args, schedule, store.Namespace())
})

var _ = gauge.Step("Wait for cron job to be active", func() {
	k8s.WaitForCronJobToBeSceduled(store.Clients(), 1, store.GetScenarioData("cronjob"), store.Namespace())
})

var _ = gauge.Step("Delete cron job", func() {
	k8s.DeleteCronJob(store.Clients(), store.GetScenarioData("cronjob"), store.Namespace())
})
