package olm

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"

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
func DeleteCSV(version string) {
	log.Printf("output: %s\n",
		cmd.MustSucceed(
			"oc", "delete", "csv", "openshift-pipelines-operator."+version,
			"-n", "openshift-operators",
		).Stdout())
}

// Subscribe helps you to subscribe specific version pipelines operator from canary channel to OCP cluster
func Subscribe() {
	path, err := config.TempFile("subscription.yaml")
	assert.FailOnError(err)
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "apply", "-f", path).Stdout())
}

// Unsubscribe helps you to subscribe specific version pipelines operator from canary channel to OCP cluster
func Unsubscribe() {
	path, err := config.TempFile("subscription.yaml")
	if err != nil {
		log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", "subscription", "openshift-pipelines-operator", "-n", "openshift-operators").Stdout())
	} else {
		log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", "-f", path).Stdout())
	}
}
