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
	"os"

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

func getServiceList(c *clients.Clients, elname, namespace string) string {
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
	return serviceList.Items[0].Name
}

func ExposeEventListner(c *clients.Clients, elname, namespace string) string {
	cmd.MustSucceed("oc", "expose", "service", getServiceList(c, elname, namespace), "-n", namespace)

	route := cmd.MustSucceed("oc", "-n", namespace, "get", "route", "--selector=eventlistener="+elname, "-o", "jsonpath='{range .items[*]}{.metadata.name}'").Stdout()

	route_url := cmd.MustSucceed("oc", "-n", namespace, "get", "route", strings.Trim(route, "'"), "--template='http://{{.spec.host}}'").Stdout()
	log.Printf("Route url: %s", route_url)

	time.Sleep(5 * time.Second)
	return strings.Trim(route_url, "'")
}

func ExposeEventListnerForTLS(c *clients.Clients, elname, namespace string) string {
	svcName := getServiceList(c, elname, namespace)
	domain := getDomain(svcName, namespace)
	fmt.Println("domain values are", domain)
	gopath := os.Getenv("GOPATH")
	fmt.Println("gopathgopathgopathgopath", gopath)
	rootcaKey := gopath+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs/rootCA.key"
	rootcaCert := gopath+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs/rootCA.crt"
	tlsKey := gopath+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs/tls.key"
	tlsCert := gopath+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs/tls.crt"
	tlsCsr := gopath+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs/tls.csr"

	cmd.MustSucceed("openssl", "genrsa", "-out", rootcaKey, "4096").Stdout()

	cmd.MustSucceed("openssl", "req", "-x509", "-new", "-nodes", "-key", rootcaKey,
		"-sha256", "-days", "1024", "-out", rootcaCert, "-subj",
		"/C=IN/ST=Kar/L=Blr/O=RedHat/CN=client").Stdout()

	//cmd.MustSucceed("openssl", "req", "-new", "-x509", "-key", rootcaKey,
	//	"-days", "1024", "-out", rootcaCert, "-config",
	//	gopath+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs/csr_ca.txt").Stdout()

	cmd.MustSucceed("openssl", "genrsa", "-out", tlsKey, "4096").Stdout()

	cmd.MustSucceed("openssl", "req", "-new", "-key", tlsKey,
		"-subj", "/C=IN/ST=Kar/L=Blr/O=RedHat/CN=tls.test.apps.savita47new.tekton.codereadyqe.com",
		"-addext", "subjectAltName=DNS:apps.savita47new.tekton.codereadyqe.com",
		//"-config", gopath+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs/tls_answer.txt",
		"-out", tlsCsr).Stdout()
		//"-subj", "/C=IN/ST=Kar/L=Blr/O=RedHat/CN=tls.test.apps.savita47new.tekton.codereadyqe.com").Stdout()
		//"-subj", fmt.Sprintf("/C=IN/ST=Kar/L=Blr/O=RedHat/CN=%s",domain)).Stdout()


	cmd.MustSucceed("openssl", "x509", "-req", "-in", tlsCsr,
		"-CA", rootcaCert, "-CAkey",
		rootcaKey, "-CAcreateserial", "-out",
		tlsCert,
		"-days", "1024").Stdout()
		//"-days", "1024", "-extensions", "req_ext", "-extfile", gopath+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs/tls_answer.txt").Stdout()
		//"-days", "1024", "-sha256", "-extensions", "req_ext").Stdout()

	routeName := cmd.MustSucceed("oc", "create", "route", "reencrypt", "--ca-cert="+rootcaCert,
		"--cert="+tlsCert, "--key="+tlsKey,
		"--service="+svcName, "--hostname=tls.test.apps.savita47new.tekton.codereadyqe.com", "--port=listener", "-n", namespace).Stdout()
		//"--service="+svcName, "--hostname="+domain, "--port=listener", "-n", namespace).Stdout()

	fmt.Println("routename is", routeName)
	r := strings.Split(routeName, " ")
	fmt.Println("r", r, "*********", r[0])
	//route_url := cmd.MustSucceed("oc", "-n", namespace, "get", strings.Trim(routeName, " "), "--template='http://{{.spec.host}}'").Stdout()
	route_url := cmd.MustSucceed("oc", "-n", namespace, "get", r[0], "--template='http://{{.spec.host}}'").Stdout()
	log.Printf("Route url: %s", route_url)

	time.Sleep(5 * time.Second)
	return strings.Trim(route_url, "'")
}

func getDomain(svcName, namespace string) string {
	currentContext := cmd.MustSucceed("oc", "config", "current-context").Stdout()
	splittedValue := strings.Split(currentContext, "/")
	splittedValue = strings.Split(splittedValue[1], ":")
	splittedValue = strings.SplitAfter(splittedValue[0], "api-")
	fmt.Println("splittedValuesplittedValue", splittedValue, splittedValue[0], splittedValue[1])
	splittedValue = strings.Split(splittedValue[1], "-")
	joinedString := strings.Join(splittedValue, ".")
	routeDomainName := "apps."+joinedString
	return "tls.test."+routeDomainName
	//return "tls."+namespace+"."+routeDomainName
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

func MockPostEvent(routeurl, interceptor, eventType, payload string, isTLS bool) *http.Response {
	var (
		req *http.Request
		err error
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
}
