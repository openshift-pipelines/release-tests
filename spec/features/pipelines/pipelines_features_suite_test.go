package pipelines

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	flags "github.com/openshift-pipelines/release-tests/spec/flags"
)

var tkn helper.TknRunner

func TestSuiteFeatures(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pipelines E2E TestSuite")
	//RunSpecsWithDefaultAndCustomReporters(t, "Pipelines e2e TestSuite", []Reporter{reporter.JunitReport(t, "../../../reports")})
}

var _ = BeforeSuite(func() {
	_, _, flags.CleanupSuite = operator.InstallOperator(flags.OperatorVersion)
	tkn = helper.NewTknRunner(filepath.Join(helper.RootDir(), fmt.Sprintf("../build/tkn/v%s/tkn", os.Getenv("TKN_VERSION"))))
})

var _ = AfterSuite(func() {
	defer flags.CleanupSuite()
})
