package cli

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"regexp"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/opc"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Start and verify pipeline <pipelineName> with param <paramName> with values stored in variable <variableName> with workspace <workspaceValue>", func(pipelineName, paramName, variableName, workspaceValue string) {
	values := store.GetScenarioDataSlice(variableName)
	params := make(map[string]string)
	workspaces := make(map[string]string)
	workspaces[strings.Split(workspaceValue, ",")[0]] = strings.Split(workspaceValue, ",")[1]
	var wg sync.WaitGroup // wait for all goroutine to finish
	for _, value := range values {
		if value == "latest" {
			continue
		}

		wg.Add(1)

		go func(value string) {
			defer wg.Done()
			log.Printf("Starting pipeline %s with param %s=%s", pipelineName, paramName, value)
			params[paramName] = value
			customPipelineRunName := pipelineName + "-run-" + value
			pipelineRunName := opc.StartPipeline(pipelineName, params, workspaces, store.Namespace(), "--use-param-defaults", "--prefix-name", customPipelineRunName)
			pipelines.ValidatePipelineRun(store.Clients(), pipelineRunName, "successful", "no", store.Namespace())
		}(value)

		time.Sleep(3 * time.Second)
	}
	wg.Wait()
})

var _ = gauge.Step("Start and verify dotnet pipeline <pipelineName> with values stored in variable <variableName> with workspace <workspaceValue>", func(pipelineName, variableName, workspaceValue string) {
	values := store.GetScenarioDataSlice(variableName)
	params := make(map[string]string)
	workspaces := make(map[string]string)
	workspaces[strings.Split(workspaceValue, ",")[0]] = strings.Split(workspaceValue, ",")[1]
	paramName := "VERSION"
	re := regexp.MustCompile(`\d+\.\d+`)
	var exampleRevision string
	var wg sync.WaitGroup // wait for all goroutine to finish
	for _, value := range values {
		if value == "latest" {
			continue
		}
		exampleRevision = re.FindString(value)
		params[paramName] = value
		versionInt, _ := strconv.ParseFloat(exampleRevision, 64)
		if versionInt >= 5.0 {
			params["EXAMPLE_REVISION"] = "dotnet-" + exampleRevision
		} else {
			params["EXAMPLE_REVISION"] = "dotnetcore-" + exampleRevision
		}
		log.Printf("Starting pipeline %s with param %s=%s and EXAMPLE_REVISION=%s", pipelineName, paramName, value, params["EXAMPLE_REVISION"])

		wg.Add(1)

		go func(pipelineName string, params map[string]string, workspaces map[string]string) {
			defer wg.Done()
			customPipelineRunName := pipelineName + "-run-" + value
			pipelineRunName := opc.StartPipeline(pipelineName, params, workspaces, store.Namespace(), "--use-param-defaults", "--prefix-name", customPipelineRunName)
			pipelines.ValidatePipelineRun(store.Clients(), pipelineRunName, "successful", "no", store.Namespace())
		}(pipelineName, params, workspaces)

		time.Sleep(3 * time.Second)
	}
	wg.Wait()
})

var _ = gauge.Step("Start the <pipelineName> pipeline with params <parameters> with workspace <workspaceValue> and store the pipelineRunName to variable <variableName>", func(pipelineName, parameters, workspaceValue, variableName string) {
	params := make(map[string]string)
	paramPairs := strings.Split(parameters, ",")
	for _, param := range paramPairs {
		keyValue := strings.Split(param, "=")
		if len(keyValue) == 2 {
			params[keyValue[0]] = keyValue[1]
		}
	}
	workspaces := make(map[string]string)
	workspaces[strings.Split(workspaceValue, ",")[0]] = strings.Split(workspaceValue, ",")[1]
	pipelineRunName := opc.StartPipeline(pipelineName, params, workspaces, store.Namespace(), "--use-param-defaults")
	prList, err := opc.GetOpcPrList(pipelineRunName, store.Namespace())
	if err != nil {
		testsuit.T.Errorf("Failed to get pipelineRun %s: %v", pipelineRunName, err)
	}
	if len(prList) == 0 || prList[0].Name != pipelineRunName {
		testsuit.T.Errorf("pipelineRun %s not found", pipelineRunName)
	}
	pipelines.ValidatePipelineRun(store.Clients(), pipelineRunName, "successful", "no", store.Namespace())
	store.PutScenarioData(variableName, pipelineRunName)
})

var _ = gauge.Step("Hub Search for <resource>", func(resource string) {
	if err := opc.HubSearch(resource); err != nil {
		testsuit.T.Errorf("Hub search error: %v", err)
	}
})

var _ = gauge.Step("Verify that <resourceType> <resourceName> exists", func(resourcetype, resourcename string) {
	if _, err := opc.VerifyResourceListMatchesName(resourcetype, resourcename, store.Namespace()); err != nil {
		testsuit.T.Errorf("Failed to verify %s with %s in %s failed: %v", resourcetype, resourcename, store.Namespace(), err)
	} else {
		log.Printf("%s %s exists", resourcetype, resourcename)
	}
})
