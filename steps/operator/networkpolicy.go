package olm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/oc"
)

type componentNPConfig struct {
	policies   []string
	deployment string
}

var networkPoliciesByComponent = map[string]componentNPConfig{
	"pipelines": {
		deployment: "tekton-pipelines-controller",
		policies: []string{
			"pipeline-default-deny",
			"pipeline-controller",
			"pipeline-webhook",
			"pipeline-events-controller",
			"pipeline-resolvers",
			"tekton-proxy-webhook-default-deny",
			"proxy-webhook",
		},
	},
	"triggers": {
		deployment: "tekton-triggers-controller",
		policies: []string{
			"tekton-default-deny",
			"triggers-controller",
			"triggers-webhook",
			"triggers-core-interceptors",
		},
	},
	"chains": {
		deployment: "tekton-chains-controller",
		policies: []string{
			"chains-controller-default-deny",
			"chains-controller",
		},
	},
	"results": {
		deployment: "tekton-results-api",
		policies: []string{
			"results-default-deny",
			"results-api",
			"results-watcher",
			"results-retention-policy-agent",
			"results-postgres",
		},
	},
	"pruner": {
		deployment: "tekton-pruner-controller",
		policies: []string{
			"tekton-pruner-default-deny",
			"pruner-controller",
			"pruner-webhook",
		},
	},
	"manual-approval-gate": {
		deployment: "manual-approval-gate-controller",
		policies: []string{
			"mag-default-deny",
			"mag-controller",
			"mag-webhook",
		},
	},
	"pac": {
		deployment: "pipelines-as-code-controller",
		policies: []string{
			"pac-default-deny",
			"pac-controller",
			"pac-watcher",
			"pac-webhook",
		},
	},
}

var _ = gauge.Step("Verify NetworkPolicies are <status> in namespace <namespace> for component <component>", func(status, namespace, component string) {
	cfg, ok := networkPoliciesByComponent[component]
	if !ok {
		testsuit.T.Errorf("Unknown component %q, valid components: %v", component, validComponentNames())
		return
	}

	if !isDeploymentPresent(cfg.deployment, namespace) {
		log.Printf("Component %q is not deployed (deployment %q not found in %q), skipping NetworkPolicy check", component, cfg.deployment, namespace)
		return
	}

	shouldBePresent := status == "present"
	for _, policyName := range cfg.policies {
		if shouldBePresent {
			waitForNetworkPolicyPresent(policyName, namespace)
		} else {
			waitForNetworkPolicyAbsent(policyName, namespace)
		}
	}
})

var _ = gauge.Step("Disable NetworkPolicy on TektonConfig", func() {
	patchData := `{"spec":{"networkPolicy":{"disabled":true}}}`
	oc.UpdateTektonConfig(patchData)
	log.Println("Disabled NetworkPolicy on TektonConfig")
	patchStandaloneComponentNetworkPolicy(true)
})

var _ = gauge.Step("Enable NetworkPolicy on TektonConfig", func() {
	patchData := `{"spec":{"networkPolicy":{"disabled":false}}}`
	oc.UpdateTektonConfig(patchData)
	log.Println("Enabled NetworkPolicy on TektonConfig")
	patchStandaloneComponentNetworkPolicy(false)
})

var _ = gauge.Step("Wait for TektonConfig to be ready", func() {
	waitForTektonConfigReady()
})

var _ = gauge.Step("Add custom NetworkPolicy <name> on TektonConfig", func(name string) {
	patchData := fmt.Sprintf(`{"spec":{"networkPolicy":{"policies":{"%s":{"podSelector":{},"policyTypes":["Ingress","Egress"]}}}}}`, name)
	oc.UpdateTektonConfig(patchData)
	log.Printf("Added custom NetworkPolicy %q on TektonConfig", name)
})

var _ = gauge.Step("Remove custom NetworkPolicy from TektonConfig", func() {
	patchData := `[{"op":"remove","path":"/spec/networkPolicy/policies"}]`
	result := cmd.Run("oc", "patch", "tektonconfig", "config", "--type=json", "-p", patchData)
	log.Printf("Removed custom NetworkPolicy policies from TektonConfig: %s", result.Stdout())
})

var _ = gauge.Step("Add custom NetworkPolicy <name> on ManualApprovalGate", func(name string) {
	result := cmd.Run("oc", "get", "manualapprovalgate", "manual-approval-gate", "-o", "name")
	if result.ExitCode != 0 || !strings.Contains(result.Stdout(), "manual-approval-gate") {
		log.Printf("ManualApprovalGate not found, skipping custom NetworkPolicy add")
		return
	}
	patchData := fmt.Sprintf(`{"spec":{"networkPolicy":{"policies":{"%s":{"podSelector":{},"policyTypes":["Ingress","Egress"]}}}}}`, name)
	log.Printf("output: %s\n", cmd.Run("oc", "patch", "manualapprovalgate", "manual-approval-gate",
		"--type=merge", "-p", patchData).Stdout())
	log.Printf("Added custom NetworkPolicy %q on ManualApprovalGate", name)
})

