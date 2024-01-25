package olm

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/oc"
	"github.com/openshift-pipelines/release-tests/pkg/olm"
	"github.com/openshift-pipelines/release-tests/pkg/openshift"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/openshift-pipelines/release-tests/pkg/tkn"
)

var once sync.Once
var _ = gauge.Step("Validate Operator should be installed", func() {
	once.Do(func() {
		operator.ValidateOperatorInstallStatus(store.Clients(), store.GetCRNames())
	})
})

var _ = gauge.Step("Subscribe to operator", func() {
	// Creates subscription yaml with configured details from env/test/test.properties
	if _, err := olm.SubscribeAndWaitForOperatorToBeReady(store.Clients(), config.Flags.SubscriptionName, config.Flags.Channel, config.Flags.CatalogSource); err != nil {
		testsuit.T.Fail(fmt.Errorf("operator not ready after creating subscription \n %v", err))
	}
})

var _ = gauge.Step("Wait for TektonConfig CR availability", func() {
	operator.EnsureTektonConfigExists(store.Clients().TektonConfig(), store.GetCRNames())
})

var _ = gauge.Step("Upgrade operator subscription", func() {
	// Creates subscription yaml with configured details from env/test/test.properties
	if _, err := olm.UptadeSubscriptionAndWaitForOperatorToBeReady(store.Clients(), config.Flags.SubscriptionName, config.Flags.Channel); err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to update subscription \n %v", err))
	}
})

var _ = gauge.Step("Validate RBAC", func() {
	operator.ValidateRBAC(store.Clients(), store.GetCRNames())
})

var _ = gauge.Step("Validate pipelines deployment", func() {
	operator.ValidatePipelineDeployments(store.Clients(), store.GetCRNames())
})

var _ = gauge.Step("Validate triggers deployment", func() {
	operator.ValidateTriggerDeployments(store.Clients(), store.GetCRNames())
})

var _ = gauge.Step("Validate chains deployment", func() {
	operator.ValidateChainsDeployments(store.Clients(), store.GetCRNames())
})

var _ = gauge.Step("Validate hub deployment", func() {
	operator.ValidateHubDeployments(store.Clients(), store.GetCRNames())
})

var _ = gauge.Step("Uninstall Operator", func() {
	//cleanup operator Traces
	operator.Uninstall(store.Clients(), store.GetCRNames())
})

var _ = gauge.Step("Verify TektonAddons Install status", func() {
	operator.EnsureTektonAddonsStatusInstalled(store.Clients().TektonAddon(), store.GetCRNames())
})

var _ = gauge.Step("Validate PAC deployment", func() {
	rnames := store.GetCRNames()
	cs := store.Clients()
	k8s.ValidateDeployments(cs, rnames.TargetNamespace, config.PacControllerName)
	k8s.ValidateDeployments(cs, rnames.TargetNamespace, config.PacWatcherName)
	k8s.ValidateDeployments(cs, rnames.TargetNamespace, config.PacWebhookName)
})

var _ = gauge.Step("Validate tkn server cli deployment", func() {
	rnames := store.GetCRNames()
	cs := store.Clients()

	if openshift.IsCapabilityEnabled(cs, "Console") {
		k8s.ValidateDeployments(cs, rnames.TargetNamespace, config.TknDeployment)
	} else {
		log.Printf("OpenShift Console is not enabled, skipping validation of tkn serve CLI deployment")
	}
})

var _ = gauge.Step("Validate tektoninstallersets status", func() {
	k8s.ValidateTektonInstallersetStatus(store.Clients())
})

var _ = gauge.Step("Validate tektoninstallersets names", func() {
	k8s.ValidateTektonInstallersetNames(store.Clients())
})

var _ = gauge.Step("Check version of component <component>", func(component string) {
	defaultVersion := os.Getenv(strings.ToUpper(component + "_version"))
	tkn.AssertComponentVersion(defaultVersion, component)
})

var _ = gauge.Step("Check version of OSP", func() {
	defaultVersion := os.Getenv("OSP_VERSION")
	tkn.AssertComponentVersion(defaultVersion, "OSP")
})

var _ = gauge.Step("Download and extract CLI from cluster", func() {
	tkn.DownloadCLIFromCluster()
})

var _ = gauge.Step("Check <binary> client version", func(binary string) {
	tkn.AssertClientVersion(binary)
})

var _ = gauge.Step("Check <binary> version", func(binary string) {
	tkn.AssertClientVersion(binary)
})

var _ = gauge.Step("Validate quickstarts", func() {
	tkn.ValidateQuickstarts()
})

var _ = gauge.Step("Create secrets for Tekton Results", func() {
	if !oc.SecretExists("tekton-results-postgres", "openshift-pipelines") && !oc.SecretExists("tekton-results-tls", "openshift-pipelines"){
		operator.CreateSecretsForTektonResults()
	} else {
		log.Printf("\"tekton-results-postgres\" or \"tekton-results-tls\" secrets already exist")
	}
})

var _ = gauge.Step("Ensure that Tekton Results is ready", func() {
	operator.EnsureResutsReady()
})

var _ = gauge.Step("Create Results route", func(){
	operator.CreateResultsRoute()
})

var _ = gauge.Step("Verify <resourceType> Results records", func(resourceType string){
	operator.VerifyResultsRecords(resourceType)
})

var _ = gauge.Step("Verify <resourceType> Results logs", func(resourceType string){
	operator.VerifyResultsLogs(resourceType)
})