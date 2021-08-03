package triggers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	"github.com/tektoncd/pipeline/pkg/names"
	eventReconciler "github.com/tektoncd/triggers/pkg/reconciler/v1alpha1/eventlistener"
	"github.com/tektoncd/triggers/pkg/sink"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

func getServiceNameAndPort(c *clients.Clients, elname, namespace string) (string, string) {
	// Verify the EventListener to be ready
	err := wait.WaitFor(c.Ctx, wait.EventListenerReady(c, namespace, elname))
	assert.NoError(err, fmt.Sprintf("EventListener not %s ready", elname))

	labelSelector := fields.SelectorFromSet(eventReconciler.GenerateResourceLabels(elname)).String()
	// Grab EventListener sink pods
	sinkPods, err := c.KubeClient.Kube.CoreV1().Pods(namespace).List(c.Ctx, metav1.ListOptions{LabelSelector: labelSelector})
	assert.NoError(err, fmt.Sprintf("Error listing EventListener sink pods"))
	log.Printf("sinkpod name: %s", sinkPods.Items[0].Name)

	serviceList, err := c.KubeClient.Kube.CoreV1().Services(namespace).List(c.Ctx, metav1.ListOptions{LabelSelector: labelSelector})
	assert.NoError(err, fmt.Sprintf("Error listing services"))
	return serviceList.Items[0].Name, serviceList.Items[0].Spec.Ports[0].Name
}

func ExposeEventListner(c *clients.Clients, elname, namespace string) string {
	svcName, _ := getServiceNameAndPort(c, elname, namespace)
	cmd.MustSucceed("oc", "expose", "service", svcName, "-n", namespace)

	route := cmd.MustSucceed("oc", "-n", namespace, "get", "route", "--selector=eventlistener="+elname, "-o", "jsonpath='{range .items[*]}{.metadata.name}'").Stdout()

	route_url := cmd.MustSucceed("oc", "-n", namespace, "get", "route", strings.Trim(route, "'"), "--template='http://{{.spec.host}}'").Stdout()
	log.Printf("Route url: %s", route_url)

	time.Sleep(5 * time.Second)
	return strings.Trim(route_url, "'")
}

func ExposeEventListnerForTLS(c *clients.Clients, elname, namespace string) string {
	svcName, portName := getServiceNameAndPort(c, elname, namespace)
	domain := getDomain()
	cmd.MustSucceed("mkdir", "-p", resource.Path("testdata/triggers/certs")).Stdout()

	rootcaKey := resource.Path("testdata/triggers/certs/ca.key")
	rootcaCert := resource.Path("testdata/triggers/certs/ca.crt")
	tlsKey := resource.Path("testdata/triggers/certs/server.key")
	tlsCert := resource.Path("testdata/triggers/certs/server.crt")
	tlsCsr := resource.Path("testdata/triggers/certs/server.csr")
	serverEXT := resource.Path("testdata/triggers/certs/server.ext")

	cmd.MustSucceed("openssl", "genrsa", "-out", rootcaKey, "4096").Stdout()

	cmd.MustSucceed("openssl", "req", "-x509", "-new", "-nodes", "-key", rootcaKey,
		"-sha256", "-days", "4096", "-out", rootcaCert, "-subj",
		"/C=IN/ST=Kar/L=Blr/O=RedHat").Stdout()

	cmd.MustSucceed("openssl", "genrsa", "-out", tlsKey, "4096").Stdout()

	cmd.MustSucceed("openssl", "req", "-new", "-key", tlsKey, "-out", tlsCsr,
		"-subj", fmt.Sprintf("/C=IN/ST=Kar/L=Blr/O=RedHat/CN=%s", domain)).Stdout()

	extData := fmt.Sprintf("authorityKeyIdentifier=keyid,issuer\nbasicConstraints=CA:FALSE\nkeyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment\n"+
		"subjectAltName = @alt_names\n\n\n[alt_names]\nDNS.1 = %s\n", domain)

	err := ioutil.WriteFile(serverEXT, []byte(extData), 0644)
	assert.FailOnError(err)

	cmd.MustSucceed("openssl", "x509", "-req", "-in", tlsCsr,
		"-CA", rootcaCert, "-CAkey",
		rootcaKey, "-CAcreateserial", "-out",
		tlsCert,
		"-days", "4096", "-sha256", "-extfile", serverEXT).Stdout()

	routeName := cmd.MustSucceed("oc", "create", "route", "reencrypt", "--ca-cert="+rootcaCert,
		"--cert="+tlsCert, "--key="+tlsKey,
		"--service="+svcName, "--hostname="+domain, "--port="+portName, "-n", namespace).Stdout()

	route_url := cmd.MustSucceed("oc", "-n", namespace, "get", strings.Split(routeName, " ")[0], "--template='http://{{.spec.host}}'").Stdout()
	log.Printf("Route url: %s", route_url)

	time.Sleep(5 * time.Second)
	return strings.Trim(route_url, "'")
}

