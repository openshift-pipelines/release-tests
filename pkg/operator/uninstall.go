package operator

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/openshift-pipelines/release-tests/pkg/client"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"gotest.tools/v3/icmd"
)

func DeleteClusterCR() {
	res := icmd.RunCmd(icmd.Cmd{Command: append([]string{"oc"}, "delete", "config.operator.tekton.dev", "cluster"), Timeout: 10 * time.Minute})
	res.Assert(t, icmd.Expected{
		ExitCode: 0,
		Err:      icmd.None,
	})
}

func DeleteCSV(version string) {
	res := icmd.RunCmd(icmd.Cmd{Command: append([]string{"oc"}, "delete", "csv", "openshift-pipelines-operator."+version, "-n", "openshift-operators"), Timeout: 10 * time.Minute})
	res.Assert(t, icmd.Expected{
		ExitCode: 0,
		Err:      icmd.None,
	})
	log.Printf("%s\n", res.Stdout())

}

func DeleteInstallPlan() {
	installPlanName := icmd.RunCmd(icmd.Cmd{Command: append([]string{"oc"}, "get", "subscription",
		"openshift-pipelines-operator", "-n", "openshift-operators",
		`-o=jsonpath={.status.installplan.name}`),
		Timeout: 10 * time.Minute}).Stdout()
	res := icmd.RunCmd(icmd.Cmd{Command: append([]string{"oc"}, "delete",
		"-n",
		"openshift-operators",
		"installplan",
		installPlanName),
		Timeout: 10 * time.Minute})
	res.Assert(t, icmd.Expected{
		ExitCode: 0,
		Err:      icmd.None,
	})
	log.Printf("Deleted install plan %s\n", installPlanName)

}

func DeleteSubscription() {
	res := icmd.RunCmd(icmd.Cmd{Command: append([]string{"oc"}, "delete",
		"subscription",
		"openshift-pipelines-operator",
		"-n",
		"openshift-operators"),
		Timeout: 10 * time.Minute})
	res.Assert(t, icmd.Expected{
		ExitCode: 0,
		Err:      icmd.None,
	})
	log.Printf("Deleted Subscription %s", res.Stdout())
}

func DeleteOperator(t *testing.T, c *client.Clients) {

	cr := helper.WaitForClusterCR(t, config.ClusterCRName, c)

	helper.DeleteClusterCR(t, config.ClusterCRName, c)

	helper.ValidatePipelineOrTriggerCleanup(t, cr,
		c,
		config.PipelineControllerName,
		config.PipelineWebhookName)

	helper.ValidatePipelineOrTriggerCleanup(t, cr,
		c,
		config.TriggerControllerName,
		config.TriggerWebhookName)

	helper.ValidateSCCRemoved(t, cr.Spec.TargetNamespace, config.PipelineControllerName, c)
	DeleteCSV(os.Getenv("OPERATOR_VERSION"))
	DeleteInstallPlan()
	DeleteSubscription()
}
