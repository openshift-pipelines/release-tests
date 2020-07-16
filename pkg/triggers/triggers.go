package triggers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/wait"
	eventReconciler "github.com/tektoncd/triggers/pkg/reconciler/v1alpha1/eventlistener"
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

func MockGetEvent(routeurl string) {
	// Send GET request to EventListener sink
	req, err := http.NewRequest("GET", routeurl, nil)
	req.Header.Add("Accept", "application/json")
	assert.NoError(err, fmt.Sprintf("Error creating GET request for trigger "))
	resp, err := CreateHTTPClient().Do(req)
	assert.NoError(err, fmt.Sprintf("Error Sending GET request for trigger "))
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		testsuit.T.Errorf(fmt.Sprintf("sink returned 401/403 response: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	log.Println(resp.Status)
	gauge.WriteMessage("%s", string(resp_body))
}
