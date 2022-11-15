package chains

import (
	"bytes"
	"log"
	"os"
	"text/template"

	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
)

func CreateChainsCR(cs *clients.Clients, format, storage, transparencyEnabled string) {
	var chain = struct {
		ArtifactsTaskrunFormat  string
		ArtifactsTaskrunStorage string
		TargetNamespace         string
		TransparencyEnabled     string
	}{
		ArtifactsTaskrunFormat:  format,
		ArtifactsTaskrunStorage: storage,
		TargetNamespace:         config.TargetNamespace,
		TransparencyEnabled:     transparencyEnabled,
	}

	if _, err := config.TempDir(); err != nil {
		assert.FailOnError(err)
	}
	defer config.RemoveTempDir()

	tmpl, err := config.Read("chains.yaml.tmp")
	if err != nil {
		assert.FailOnError(err)
	}

	sub, err := template.New("chain").Parse(string(tmpl))
	if err != nil {
		assert.FailOnError(err)
	}

	var buffer bytes.Buffer
	if err = sub.Execute(&buffer, chain); err != nil {
		assert.FailOnError(err)
	}
	file, err := config.TempFile("chain.yaml")
	assert.FailOnError(err)
	if err = os.WriteFile(file, buffer.Bytes(), 0666); err != nil {
		assert.FailOnError(err)
	}

	log.Printf("output: %s\n", cmd.MustSucceed("oc", "apply", "-f", file).Stdout())
}
