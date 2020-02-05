package steps

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
)

var _ = gauge.Step("Verify ServiceAccount <sa> does not exist", func(sa string) {
	k8s.VerifyNoServiceAccount(Clients().KubeClient, sa, Namespace())
})
