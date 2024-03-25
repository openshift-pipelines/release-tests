package triggers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/gauge"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	resource "github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/wait"
	"github.com/tektoncd/pipeline/pkg/names"
	eventReconciler "github.com/tektoncd/triggers/pkg/reconciler/eventlistener"
	"github.com/tektoncd/triggers/pkg/reconciler/eventlistener/resources"
	"github.com/tektoncd/triggers/pkg/sink"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

func getServiceNameAndPort(c *clients.Clients, elname, namespace string) (string, string) {
	// Verify the EventListener to be ready
	err := wait.WaitFor(c.Ctx, wait.EventListenerReady(c, namespace, elname))
	if err != nil {
		testsuit.T.Errorf("event listener %s in namespace %s not ready \n %v", elname, namespace, err)
	}

	labelSelector := fields.SelectorFromSet(resources.GenerateLabels(elname, resources.DefaultStaticResourceLabels)).String()
	// Grab EventListener sink pods
	sinkPods, err := c.KubeClient.Kube.CoreV1().Pods(namespace).List(c.Ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		testsuit.T.Errorf("failed to list event listener %s sink pods in namespace %s \n %v", elname, namespace, err)
	}

	log.Printf("sinkpod name: %s", sinkPods.Items[0].Name)

	serviceList, err := c.KubeClient.Kube.CoreV1().Services(namespace).List(c.Ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		testsuit.T.Errorf("failed to list services with label selector %s in namespace %s \n %v", labelSelector, namespace, err)
	}
	return serviceList.Items[0].Name, serviceList.Items[0].Spec.Ports[0].Name
}

func ExposeEventListner(c *clients.Clients, elname, namespace string) string {
	svcName, _ := getServiceNameAndPort(c, elname, namespace)
	cmd.MustSucceed("oc", "expose", "service", svcName, "-n", namespace)

	return GetRoute(elname, namespace)
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

	err := os.WriteFile(serverEXT, []byte(extData), 0644)
	if err != nil {
		testsuit.T.Fail(err)
	}

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
	// extract cluster's domain from ingress config, e.g. apps.mycluster.example.com
	routeDomainName := cmd.MustSucceed("oc", "get", "ingresses.config/cluster", "-o", "jsonpath={.spec.domain}").Stdout()
	randomName := names.SimpleNameGenerator.RestrictLengthWithRandomSuffix("rt")
	return "tls." + randomName + "." + routeDomainName
}

func MockPostEventWithEmptyPayload(routeurl string) *http.Response {
	// Send empty POST request to EventListener sink
	req, err := http.NewRequest("POST", routeurl, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		testsuit.T.Fail(err)
	}

	req.Header.Add("Accept", "application/json")
	resp, err := CreateHTTPClient().Do(req)
	if err != nil {
		testsuit.T.Fail(err)
	}

	if resp.StatusCode > http.StatusAccepted {
		testsuit.T.Errorf("sink did not return 2xx response. Got status code: %d", resp.StatusCode)
	}
	return resp
}

