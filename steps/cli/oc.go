package cli

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/getgauge-contrib/gauge-go/gauge"
	m "github.com/getgauge-contrib/gauge-go/models"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/oc"
	"github.com/openshift-pipelines/release-tests/pkg/openshift"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Create <table>", func(table *m.Table) {
	for _, row := range table.Rows {
		resource := row.Cells[1]
		oc.Create(resource, store.Namespace())
	}
})

var _ = gauge.Step("Create remote <table>", func(table *m.Table) {
	for _, row := range table.Rows {
		resource := row.Cells[1]
		expr := regexp.MustCompile("{.+}")
		matchedString := expr.FindString(resource)
		if matchedString != "" {
			envVariable := strings.Replace(matchedString, "{", "", 1)
			envVariable = strings.Replace(envVariable, "}", "", 1)
			resource = strings.Replace(resource, matchedString, os.Getenv(envVariable), 1)
		}
		oc.CreateRemote(resource, store.Namespace())
	}
})

var _ = gauge.Step("Apply <table>", func(table *m.Table) {
	for _, row := range table.Rows {
		resource := row.Cells[1]
		oc.Apply(resource, store.Namespace())
	}
})

var _ = gauge.Step("Apply in namespace <ns> <table>", func(ns string, table *m.Table) {
	for _, row := range table.Rows {
		resource := row.Cells[1]
		oc.Apply(resource, ns)
	}
})

var _ = gauge.Step("Enable TLS config for eventlisteners", func() {
	oc.EnableTLSConfigForEventlisteners(store.Namespace())

})

var _ = gauge.Step("Verify kubernetes events for eventlistener", func() {
	oc.VerifyKubernetesEventsForEventListener(store.Namespace())
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

var _ = gauge.Step("Update pruner config <keepPresence> keep <keep> schedule <schedule> resources <resources> and <keepSincePresence> keep-since <keepSince>", func(keepPresence, keep, schedule, resources, keepSincePresence, keepSince string) {
	resourcesSplit := strings.Split(resources, ",")
	resourcesList := strings.Join(resourcesSplit, "\",\"")
	patch_data := ""
	if keepPresence == "with" && keepSincePresence == "without" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep\":%s,\"keep-since\":null,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", keep, schedule, resourcesList)
	} else if keepPresence == "without" && keepSincePresence == "with" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep\":null,\"keep-since\":%s,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", keepSince, schedule, resourcesList)
	} else if keepPresence == "with" && keepSincePresence == "with" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep\":%s,\"keep-since\":%s,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", keep, keepSince, schedule, resourcesList)
	} else if keepPresence == "without" && keepSincePresence == "without" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep\":null,\"keep-since\":null,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", schedule, resourcesList)
	}
	oc.UpdateTektonConfig(patch_data)
})

var _ = gauge.Step("Update pruner config with invalid data <keepPresence> keep <keep> schedule <schedule> resources <resources> and <keepSincePresence> keep-since <keepSince> and expect error message <errorMessage>", func(keepPresence, keep, schedule, resources, keepSincePresence, keepSince, errorMessage string) {
	resourcesSplit := strings.Split(resources, ",")
	resourcesList := strings.Join(resourcesSplit, "\",\"")
	patch_data := ""
	if keepPresence == "with" && keepSincePresence == "without" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep\":%s,\"keep-since\":null,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", keep, schedule, resourcesList)
	} else if keepPresence == "without" && keepSincePresence == "with" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep\":null,\"keep-since\":%s,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", keepSince, schedule, resourcesList)
	} else if keepPresence == "with" && keepSincePresence == "with" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep\":%s,\"keep-since\":%s,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", keep, keepSince, schedule, resourcesList)
	} else if keepPresence == "without" && keepSincePresence == "without" {
		patch_data = fmt.Sprintf("{\"spec\":{\"pruner\":{\"keep\":null,\"keep-since\":null,\"schedule\":\"%s\",\"resources\":[\"%s\"]}}}", schedule, resourcesList)
	}
	oc.UpdateTektonConfigwithInvalidData(patch_data, errorMessage)
})

var _ = gauge.Step("Remove auto pruner configuration from config CR", func() {
	log.Print("Removing pruner configuration from config CR")
	oc.RemovePrunerConfig()
})

