package olm

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/oc"
	"github.com/openshift-pipelines/release-tests/pkg/olm"
	"github.com/openshift-pipelines/release-tests/pkg/opc"
	"github.com/openshift-pipelines/release-tests/pkg/openshift"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/pkg/statefulset"
	"github.com/openshift-pipelines/release-tests/pkg/store"
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
	if _, err := operator.EnsureTektonConfigExists(store.Clients().TektonConfig(), store.GetCRNames()); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonConfig doesn't exists\n %v", err))
	}
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

var _ = gauge.Step("Validate manual approval gate deployment", func() {
	operator.ValidateManualApprovalGateDeployments(store.Clients(), store.GetCRNames())
})

var _ = gauge.Step("Validate <deploymentName> statefulset deployment", func(deploymentName string) {
	log.Printf("Validating statefulset %v deployment\n", deploymentName)
	statefulset.ValidateStatefulSetDeployment(store.Clients(), deploymentName)
})

var _ = gauge.Step("Uninstall Operator", func() {
	// cleanup operator Traces
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

var _ = gauge.Step("Validate console plugin deployment", func() {
	rnames := store.GetCRNames()
	cs := store.Clients()

	if openshift.IsCapabilityEnabled(cs, "Console") {
		k8s.ValidateDeployments(cs, rnames.TargetNamespace, config.ConsolePluginDeployment)
	} else {
		log.Printf("OpenShift Console is not enabled, skipping validation of console plugin deployment")
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
	opc.AssertComponentVersion(defaultVersion, component)
})

var _ = gauge.Step("Check version of OSP", func() {
	defaultVersion := os.Getenv("OSP_VERSION")
	opc.AssertComponentVersion(defaultVersion, "OSP")
})

var _ = gauge.Step("Download and extract CLI from cluster", func() {
	opc.DownloadCLIFromCluster()
})

var _ = gauge.Step("Check <binary> client version", func(binary string) {
	opc.AssertClientVersion(binary)
})

var _ = gauge.Step("Check <binary> server version", func(binary string) {
	opc.AssertServerVersion(binary)
})

var _ = gauge.Step("Check <binary> version", func(binary string) {
	opc.AssertClientVersion(binary)
})

var _ = gauge.Step("Validate quickstarts", func() {
	opc.ValidateQuickstarts()
})

var _ = gauge.Step("Ensure that Tekton Results is ready", func() {
	operator.EnsureResultsReady()
})

var _ = gauge.Step("Create Results route", func() {
	operator.CreateResultsRoute()
})

var _ = gauge.Step("Verify <resourceType> Results stored", func(resourceType string) {
	operator.VerifyResultsAnnotationStored(resourceType)
})

var _ = gauge.Step("Verify <resourceType> Results records", func(resourceType string) {
	operator.VerifyResultsRecords(resourceType)
})

var _ = gauge.Step("Verify <resourceType> Results logs", func(resourceType string) {
	operator.VerifyResultsLogs(resourceType)
})

var _ = gauge.Step("Enable generateSigningSecret for Tekton Chains in TektonConfig", func() {
	patch_data := "{\"spec\":{\"chain\":{\"generateSigningSecret\":true}}}"
	if oc.SecretExists("signing-secrets", "openshift-pipelines") {
		log.Printf("Secrets \"signing-secrets\" already exists")
		if oc.GetSecretsData("signing-secrets", "openshift-pipelines") == "\"\"" {
			log.Printf("The \"signing-secrets\" does not contain any data")
			oc.UpdateTektonConfig(patch_data)
		}
	} else {
		cmd.MustSucceed("oc", "create", "secret", "generic", "signing-secrets", "--namespace", "openshift-pipelines")
		oc.UpdateTektonConfig(patch_data)
	}
})

var _ = gauge.Step("Store Cosign public key in file", func() {
	operator.CreateFileWithCosignPubKey()
})

var _ = gauge.Step("Verify <binary> version from the pipelinerun logs", func(binary string) {
	pipelines.CheckLogVersion(store.Clients(), binary, store.Namespace())
})

var _ = gauge.Step("Get olm-skip-range <upgradeType> and save to field <fieldName> in file <fileName>", func(upgradeType string, fieldName string, filename string) {
	oc.GetOlmSkipRange(upgradeType, fieldName, filename)
})

var _ = gauge.Step("Validate skipRange diff between fields <preUpgradeSkipRange> and <postUpgradeSkipRange> in file <fileName>", func(preUpgradeSkipRange string, postUpgradeSkipRange string, fileName string) {
	oc.ValidateOlmSkipRangeDiff(fileName, preUpgradeSkipRange, postUpgradeSkipRange)
})

var _ = gauge.Step("Validate OSP Version in OlmSkipRange", func() {
	oc.ValidateOlmSkipRange()
})
