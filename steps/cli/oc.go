package cli

import (
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