var _ = gauge.Step("Annotate namespace with <annotation>", func(annotation string) {
	log.Printf("Annotating namespace %v with %v", store.Namespace(), annotation)
	oc.AnnotateNamespace(store.Namespace(), annotation)
})

var _ = gauge.Step("Remove annotation <annotation> from namespace", func(annotation string) {
	log.Printf("Removing annotation %v from namespace %v", store.Namespace(), annotation)
	oc.AnnotateNamespace(store.Namespace(), annotation+"-")
})

var _ = gauge.Step("Add label <label> to namespace", func(label string) {
	log.Printf("Labelling namespace %v with %v", store.Namespace(), label)
	oc.LabelNamespace(store.Namespace(), label)
})

var _ = gauge.Step("Remove label <label> from the namespace", func(label string) {
	log.Printf("Removing annotation %v from namespace %v", store.Namespace(), label)
	oc.AnnotateNamespace(store.Namespace(), label+"-")
})

var _ = gauge.Step("Update addon config with clusterTasks as <clusterTaskStatus> communityClustertasks as <comClusterTaskStatus> and pipelineTemplates as <pipelineTemplateStatus> and expect message <expectedMessage>", func(clusterTaskStatus, commClustertaskStatus, pipeTemplateStatus, expectedMessage string) {
	patchData := fmt.Sprintf("{\"spec\":{\"addon\":{\"params\":[{\"name\":\"clusterTasks\",\"value\":\"%s\"},{\"name\":\"communityClusterTasks\",\"value\":\"%s\"},{\"name\":\"pipelineTemplates\",\"value\":\"%s\"}]}}}", clusterTaskStatus, commClustertaskStatus, pipeTemplateStatus)
	if expectedMessage == "" {
		oc.UpdateTektonConfig(patchData)
	} else {
		oc.UpdateTektonConfigwithInvalidData(patchData, expectedMessage)
	}
})

var _ = gauge.Step("Update addon config with resolverTasks as <resolverTaskStatus> and pipelineTemplates as <pipelineTemplateStatus> and expect message <expectedMessage>", func(resolverTaskStatus, pipeTemplateStatus, expectedMessage string) {
	patchData := fmt.Sprintf("{\"spec\":{\"addon\":{\"params\":[{\"name\":\"resolverTasks\",\"value\":\"%s\"},{\"name\":\"pipelineTemplates\",\"value\":\"%s\"}]}}}", resolverTaskStatus, pipeTemplateStatus)
	if expectedMessage == "" {
		oc.UpdateTektonConfig(patchData)
	} else {
		oc.UpdateTektonConfigwithInvalidData(patchData, expectedMessage)
	}
})

var _ = gauge.Step("Update addon config with resolverTasks as <resolverTasksStatus> and expect message <expectedMessage>", func(resolverTasksStatus, expectedMessage string) {
	patchData := fmt.Sprintf("{\"spec\":{\"addon\":{\"params\":[{\"name\":\"resolverTasks\",\"value\":\"%s\"}]}}}", resolverTasksStatus)
	if expectedMessage == "" {
		oc.UpdateTektonConfig(patchData)
	} else {
		oc.UpdateTektonConfigwithInvalidData(patchData, expectedMessage)
	}
})

var _ = gauge.Step("Update addon config with resolverStepActions as <resolverStepActionsStatus> and expect message <expectedMessage>", func(resolverStepActionsStatus, expectedMessage string) {
	patchData := fmt.Sprintf("{\"spec\":{\"addon\":{\"params\":[{\"name\":\"resolverStepActions\",\"value\":\"%s\"}]}}}", resolverStepActionsStatus)
	if expectedMessage == "" {
		oc.UpdateTektonConfig(patchData)
	} else {
		oc.UpdateTektonConfigwithInvalidData(patchData, expectedMessage)
	}
})

var _ = gauge.Step("Verify versioned ecosystem tasks", func() {
	operator.VerifyVersionedTasks()
})

var _ = gauge.Step("Verify versioned ecosystem step actions", func() {
	operator.VerifyVersionedStepActions()
})

var _ = gauge.Step("Create project <projectName>", func(projectName string) {
	log.Printf("Check if project %v already exists", projectName)
	if oc.CheckProjectExists(projectName) {
		log.Printf("Switch to project %v", projectName)
	} else {
		log.Printf("Creating project %v", projectName)
		oc.CreateNewProject(projectName)
	}
	store.Clients().NewClientSet(projectName)
	gauge.GetScenarioStore()["namespace"] = projectName
})

