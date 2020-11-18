package olm

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"strings"

	"github.com/openshift-pipelines/release-tests/pkg/assert"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
)

func CreateSubscriptionYaml(channel string) {
	// TODO: convert to Go code
	var data = struct {
		Channel string
	}{
		Channel: channel,
	}

	if _, err := config.TempDir(); err != nil {
		testsuit.T.Fail(err)
	}

	b, err := config.Read("subscription.yaml.tmp")
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("subscription").Parse(string(b))
	if err != nil {
		testsuit.T.Fail(err)
	}

	var buffer bytes.Buffer
	if err = tmpl.Execute(&buffer, data); err != nil {
		testsuit.T.Fail(err)
	}
	file, err := config.TempFile("subscription.yaml")
	assert.FailOnError(err)
	if err = ioutil.WriteFile(file, buffer.Bytes(), 0666); err != nil {
		testsuit.T.Fail(err)
	}

}

// DeleteCSV deletes cluster Service Version(v0.9.*) Resource from TargetNamespace
func DeleteCSV() {
	CSV := cmd.MustSucceed("oc", "get", "csv", "-n", "openshift-operators", "-o", "name").Stdout()
	log.Printf("output: %s\n",
		cmd.MustSucceed("oc", "-n", "openshift-operators", "delete", strings.TrimSuffix(CSV, "\n")).Stdout())
}

// Subscribe helps you to subscribe specific version pipelines operator from canary channel to OCP cluster
func Subscribe() {
	path, err := config.TempFile("subscription.yaml")
	assert.FailOnError(err)
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "apply", "-f", path).Stdout())
}

// Unsubscribe helps you to subscribe specific version pipelines operator from canary channel to OCP cluster
func Unsubscribe() {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", "subscription", "openshift-pipelines-operator-rh", "-n", "openshift-operators").Stdout())
}

func DeleteInstallPlan() {
	InstallPlan := cmd.MustSucceed("oc", "get", "subscription", "openshift-pipelines-operator-rh", "-n", "openshift-operators", "-o", "jsonpath='{.status.installplan.name}'").Stdout()
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "-n", "openshift-operators", "delete", "installplan", strings.Trim(strings.TrimSuffix(InstallPlan, "\n"), "'")).Stdout())
}