var _ = gauge.Step("Remove custom NetworkPolicy from ManualApprovalGate", func() {
	result := cmd.Run("oc", "get", "manualapprovalgate", "manual-approval-gate", "-o", "name")
	if result.ExitCode != 0 || !strings.Contains(result.Stdout(), "manual-approval-gate") {
		log.Printf("ManualApprovalGate not found, skipping custom NetworkPolicy remove")
		return
	}
	patchData := `[{"op":"remove","path":"/spec/networkPolicy/policies"}]`
	result = cmd.Run("oc", "patch", "manualapprovalgate", "manual-approval-gate", "--type=json", "-p", patchData)
	log.Printf("Removed custom NetworkPolicy policies from ManualApprovalGate: %s", result.Stdout())
})

var _ = gauge.Step("Verify NetworkPolicy <name> is <status> in namespace <namespace>", func(name, status, namespace string) {
	shouldBePresent := status == "present"
	if shouldBePresent {
		waitForNetworkPolicyPresent(name, namespace)
	} else {
		waitForNetworkPolicyAbsent(name, namespace)
	}
})

func isDeploymentPresent(name, namespace string) bool {
	result := cmd.Run("oc", "get", "deployment", name, "-n", namespace, "-o", "name")
	return result.ExitCode == 0 && strings.Contains(result.Stdout(), name)
}

func waitForNetworkPolicyPresent(name, namespace string) {
	log.Printf("Waiting for NetworkPolicy %q to be present in namespace %q", name, namespace)
	deadline := time.Now().Add(5 * time.Minute)
	for time.Now().Before(deadline) {
		result := cmd.Run("oc", "get", "networkpolicy", name, "-n", namespace, "-o", "name")
		if result.ExitCode == 0 && strings.Contains(result.Stdout(), name) {
			log.Printf("NetworkPolicy %q found in namespace %q", name, namespace)
			return
		}
		time.Sleep(5 * time.Second)
	}
	testsuit.T.Errorf("NetworkPolicy %q not found in namespace %q after timeout", name, namespace)
}

func waitForNetworkPolicyAbsent(name, namespace string) {
	log.Printf("Waiting for NetworkPolicy %q to be absent from namespace %q", name, namespace)
	deadline := time.Now().Add(5 * time.Minute)
	for time.Now().Before(deadline) {
		result := cmd.Run("oc", "get", "networkpolicy", name, "-n", namespace, "-o", "name")
		if result.ExitCode != 0 || !strings.Contains(result.Stdout(), name) {
			log.Printf("NetworkPolicy %q absent from namespace %q", name, namespace)
			return
		}
		time.Sleep(5 * time.Second)
	}
	testsuit.T.Errorf("NetworkPolicy %q still present in namespace %q after timeout", name, namespace)
}

func waitForTektonConfigReady() {
	log.Println("Waiting for TektonConfig to be ready")
	deadline := time.Now().Add(10 * time.Minute)
	for time.Now().Before(deadline) {
		result := cmd.Run("oc", "get", "tektonconfig", "config", "-o", "jsonpath={.status.conditions[?(@.type==\"Ready\")].status}")
		if result.ExitCode == 0 && strings.TrimSpace(result.Stdout()) == "True" {
			log.Println("TektonConfig is ready")
			return
		}
		time.Sleep(10 * time.Second)
	}
	testsuit.T.Errorf("TektonConfig did not reach ready state after timeout")
}

// patchStandaloneComponentNetworkPolicy patches the networkPolicy field on
// standalone CRs (ManualApprovalGate) that are not managed by TektonConfig.
func patchStandaloneComponentNetworkPolicy(disabled bool) {
	patchData := fmt.Sprintf(`{"spec":{"networkPolicy":{"disabled":%t}}}`, disabled)
	result := cmd.Run("oc", "get", "manualapprovalgate", "manual-approval-gate", "-o", "name")
	if result.ExitCode == 0 && strings.Contains(result.Stdout(), "manual-approval-gate") {
		log.Printf("output: %s\n", cmd.Run("oc", "patch", "manualapprovalgate", "manual-approval-gate",
			"--type=merge", "-p", patchData).Stdout())
		log.Printf("Patched ManualApprovalGate networkPolicy.disabled=%t", disabled)
	}
}

func validComponentNames() []string {
	names := make([]string, 0, len(networkPoliciesByComponent))
	for k := range networkPoliciesByComponent {
		names = append(names, fmt.Sprintf("%q", k))
	}
	return names
}
