package config

import (
	"flag"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
)

const (
	// APIRetry defines the frequency at which we check for updates against the
	// k8s api when waiting for a specific condition to be true.
	APIRetry = time.Second * 5

	// APITimeout defines the amount of time we should spend querying the k8s api
	// when waiting for a specific condition to be true.
	APITimeout = time.Minute * 10
	// CLITimeout defines the amount of maximum execution time for CLI commands
	CLITimeout = time.Second * 15

	// ConsistentlyDuration sets  the default duration for Consistently. Consistently will verify that your condition is satisfied for this long.
	ConsistentlyDuration = 30 * time.Second

	ResourceTimeout = 60 * time.Second

	//TargetNamespace specify the name of Target namespace
	TargetNamespace = "openshift-pipelines"

	// Name of the pipeline controller deployment
	PipelineControllerName = "tekton-pipelines-controller"
	PipelineControllerSA   = "tekton-pipelines-controller"

	PipelineWebhookName          = "tekton-pipelines-webhook"
	PipelineWebhookConfiguration = "webhook.tekton.dev"
	SccAnnotationKey             = "operator.tekton.dev"

	// Name of the trigger deployment
	TriggerControllerName = "tekton-triggers-controller"
	TriggerWebhookName    = "tekton-triggers-webhook"

	// Name of the chains deployment
	ChainsControllerName = "tekton-chains-controller"

	// Name of the hub deployment
	HubApiName = "tekton-hub-api"
	HubDbName  = "tekton-hub-db"
	HubUiName  = "tekton-hub-ui"

	// Default config for auto pruner
	PrunerSchedule   = "0 8 * * *"
	PrunerNamePrefix = "tekton-resource-pruner-"

	// Name of PAC deployment
	PacControllerName = "pipelines-as-code-controller"
	PacWatcherName    = "pipelines-as-code-watcher"
	PacWebhookName    = "pipelines-as-code-webhook"

	// Name of tkn deployment
	TknDeployment = "tkn-cli-serve"

	// Name of console deployment
	ConsolePluginDeployment = "pipelines-console-plugin"

	// Community Clustertasks
	CommunityClustertasks = "jib-maven,helm-upgrade-from-source,helm-upgrade-from-repo,trigger-jenkins-job,git-cli,pull-request,kubeconfig-creator,argocd-task-sync-and-wait"

	// A token used in triggers tests
	TriggersSecretToken = "1234567"
)

// Name prefixes of installerset
var TektonInstallersetNamePrefixes [27]string = [27]string{
	"addon-custom-clustertask",
	"addon-custom-communityclustertask",
	"addon-custom-consolecli",
	"addon-custom-openshiftconsole",
	"addon-custom-pipelinestemplate",
	"addon-custom-triggersresources",
	"addon-versioned-clustertasks",
	"chain",
	"chain-secret",
	"console-link-hub",
	"openshiftpipelinesascode-main-deployment",
	"openshiftpipelinesascode-main-static",
	"openshiftpipelinesascode-post",
	"pipeline-main-deployment",
	"pipeline-main-static",
	"pipeline-post",
	"pipeline-pre",
	"rhosp-rbac",
	"tekton-hub-api",
	"tekton-hub-db",
	"tekton-hub-db-migration",
	"tekton-hub-ui",
	"tektoncd-pruner",
	"trigger-main-deployment",
	"trigger-main-static",
	"validating-mutating-webhook",
	"tekton-config-console-plugin-manifests",
}

var PrefixesOfDefaultPipelines [9]string = [9]string{"buildah", "s2i-dotnet", "s2i-go", "s2i-java", "s2i-nodejs", "s2i-perl", "s2i-php", "s2i-python", "s2i-ruby"}

// Flags holds the command line flags or defaults for settings in the user's environment.
// See EnvironmentFlags for a list of supported fields
// Todo: change initialization of falgs when required by parsing them or from environment variable
var Flags = initializeFlags()

// EnvironmentFlags define the flags that are needed to run the e2e tests.
type EnvironmentFlags struct {
	Cluster          string // K8s cluster (defaults to cluster in kubeconfig)
	Kubeconfig       string // Path to kubeconfig (defaults to ./kube/config)
	DockerRepo       string // Docker repo (defaults to $KO_DOCKER_REPO)
	CSV              string // Default csv openshift-pipelines-operator.v0.9.1
	Channel          string // Default channel canary
	CatalogSource    string
	SubscriptionName string
	InstallPlan      string // Default Installationplan Automatic
	OperatorVersion  string
	TknVersion       string
	ClusterArch      string // Architecture of the cluster
	IsDisconnected   bool
}

