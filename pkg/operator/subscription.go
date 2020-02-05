package operator

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"gotest.tools/v3/icmd"
)

func CreateSubscriptionYaml(channel, installPlan, csv string) {
	// TODO: convert to Go code
	var err error
	var data = struct {
		Channel     string
		InstallPlan string
		CSV         string
	}{
		Channel:     channel,
		InstallPlan: installPlan,
		CSV:         csv,
	}

	var tmplBytes bytes.Buffer

	b, err := ioutil.ReadFile(filepath.Join(helper.RootDir(), "../config/subscription.yaml.tmp")) // just pass the file name
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("subscription").Parse(string(b))
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(&tmplBytes, data)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filepath.Join(helper.RootDir(), "../config/subscription.yaml"), tmplBytes.Bytes(), 0666)
	// TODO: handle this error
	if err != nil {
		// print it out
		log.Fatal(err)
	}
}

// // DeleteClusterCR deletes Cluster config from the cluster
// func DeleteClusterCR() {
//
// 	log.Printf("output: %s\n",
// 		cmd.AssertOutput(
// 			&cmd.Cmd{
// 				Args: []string{"oc", "delete", "config.operator.tekton.dev", "cluster"},
// 				Expected: icmd.Expected{
// 					ExitCode: 0,
// 					Err:      icmd.None,
// 				},
// 			}).Stdout())
//
// }

// DeleteCSV deletes cluster Service Version(v0.9.*) Resource from TargetNamespace
func DeleteCSV(version string) {
	log.Printf("output: %s\n",
		cmd.AssertOutput(
			&cmd.Cmd{
				Args: []string{"oc", "delete",
					"csv",
					"openshift-pipelines-operator." + version,
					"-n",
					"openshift-operators"},
				Expected: icmd.Success,
			}).Stdout())
}

// DeleteInstallPlan deletes installation plan
func DeleteInstallPlan() {

	installPlan := cmd.AssertOutput(
		&cmd.Cmd{
			Args: []string{"oc", "get", "-n", "openshift-operators",
				"subscription", "openshift-pipelines-operator",
				`-o=jsonpath={.status.installplan.name}`},
			Expected: icmd.Success,
		}).Stdout()

	log.Printf("install paln %s\n", installPlan)
	res := cmd.AssertOutput(
		&cmd.Cmd{
			Args: []string{"oc", "delete",
				"-n", "openshift-operators",
				"installplan",
				installPlan},
			Expected: icmd.Success,
		})
	log.Printf("Deleted install plan : %s\n", res.Stdout())
}

// Subscribe helps you to subscribe specific version pipelines operator from canary channel to OCP cluster
func Subscribe(version string) (*clients.Clients, string, func()) {
	path := filepath.Join(helper.RootDir(), "../config/subscription.yaml")
	log.Printf("output: %s\n",
		cmd.AssertOutput(
			&cmd.Cmd{
				Args:     []string{"oc", "apply", "-f", path},
				Expected: icmd.Success,
			}))

	cs, ns, cleanupNs := k8s.NewClientSet()

	cleanup := func() {
		Delete(cs, version)
		cleanupNs()
	}

	return cs, ns, cleanup
}

// DeleteSubscription deletes operator subscription from cluster
func DeleteSubscription() {
	log.Printf("Output %s \n", cmd.AssertOutput(
		&cmd.Cmd{
			Args: []string{
				"oc", "delete", "-n", "openshift-operators",
				"subscription", "openshift-pipelines-operator"},
			Expected: icmd.Success,
		}).Stdout())

}
