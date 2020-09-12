package triggers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/gauge"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	resource "github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/wait"
	eventReconciler "github.com/tektoncd/triggers/pkg/reconciler/v1alpha1/eventlistener"
	"github.com/tektoncd/triggers/pkg/sink"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

func ExposeEventListner(c *clients.Clients, elname, namespace string) string {
	// Verify the EventListener to be ready
	err := wait.WaitFor(wait.EventListenerReady(c, namespace, elname))
	assert.NoError(err, fmt.Sprintf("EventListener not %s ready", elname))

	labelSelector := fields.SelectorFromSet(eventReconciler.GenerateResourceLabels(elname)).String()
	// Grab EventListener sink pods
	sinkPods, err := c.KubeClient.Kube.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
	assert.NoError(err, fmt.Sprintf("Error listing EventListener sink pods"))
	log.Printf("sinkpod name: %s", sinkPods.Items[0].Name)

	serviceList, err := c.KubeClient.Kube.CoreV1().Services(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
	assert.NoError(err, fmt.Sprintf("Error listing services"))

	cmd.MustSucceed("oc", "expose", "service", serviceList.Items[0].Name, "-n", namespace)

	route := cmd.MustSucceed("oc", "-n", namespace, "get", "route", "--selector=eventlistener="+elname, "-o", "jsonpath='{range .items[*]}{.metadata.name}'").Stdout()

	route_url := cmd.MustSucceed("oc", "-n", namespace, "get", "route", strings.Trim(route, "'"), "--template='http://{{.spec.host}}'").Stdout()
	log.Printf("Route url: %s", route_url)

	time.Sleep(5 * time.Second)
	return strings.Trim(route_url, "'")
}

func MockGetEvent(routeurl string) *http.Response {
	// Send GET request to EventListener sink
	req, err := http.NewRequest("GET", routeurl, nil)
	req.Header.Add("Accept", "application/json")
	assert.FailOnError(err)

	resp, err := CreateHTTPClient().Do(req)
	assert.FailOnError(err)
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		testsuit.T.Errorf(fmt.Sprintf("sink returned 401/403 response: %d", resp.StatusCode))
	}
	return resp
}

// func MockPostEvent(routeurl, headers, payload string) *http.Response {
// 	eventBodyJSON, err := ioutil.ReadFile(resource.Path(payload))
// 	assert.NoError(err, fmt.Sprintf("Couldn't load test data"))
// 	gauge.GetScenarioStore()["payload"] = eventBodyJSON

// 	// Send POST request to EventListener sink
// 	req, err := http.NewRequest("POST", routeurl, bytes.NewBuffer(eventBodyJSON))
// 	assert.FailOnError(err)
// 	req = buildHeaders(req, headers)

// 	resp, err := CreateHTTPClient().Do(req)
// 	assert.FailOnError(err)
// 	if resp.StatusCode > http.StatusAccepted {
// 		testsuit.T.Errorf(fmt.Sprintf("sink did not return 2xx response. Got status code: %d", resp.StatusCode))
// 	}
// 	return resp
// }

func MockPostEvent(routeurl, interceptor, eventType, payload string) *http.Response {
	eventBodyJSON, err := ioutil.ReadFile(resource.Path(payload))
	assert.NoError(err, fmt.Sprintf("Couldn't load test data"))
	gauge.GetScenarioStore()["payload"] = eventBodyJSON

	// Send POST request to EventListener sink
	req, err := http.NewRequest("POST", routeurl, bytes.NewBuffer(eventBodyJSON))
	assert.FailOnError(err)

	req = buildHeaders(req, interceptor, eventType)

	resp, err := CreateHTTPClient().Do(req)
	assert.FailOnError(err)
	if resp.StatusCode > http.StatusAccepted {
		testsuit.T.Errorf(fmt.Sprintf("sink did not return 2xx response. Got status code: %d", resp.StatusCode))
	}
	return resp
}

func AssertElResponse(c *clients.Clients, resp *http.Response, elname, namespace string) {
	wantBody := sink.Response{
		EventListener: elname,
		Namespace:     namespace,
	}

	defer resp.Body.Close()
	var gotBody sink.Response
	err := json.NewDecoder(resp.Body).Decode(&gotBody)
	assert.FailOnError(err)

	if diff := cmp.Diff(wantBody, gotBody, cmpopts.IgnoreFields(sink.Response{}, "EventID")); diff != "" {
		testsuit.T.Errorf(fmt.Sprintf("unexpected sink response -want/+got: %s", diff))
	}

	if gotBody.EventID == "" {
		testsuit.T.Errorf("sink response no eventID")
	}

	labelSelector := fields.SelectorFromSet(eventReconciler.GenerateResourceLabels(elname)).String()
	// Grab EventListener sink pods
	sinkPods, err := c.KubeClient.Kube.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
	assert.NoError(err, fmt.Sprintf("Error listing EventListener sink pods"))
	logs := cmd.MustSucceed("oc", "-n", namespace, "logs", "pods/"+sinkPods.Items[0].Name, "--all-containers", "--tail=2").Stdout()
	if strings.Contains(logs, "error") {
		testsuit.T.Errorf("Error: sink logs: \n %s", logs)
		gauge.WriteMessage(fmt.Sprintf("sink logs: \n %s", logs))
	}
}

func CleanupTriggers(c *clients.Clients, elName, namespace string) {
	// Delete EventListener
	err := c.TriggersClient.TriggersV1alpha1().EventListeners(namespace).Delete(elName, &metav1.DeleteOptions{})
	assert.FailOnError(err)

	log.Println("Deleted EventListener")

	// Verify the EventListener's Deployment is deleted
	err = wait.WaitFor(wait.DeploymentNotExist(c, namespace, fmt.Sprintf("%s-%s", eventReconciler.GeneratedResourcePrefix, elName)))
	assert.FailOnError(err)

	log.Println("EventListener's Deployment was deleted")

	// Verify the EventListener's Service is deleted
	err = wait.WaitFor(wait.ServiceNotExist(c, namespace, fmt.Sprintf("%s-%s", eventReconciler.GeneratedResourcePrefix, elName)))
	assert.FailOnError(err)

	log.Println("EventListener's Service was deleted")
}
