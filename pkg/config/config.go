package config

import (
	"flag"
	"os"
	"os/user"
	"path"
	"time"
)

const (
	// APIRetry defines the frequency at which we check for updates against the
	// k8s api when waiting for a specific condition to be true.
	APIRetry = time.Second * 5

	// APITimeout defines the amount of time we should spend querying the k8s api
	// when waiting for a specific condition to be true.
	APITimeout = time.Minute * 20

	// CleanupRetry is the interval at which test framework attempts cleanup
	CleanupRetry = time.Second * 10

	// CleanupTimeout is the wait time for test framework cleanup
	CleanupTimeout = time.Second * 180

	// TestOperatorName specifies the name of the operator being tested
	TestOperatorName = "openshift-pipelines-operator"

	DefaultSA            = "pipeline"
	DefaultIgnorePattern = "^(openshift|kube)-"

	ClusterCRName   = "cluster"
	DefaultTargetNs = "openshift-pipelines"

	// Name of the pipeline controller deployment
	PipelineControllerName = "tekton-pipelines-controller"
	PipelineControllerSA   = "tekton-pipelines-controller"

	PipelineWebhookName          = "tekton-pipelines-webhook"
	PipelineWebhookConfiguration = "webhook.tekton.dev"
	SccAnnotationKey             = "operator.tekton.dev"

	// Name of the trigger deployment
	TriggerControllerName = "tekton-triggers-controller"
	TriggerWebhookName    = "tekton-triggers-webhook"

	//Interval time defines, frequency at which it checks for resources
	Interval = 1 * time.Second

	//Timeout defines the amount of time we should spend waiting for the resource when condition is true
	Timeout = 10 * time.Minute

	// ConsistentlyDuration sets  the default duration for Consistently. Consistently will verify that your condition is satisfied for this long.
	ConsistentlyDuration = 30 * time.Second
)

// Flags holds the command line flags or defaults for settings in the user's environment.
// See EnvironmentFlags for a list of supported fields.
var Flags = initializeFlags()

// EnvironmentFlags define the flags that are needed to run the e2e tests.
type EnvironmentFlags struct {
	Cluster    string // K8s cluster (defaults to cluster in kubeconfig)
	Kubeconfig string // Path to kubeconfig (defaults to ./kube/config)
	DockerRepo string // Docker repo (defaults to $KO_DOCKER_REPO)
}

func initializeFlags() *EnvironmentFlags {
	var f EnvironmentFlags
	flag.StringVar(&f.Cluster, "cluster", "",
		"Provide the cluster to test against. Defaults to the current cluster in kubeconfig.")

	var defaultKubeconfig string
	if usr, err := user.Current(); err == nil {
		defaultKubeconfig = path.Join(usr.HomeDir, ".kube/config")
	}

	flag.StringVar(&f.Kubeconfig, "kubeconfig", defaultKubeconfig,
		"Provide the path to the `kubeconfig` file you'd like to use for these tests. The `current-context` will be used.")

	defaultRepo := os.Getenv("KO_DOCKER_REPO")
	flag.StringVar(&f.DockerRepo, "dockerrepo", defaultRepo,
		"Provide the uri of the docker repo you have uploaded the test image to using `uploadtestimage.sh`. Defaults to $KO_DOCKER_REPO")

	return &f
}
