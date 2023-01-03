package clients

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	configV1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	consolev1 "github.com/openshift/client-go/console/clientset/versioned/typed/console/v1"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	olmversioned "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned"
	monclientv1 "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	"github.com/tektoncd/operator/pkg/client/clientset/versioned"
	operatorv1alpha1 "github.com/tektoncd/operator/pkg/client/clientset/versioned/typed/operator/v1alpha1"
	pversioned "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1beta1"
	resourceversioned "github.com/tektoncd/pipeline/pkg/client/resource/clientset/versioned"
	resourcev1alpha1 "github.com/tektoncd/pipeline/pkg/client/resource/clientset/versioned/typed/resource/v1alpha1"
	triggersclientset "github.com/tektoncd/triggers/pkg/client/clientset/versioned"
)

// KubeClient holds instances of interfaces for making requests to kubernetes client.
type KubeClient struct {
	Kube *kubernetes.Clientset
}

// Clients holds instances of interfaces for making requests to Tekton Pipelines.
type Clients struct {
	KubeClient             *KubeClient
	Ctx                    context.Context
	Dynamic                dynamic.Interface
	Operator               operatorv1alpha1.OperatorV1alpha1Interface
	KubeConfig             *rest.Config
	Scheme                 *runtime.Scheme
	OLM                    olmversioned.Interface
	Route                  routev1.RouteV1Interface
	ProxyConfig            configV1.ConfigV1Interface
	ClusterVersion         configV1.ClusterVersionInterface
	ConsoleCLIDownload     consolev1.ConsoleCLIDownloadInterface
	MonitoringClient       monclientv1.MonitoringV1Interface
	Tekton                 pversioned.Interface
	PipelineClient         v1beta1.PipelineInterface
	TaskClient             v1beta1.TaskInterface
	TaskRunClient          v1beta1.TaskRunInterface
	PipelineRunClient      v1beta1.PipelineRunInterface
	PipelineResourceClient resourcev1alpha1.PipelineResourceInterface
	TriggersClient         triggersclientset.Interface
	ClustertaskClient      v1beta1.ClusterTaskInterface
}

// NewClients instantiates and returns several clientsets required for making request to the
// TektonPipeline cluster specified by the combination of clusterName and configPath.
func NewClients(configPath string, clusterName, namespace string) (*Clients, error) {
	var err error
	clients := &Clients{}

	clients.KubeClient, clients.KubeConfig, err = NewKubeClient(configPath, clusterName)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubeclient from config file at %s: %s", configPath, err)
	}

	// We poll, so set our limits high.
	clients.KubeConfig.QPS = 100
	clients.KubeConfig.Burst = 200

	ctx := context.Background()
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()
	clients.Ctx = ctx

	clients.Dynamic, err = dynamic.NewForConfig(clients.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create dynamic clients from config file at %s: %s", configPath, err)
	}

	clients.Operator, err = newTektonOperatorAlphaClients(clients.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Operator v1alpha1 clients from config file at %s: %s", configPath, err)
	}

	clients.OLM, err = olmversioned.NewForConfig(clients.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create olm clients from config file at %s: %s", configPath, err)
	}

	clients.Tekton, err = pversioned.NewForConfig(clients.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pipeline clientset from config file at %s: %s", configPath, err)
	}

	rcs, err := resourceversioned.NewForConfig(clients.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create resource clientset from config file at %s: %s", configPath, err)
	}

	clients.TriggersClient, err = triggersclientset.NewForConfig(clients.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create triggers clientset from config file at %s: %s", configPath, err)
	}

	clients.PipelineClient = clients.Tekton.TektonV1beta1().Pipelines(namespace)
	clients.TaskClient = clients.Tekton.TektonV1beta1().Tasks(namespace)
	clients.TaskRunClient = clients.Tekton.TektonV1beta1().TaskRuns(namespace)
	clients.PipelineRunClient = clients.Tekton.TektonV1beta1().PipelineRuns(namespace)
	clients.PipelineResourceClient = rcs.TektonV1alpha1().PipelineResources(namespace)
	clients.Route = routev1.NewForConfigOrDie(clients.KubeConfig)
	clients.ProxyConfig = configV1.NewForConfigOrDie(clients.KubeConfig)
	clients.ClusterVersion = configV1.NewForConfigOrDie(clients.KubeConfig).ClusterVersions()
	clients.ConsoleCLIDownload = consolev1.NewForConfigOrDie(clients.KubeConfig).ConsoleCLIDownloads()
	clients.ClustertaskClient = clients.Tekton.TektonV1beta1().ClusterTasks()
	if err != nil {
		return nil, err
	}
	clients.MonitoringClient = monclientv1.NewForConfigOrDie(clients.KubeConfig)

	return clients, nil
}

// NewKubeClient instantiates and returns several clientsets required for making request to the
// kube client specified by the combination of clusterName and configPath. Clients can make requests within namespace.
func NewKubeClient(configPath string, clusterName string) (*KubeClient, *rest.Config, error) {
	cfg, err := BuildClientConfig(configPath, clusterName)
	if err != nil {
		return nil, nil, err
	}

	k, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, nil, err
	}
	return &KubeClient{Kube: k}, cfg, nil
}

// BuildClientConfig builds the client config specified by the config path and the cluster name
func BuildClientConfig(kubeConfigPath string, clusterName string) (*rest.Config, error) {
	overrides := clientcmd.ConfigOverrides{}
	// Override the cluster name if provided.
	if clusterName != "" {
		overrides.Context.Cluster = clusterName
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath},
		&overrides).ClientConfig()
}

func newTektonOperatorAlphaClients(cfg *rest.Config) (operatorv1alpha1.OperatorV1alpha1Interface, error) {
	cs, err := versioned.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return cs.OperatorV1alpha1(), nil
}

func (c *Clients) TektonPipeline() operatorv1alpha1.TektonPipelineInterface {
	return c.Operator.TektonPipelines()
}

func (c *Clients) TektonTrigger() operatorv1alpha1.TektonTriggerInterface {
	return c.Operator.TektonTriggers()
}

func (c *Clients) TektonDashboard() operatorv1alpha1.TektonDashboardInterface {
	return c.Operator.TektonDashboards()
}

func (c *Clients) TektonAddon() operatorv1alpha1.TektonAddonInterface {
	return c.Operator.TektonAddons()
}

func (c *Clients) TektonConfig() operatorv1alpha1.TektonConfigInterface {
	return c.Operator.TektonConfigs()
}
