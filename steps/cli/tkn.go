package cli

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"regexp"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/openshift-pipelines/release-tests/pkg/tkn"
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
			pipelineRunName := tkn.StartPipeline(pipelineName, params, workspaces, store.Namespace(), "--use-param-defaults")
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
			pipelineRunName := tkn.StartPipeline(pipelineName, params, workspaces, store.Namespace(), "--use-param-defaults")
			pipelines.ValidatePipelineRun(store.Clients(), pipelineRunName, "successful", "no", store.Namespace())
		}(pipelineName, params, workspaces)

		time.Sleep(3 * time.Second)
	}
	wg.Wait()
})
