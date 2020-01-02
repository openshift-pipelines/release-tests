package operator

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/openshift-pipelines/release-tests/pkg/client"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	op "github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/names"
	"gotest.tools/v3/icmd"
	knativetest "knative.dev/pkg/test"
)

var CR *v1alpha1.Config
var Clt *client.Clients
var t = &testing.T{}

func init() {
	if os.Getenv("OPERATOR_VERSION") == "" {
		log.Println("\"OPERATOR_VERSION\" env variable is required, Cannot Procced E2E Tests")
		os.Exit(0)
	}

	Clt = client.NewClients(knativetest.Flags.Kubeconfig, knativetest.Flags.Cluster, config.DefaultTargetNs)
}

func ValidateClusterCR(c *client.Clients) *v1alpha1.Config {
	return helper.WaitForClusterCR(t, config.ClusterCRName, c)
}

func ValidatePipelineAndTriggerSetup(c *client.Clients, cr *v1alpha1.Config, webhookname string, controllername string) {
	helper.ValidatePipelineSetup(t,
		cr,
		c,
		webhookname, controllername)
	CR = ValidateClusterCR(c)
	if code := CR.Status.Conditions[0].Code; code != op.InstalledStatus {
		t.Errorf("Expected code to be %s but got %s", op.InstalledStatus, code)
	}
}

func VerifyPipelineVersion(c *client.Clients, version string) {
	CR = ValidateClusterCR(c)
	So(CR.Status.Conditions[0].Version, ShouldStartWith, version)
}

func ValidateSCC(c *client.Clients, targetNamespace string, PipelineControllerName string) {
	helper.ValidateSCCAdded(t, targetNamespace, PipelineControllerName, c)
}

func SubscribeToChannel() {
	res := icmd.RunCmd(icmd.Cmd{Command: append([]string{"oc"}, "apply", "-f", fmt.Sprintf("%s/src/github.com/openshift-pipelines/release-tests/config/subscription.yaml", os.Getenv("GOPATH"))), Timeout: 10 * time.Minute})
	res.Assert(t, icmd.Expected{
		ExitCode: 0,
		Err:      icmd.None,
	})
	log.Printf("%s\n", res.Stdout())
}

func ValidateOperatorInstall(c *client.Clients) {
	log.Printf("Waiting for operator to be up and running....\n")
	CR = ValidateClusterCR(c)
	helper.ValidatePipelineSetup(t,
		CR,
		c,
		config.PipelineWebhookName, config.PipelineControllerName)
	log.Printf("Operator is up\n")

}

func Setup(t *testing.T) (*client.Clients, string) {
	t.Helper()
	namespace := names.SimpleNameGenerator.RestrictLengthWithRandomSuffix("testrelease")
	Clt = client.NewClients(knativetest.Flags.Kubeconfig, knativetest.Flags.Cluster, namespace)
	SubscribeToChannel()
	ValidateOperatorInstall(Clt)
	helper.CreateNamespace(namespace, Clt.KubeClient)
	helper.VerifyServiceAccountExistence(namespace, Clt.KubeClient)
	return Clt, namespace
}