// This function returns the formatted hostname.
func getDomain() string {
	/* each cluster have different domain so below logic is to extract domain from existing route
	each openshift installation have `console` route in `openshift-console` namespace.
	ex: http://console-openshift-console.apps.tt3.testing */
	route_url := cmd.MustSucceed("oc", "-n", "openshift-console", "get", "route", "console", "--template=http://{{.spec.host}}").Stdout()
	splittedValue := strings.SplitAfter(route_url, ".apps")
	routeDomainName := "apps" + splittedValue[1]
	randomName := names.SimpleNameGenerator.RestrictLengthWithRandomSuffix("releasetest")
	return "tls." + randomName + "." + routeDomainName
}

func MockPostEventWithEmptyPayload(routeurl string) *http.Response {
	// Send empty POST request to EventListener sink
	req, err := http.NewRequest("POST", routeurl, bytes.NewBuffer([]byte("{}")))
	assert.FailOnError(err)
	req.Header.Add("Accept", "application/json")
	assert.FailOnError(err)

	resp, err := CreateHTTPClient().Do(req)
	assert.FailOnError(err)
	if resp.StatusCode > http.StatusAccepted {
		testsuit.T.Errorf(fmt.Sprintf("sink did not return 2xx response. Got status code: %d", resp.StatusCode))
	}
	return resp
}

func MockPostEvent(routeurl, interceptor, eventType, payload string, isTLS bool) *http.Response {
	var (
		req  *http.Request
		err  error
		resp *http.Response
	)
	eventBodyJSON, err := ioutil.ReadFile(resource.Path(payload))
	assert.NoError(err, fmt.Sprintf("Couldn't load test data"))
	gauge.GetScenarioStore()["payload"] = eventBodyJSON

	// Send POST request to EventListener sink
	if isTLS {
		req, err = http.NewRequest("POST", "https://"+strings.Split(routeurl, "//")[1], bytes.NewBuffer(eventBodyJSON))
	} else {
		req, err = http.NewRequest("POST", routeurl, bytes.NewBuffer(eventBodyJSON))
	}
	assert.FailOnError(err)

	req = buildHeaders(req, interceptor, eventType)

	if isTLS {
		resp, err = CreateHTTPSClient().Do(req)
	} else {
		resp, err = CreateHTTPClient().Do(req)
	}
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
	sinkPods, err := c.KubeClient.Kube.CoreV1().Pods(namespace).List(c.Ctx, metav1.ListOptions{LabelSelector: labelSelector})
	assert.NoError(err, fmt.Sprintf("Error listing EventListener sink pods"))
	logs := cmd.MustSucceed("oc", "-n", namespace, "logs", "pods/"+sinkPods.Items[0].Name, "--all-containers", "--tail=2").Stdout()
	if strings.Contains(logs, "error") {
		testsuit.T.Errorf("Error: sink logs: \n %s", logs)
		gauge.WriteMessage(fmt.Sprintf("sink logs: \n %s", logs))
	}
}

func CleanupTriggers(c *clients.Clients, elName, namespace string) {
	// Delete EventListener
	err := c.TriggersClient.TriggersV1alpha1().EventListeners(namespace).Delete(c.Ctx, elName, metav1.DeleteOptions{})
	assert.FailOnError(err)

	log.Println("Deleted EventListener")

	// Verify the EventListener's Deployment is deleted
	err = wait.WaitFor(c.Ctx, wait.DeploymentNotExist(c, namespace, fmt.Sprintf("%s-%s", eventReconciler.GeneratedResourcePrefix, elName)))
	assert.FailOnError(err)

	log.Println("EventListener's Deployment was deleted")

	// Verify the EventListener's Service is deleted
	err = wait.WaitFor(c.Ctx, wait.ServiceNotExist(c, namespace, fmt.Sprintf("%s-%s", eventReconciler.GeneratedResourcePrefix, elName)))
	assert.FailOnError(err)

	log.Println("EventListener's Service was deleted")

	//Delete Route exposed earlier
	err = c.Route.Routes(namespace).Delete(c.Ctx, fmt.Sprintf("%s-%s", eventReconciler.GeneratedResourcePrefix, elName), metav1.DeleteOptions{})
	assert.FailOnError(err)

	// Verify the EventListener's Route is deleted
	err = wait.WaitFor(c.Ctx, wait.RouteNotExist(c, namespace, fmt.Sprintf("%s-%s", eventReconciler.GeneratedResourcePrefix, elName)))
	assert.FailOnError(err)
	log.Println("EventListener's Route got deleted successfully...")

	// This is required when EL runs as TLS
	cmd.MustSucceed("rm", "-rf", os.Getenv("GOPATH")+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs")
}