var _ = gauge.Step("Switch to autogenerated namespace", func() {
	gauge_store := gauge.GetScenarioStore()
	autogenerated_ns := gauge_store["autogenerated"].(string)
	if oc.CheckProjectExists(autogenerated_ns) {
		log.Printf("Switch to project %v", autogenerated_ns)
	}
	store.Clients().NewClientSet(autogenerated_ns)
	gauge.GetScenarioStore()["namespace"] = autogenerated_ns
})

var _ = gauge.Step("Delete project <projectName>", func(projectName string) {
	log.Printf("Deleting project %v", projectName)
	oc.DeleteProjectIgnoreErors(projectName)
})

var _ = gauge.Step("Link secret <secret> to service account <sa>", func(secret, sa string) {
	oc.LinkSecretToSA(secret, sa, store.Namespace())
})

var _ = gauge.Step("Delete <resourceType> named <name>", func(resourceType, name string) {
	oc.DeleteResource(resourceType, name)
})

var _ = gauge.Step("Define the tekton-hub-api variable", func() {
	patch_data := "{\"spec\":{\"pipeline\":{\"hub-resolver-config\":{\"tekton-hub-api\":\"https://api.hub.tekton.dev/\"}}}}"
	oc.UpdateTektonConfig(patch_data)

})

var _ = gauge.Step("Configure GitHub token for git resolver in TektonConfig", func() {
	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Printf("Token for authorization to the GitHub repository was not exported as a system variable")
	} else {
		if !oc.SecretExists("github-auth-secret", "openshift-pipelines") {
			secretData := os.Getenv("GITHUB_TOKEN")
			oc.CreateSecretForGitResolver(secretData)
		} else {
			log.Printf("Secret \"github-auth-secret\" already exists")
		}
		patch_data := "{\"spec\":{\"pipeline\":{\"git-resolver-config\":{\"api-token-secret-key\":\"github-auth-key\", \"api-token-secret-name\":\"github-auth-secret\", \"api-token-secret-namespace\":\"openshift-pipelines\", \"default-revision\":\"main\", \"fetch-timeout\":\"1m\", \"scm-type\":\"github\"}}}}"
		oc.UpdateTektonConfig(patch_data)
	}
})

var _ = gauge.Step("Configure the bundles resolver", func() {
	patch_data := "{\"spec\":{\"pipeline\":{\"bundles-resolver-config\":{\"default-kind\":\"task\", \"defaut-service-account\":\"pipelines\"}}}}"
	oc.UpdateTektonConfig(patch_data)
})

var _ = gauge.Step("Enable console plugin", func() {
	openshiftVersion := openshift.GetOpenShiftVersion(store.Clients())
	if openshiftVersion == "" {
		testsuit.T.Errorf("Unknown version of OpenShift (cluster version \"%v\").", openshiftVersion)
	}

	minorVersion, err := strconv.Atoi(strings.Split(openshiftVersion, ".")[1])

	if err != nil {
		testsuit.T.Fail(err)
	}

	if minorVersion < 15 {
		log.Printf("Console plugin is not supported on OpenShift version lower than 4.15 (cluster version %v).", openshiftVersion)
		return
	}

	oc.EnableConsolePlugin()
})

var _ = gauge.Step("Enable statefulset in tektonconfig", func() {
	patch_data := "{\"spec\":{\"pipeline\":{\"performance\":{\"disable-ha\":false,\"statefulset-ordinals\":true,\"replicas\":2,\"buckets\":2}}}}"
	oc.UpdateTektonConfig(patch_data)
})

var _ = gauge.Step("Configure Results with Loki", func() {
	patch_data := "{\"spec\":{\"result\":{\"auth_disable\":true,\"disabled\":false,\"log_level\":\"debug\",\"loki_stack_name\":\"logging-loki\",\"loki_stack_namespace\":\"openshift-logging\"}}}"
	oc.UpdateTektonConfig(patch_data)
})

var _ = gauge.Step("Copy secret <secretName> from <sourceNamespace> namespace to autogenerated namespace", func(secretName string, sourceNamespace string) {
	if oc.SecretExists(secretName, sourceNamespace) {
		oc.CopySecret(secretName, sourceNamespace, store.Namespace())
	} else {
		testsuit.T.Fail(fmt.Errorf("secret %s doesn't exist in namespace %s", secretName, sourceNamespace))
	}
})