func MockPostEvent(routeurl, interceptor, eventType, payload string, isTLS bool) *http.Response {
	var (
		req  *http.Request
		err  error
		resp *http.Response
	)
	eventBodyJSON, err := os.ReadFile(resource.Path(payload))
	if err != nil {
		testsuit.T.Errorf("could not load test data from file %s \n %v", payload, err)
	}

	gauge.GetScenarioStore()["payload"] = eventBodyJSON

	// Send POST request to EventListener sink
	if isTLS {
		req, err = http.NewRequest("POST", "https://"+strings.Split(routeurl, "//")[1], bytes.NewBuffer(eventBodyJSON))
	} else {
		req, err = http.NewRequest("POST", routeurl, bytes.NewBuffer(eventBodyJSON))
	}
	if err != nil {
		testsuit.T.Fail(err)
	}

	req = buildHeaders(req, interceptor, eventType)

	if isTLS {
		resp, err = CreateHTTPSClient().Do(req)
	} else {
		resp, err = CreateHTTPClient().Do(req)
	}
	if err != nil {
		testsuit.T.Fail(err)
	}

	if resp.StatusCode > http.StatusAccepted {
		testsuit.T.Errorf("sink did not return 2xx response. Got status code: %d", resp.StatusCode)
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
	if err != nil {
		testsuit.T.Fail(err)
	}

	if diff := cmp.Diff(wantBody, gotBody, cmpopts.IgnoreFields(sink.Response{}, "EventID", "EventListenerUID")); diff != "" {
		testsuit.T.Errorf("unexpected sink response -want/+got: %s", diff)
	}

	if gotBody.EventID == "" {
		testsuit.T.Errorf("sink response no eventID")
	}

	labelSelector := fields.SelectorFromSet(resources.GenerateLabels(elname, resources.DefaultStaticResourceLabels)).String()
	// Grab EventListener sink pods
	sinkPods, err := c.KubeClient.Kube.CoreV1().Pods(namespace).List(c.Ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		testsuit.T.Errorf("failed to list event listener sink pods with label selector %s in namespace %s \n %v", labelSelector, namespace, err)
	}

	logs := cmd.MustSucceed("oc", "-n", namespace, "logs", "pods/"+sinkPods.Items[0].Name, "--all-containers", "--tail=2").Stdout()
	if strings.Contains(logs, "error") {
		testsuit.T.Errorf("Error: sink logs: \n %s", logs)
		gauge.WriteMessage(fmt.Sprintf("sink logs: \n %s", logs))
	}
}

func CleanupTriggers(c *clients.Clients, elName, namespace string) {
	// Delete EventListener
	err := c.TriggersClient.TriggersV1alpha1().EventListeners(namespace).Delete(c.Ctx, elName, metav1.DeleteOptions{})
	if err != nil {
		testsuit.T.Fail(err)
	}

	log.Println("Deleted EventListener")

	// Verify the EventListener's Deployment is deleted
	err = wait.WaitFor(c.Ctx, wait.DeploymentNotExist(c, namespace, fmt.Sprintf("%s-%s", eventReconciler.GeneratedResourcePrefix, elName)))
	if err != nil {
		testsuit.T.Fail(err)
	}

	log.Println("EventListener's Deployment was deleted")

	// Verify the EventListener's Service is deleted
	err = wait.WaitFor(c.Ctx, wait.ServiceNotExist(c, namespace, fmt.Sprintf("%s-%s", eventReconciler.GeneratedResourcePrefix, elName)))
	if err != nil {
		testsuit.T.Fail(err)
	}

	log.Println("EventListener's Service was deleted")

	//Delete Route exposed earlier
	err = c.Route.Routes(namespace).Delete(c.Ctx, fmt.Sprintf("%s-%s", eventReconciler.GeneratedResourcePrefix, elName), metav1.DeleteOptions{})
	if err != nil {
		testsuit.T.Fail(err)
	}

	// Verify the EventListener's Route is deleted
	err = wait.WaitFor(c.Ctx, wait.RouteNotExist(c, namespace, fmt.Sprintf("%s-%s", eventReconciler.GeneratedResourcePrefix, elName)))
	if err != nil {
		testsuit.T.Fail(err)
	}
	log.Println("EventListener's Route got deleted successfully...")

	// This is required when EL runs as TLS
	cmd.MustSucceed("rm", "-rf", os.Getenv("GOPATH")+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs")
}

func GetRoute(elname, namespace string) string {
	route := cmd.MustSucceed("oc", "-n", namespace, "get", "route", "--selector=eventlistener="+elname, "-o", "jsonpath='{range .items[*]}{.metadata.name}'").Stdout()

	route_url := cmd.MustSucceed("oc", "-n", namespace, "get", "route", strings.Trim(route, "'"), "--template='http://{{.spec.host}}'").Stdout()
	log.Printf("Route url: %s", route_url)

	time.Sleep(5 * time.Second)
	return strings.Trim(route_url, "'")
}
