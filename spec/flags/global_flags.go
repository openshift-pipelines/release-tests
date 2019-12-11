package flags

import (
	"github.com/openshift-pipelines/release-tests/pkg/client"
)

var (
	Clients         *client.Clients
	Namespace       string
	Cleanup         func()
	CleanupSuite    func()
	OperatorVersion = "v0.9.1"
	PipelineVersion = "v0.9"
)
