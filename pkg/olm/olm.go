package olm

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
)

func CreateSubscriptionYaml(channel, installPlan, csv string) {
	// TODO: convert to Go code
	var data = struct {
		Channel string
		CSV     string
	}{
		Channel: channel,
		CSV:     csv,
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

	if err = ioutil.WriteFile(config.TempFile("subscription.yaml"), buffer.Bytes(), 0666); err != nil {
		testsuit.T.Fail(err)
	}

}

// DeleteCSV deletes cluster Service Version(v0.9.*) Resource from TargetNamespace
func DeleteCSV(version string) {
	log.Printf("output: %s\n",
		cmd.MustSucceed(
			"oc", "delete", "csv", "openshift-pipelines-operator."+version,
			"-n", "openshift-operators",
		).Stdout())
}

// Subscribe helps you to subscribe specific version pipelines operator from canary channel to OCP cluster
func Subscribe() {
	path := config.TempFile("subscription.yaml")
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "apply", "-f", path).Stdout())
}

// Unsubscribe helps you to subscribe specific version pipelines operator from canary channel to OCP cluster
func Unsubscribe() {
	path := config.TempFile("subscription.yaml")
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", "-f", path).Stdout())
}
