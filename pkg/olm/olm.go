package olm

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
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

// DeleteCSV deletes cluster Service Version(v0.9.*) Resource from TargetNamespace
func DeleteCSV(version string) {
	log.Printf("output: %s\n",
		cmd.AssertOutput(
			&cmd.Cmd{
				Args: []string{
					"oc", "delete", "csv",
					"openshift-pipelines-operator." + version,
					"-n", "openshift-operators"},
				Expected: icmd.Success,
			}).Stdout())
}

// Subscribe helps you to subscribe specific version pipelines operator from canary channel to OCP cluster
func Subscribe(version string) {
	path := filepath.Join(helper.RootDir(), "../config/subscription.yaml")
	log.Printf("output: %s\n",
		cmd.AssertOutput(
			&cmd.Cmd{
				Args:     []string{"oc", "apply", "-f", path},
				Expected: icmd.Success,
			}))
}

// Unsubscribe helps you to subscribe specific version pipelines operator from canary channel to OCP cluster
func Unsubscribe(version string) {
	path := filepath.Join(helper.RootDir(), "../config/subscription.yaml")
	log.Printf("output: %s\n",
		cmd.AssertOutput(
			&cmd.Cmd{
				Args:     []string{"oc", "delete", "-f", path},
				Expected: icmd.Success,
			}))
}
