package install

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift-pipelines/release-tests/pkg/olm"
	flags "github.com/openshift-pipelines/release-tests/spec/flags"
)

func TestSuiteFeatures(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Olm E2E TestSuite")
	//RunSpecsWithDefaultAndCustomReporters(t, "olm e2e TestSuite", []Reporter{reporter.JunitReport(t, "../../../reports")})
}

var _ = BeforeSuite(func() {
	flags.Clients, _, flags.CleanupSuite = olm.Subscribe(flags.OperatorVersion)
})

var _ = AfterSuite(func() {
	defer flags.CleanupSuite()
})