func initializeFlags() *EnvironmentFlags {
	var f EnvironmentFlags
	flag.StringVar(&f.Cluster, "cluster", "",
		"Provide the cluster to test against. Defaults to the current cluster in kubeconfig.")

	var defaultKubeconfig string
	if os.Getenv("KUBECONFIG") != "" {
		defaultKubeconfig = os.Getenv("KUBECONFIG")
	} else if usr, err := user.Current(); err == nil {
		defaultKubeconfig = path.Join(usr.HomeDir, ".kube/config")
	}

	flag.StringVar(&f.Kubeconfig, "kubeconfig", defaultKubeconfig,
		"Provide the path to the `kubeconfig` file you'd like to use for these tests. The `current-context` will be used.")

	defaultRepo := os.Getenv("KO_DOCKER_REPO")
	flag.StringVar(&f.DockerRepo, "dockerrepo", defaultRepo,
		"Provide the uri of the docker repo you have uploaded the test image to using `uploadtestimage.sh`. Defaults to $KO_DOCKER_REPO")

	defaultChannel := os.Getenv("CHANNEL")
	flag.StringVar(&f.Channel, "channel", defaultChannel,
		"Provide channel to subcribe your operator you'd like to use for these tests. By default `canary` will be used.")

	defaultCatalogSource := os.Getenv("CATALOG_SOURCE")
	flag.StringVar(&f.CatalogSource, "catalogsource", defaultCatalogSource,
		"Provide defaultCatalogSource to subscribe operator from. By default `custom-operators` will be used.")

	defaultSubscriptionName := os.Getenv("SUBSCRIPTION_NAME")
	flag.StringVar(&f.SubscriptionName, "subscriptionName", defaultSubscriptionName,
		"Provide defaultSubscriptionName to operator, By default `openshift-pipelines-operator-rh` will be used.")

	defaultPlan := os.Getenv("INSTALL_PLAN")
	flag.StringVar(&f.InstallPlan, "installplan", defaultPlan,
		"Provide Install Approval plan for your operator you'd like to use for these tests. By default `Automatic` will be used.")

	defaultOpVersion := os.Getenv("CSV_VERSION")
	flag.StringVar(&f.OperatorVersion, "opversion", defaultOpVersion,
		"Provide Operator version for your operator you'd like to use for these tests. By default `v0.9.1` ")

	defaultCsv := os.Getenv("CSV")
	flag.StringVar(&f.CSV, "csv", defaultCsv+defaultOpVersion,
		"Provide csv for your operator you'd like to use for these tests. By default `openshift-pipelines-operator.v0.9.1` will be used.")

	defaultTkn := os.Getenv("TKN_VERSION")
	flag.StringVar(&f.TknVersion, "tknversion", defaultTkn,
		"Provide tknversion to download specified cli binary you'd like to use for these tests. By default `0.6.0` will be used.")

	defaultClusterArch := os.Getenv("ARCH")
	if defaultClusterArch != "" && strings.Contains(defaultClusterArch, "/") {
		defaultClusterArch = strings.Split(defaultClusterArch, "/")[1]
	}
	flag.StringVar(&f.ClusterArch, "clusterarch", defaultClusterArch,
		"Provide the architecture of testing cluster. By default `amd64` will be used.")

	isDiconnectedEnv := os.Getenv("IS_DISCONNECTED")
	defaultIsDiconnected, err := strconv.ParseBool(isDiconnectedEnv)
	if err != nil {
		defaultIsDiconnected = false
	}
	flag.BoolVar(&f.IsDisconnected, "isdisconnected", defaultIsDiconnected,
		"Provide the info if the testing cluster is disconnected. By default `false` will be used.")
	return &f
}

func Dir() string {
	_, b, _, _ := runtime.Caller(0)
	configDir := path.Join(path.Dir(b), "..", "..", "template")
	return configDir
}

func File(elem ...string) string {
	path := append([]string{Dir()}, elem...)
	return filepath.Join(path...)
}

func Read(path string) ([]byte, error) {
	return os.ReadFile(File(path))
}

func TempDir() (string, error) {
	tmp := filepath.Join(Dir(), "..", "tmp")
	if _, err := os.Stat(tmp); os.IsNotExist(err) {
		err := os.Mkdir(tmp, 0755)
		return tmp, err
	}
	return tmp, nil
}

func TempFile(elem ...string) (string, error) {
	tmp, err := TempDir()
	if err != nil {
		return "", err
	}
	path := append([]string{tmp}, elem...)
	return filepath.Join(path...), nil
}

func RemoveTempDir() {
	var err error
	tmp, _ := TempDir()
	err = os.RemoveAll(tmp)
	if err != nil {
		testsuit.T.Errorf("Error: In deleting directory %s: %+v ", tmp, err)
	}
}

func Path(elem ...string) string {
	td := filepath.Join(Dir(), "..")
	if _, err := os.Stat(td); os.IsNotExist(err) {
		testsuit.T.Errorf("Error: in identifying test data path %s: %+v", td, err)
	}
	return filepath.Join(append([]string{td}, elem...)...)
}
