package cli

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/chains"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Create Chains CR with format <format>, storage <storage> and transparency enabled <enabled>", func(format, storage, transparencyEnabled string) {
	chains.CreateChainsCR(store.Clients(), format, storage, transparencyEnabled)
})
