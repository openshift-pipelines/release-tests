package cli

import (
	"fmt"
	"log"
	"strings"

	"github.com/getgauge-contrib/gauge-go/gauge"
	m "github.com/getgauge-contrib/gauge-go/models"
	"github.com/openshift-pipelines/release-tests/pkg/oc"
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

var _ = gauge.Step("Update pruner config with keep <keep> schedule <schedule> resouces <resources>", func(keep, schedule, resouces string) {
	resouces_split := strings.Split(resouces, ",")
	resources_list := strings.Join(resouces_split, "\",\"")
	patch_data := fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep\":%s,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", keep, schedule, resources_list)
	oc.UpdateTektonConfig(patch_data)
})

var _ = gauge.Step("Update pruner config to default", func() {
	log.Print("Updating pruner config to default value")
	patch_data := "{\"spec\":{\"pruner\":{\"keep\":1,\"schedule\":\"\"}}}"
	oc.UpdateTektonConfig(patch_data)
})

var _ = gauge.Step("Assert if cronjob <cronJobName> is <status> in namespace <namespace>", func(cronJobName, status, namespace string) {
	log.Printf("Verifying if the cronjob %v is %v in namespace %v", cronJobName, status, namespace)
	oc.VerifyCronjobStatus(cronJobName, status, namespace)
})
