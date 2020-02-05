package steps

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
)

var _ = gauge.Step("Verify Service Account <sa> does not exist", func(sa string) {
	k8s.VerifyNoServiceAccount(GetClient(), sa, GetNameSpace())
})
