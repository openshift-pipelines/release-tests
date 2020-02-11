package olm

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
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

	b, err := config.Read("subscription.yaml.tmp")
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("subscription").Parse(string(b))
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(&tmplBytes, data)
	if err != nil {
		// TODO fail testsuit
		panic(err)
	}

	err = ioutil.WriteFile(config.File("subscription.yaml"), tmplBytes.Bytes(), 0666)
	// TODO: handle this error
	if err != nil {
		// print it out
		log.Fatal(err)
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
func Subscribe(version string) {
	path := config.File("subscription.yaml")
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "apply", "-f", path))
}

// Unsubscribe helps you to subscribe specific version pipelines operator from canary channel to OCP cluster
func Unsubscribe(version string) {
	path := config.File("subscription.yaml")
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", "-f", path))
}
