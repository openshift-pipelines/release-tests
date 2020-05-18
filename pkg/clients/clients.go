package clients

import (
	"sync"
	"time"

	goctx "context"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/tektoncd/operator/pkg/apis"
	op "github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1beta1"
	resourceversioned "github.com/tektoncd/pipeline/pkg/client/resource/clientset/versioned"
	resourcev1alpha1 "github.com/tektoncd/pipeline/pkg/client/resource/clientset/versioned/typed/resource/v1alpha1"
	triggersclientset "github.com/tektoncd/triggers/pkg/client/clientset/versioned"
	extscheme "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	cached "k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	cgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	dynclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	// Global framework struct
	// mutex for AddToFrameworkScheme
	mutex = sync.Mutex{}

	// decoder used by createFromYaml
	//dynamicDecoder runtime.Decoder
	// restMapper for the dynamic client
	restMapper *restmapper.DeferredDiscoveryRESTMapper
)

type frameworkClient struct {
	dynclient.Client
}

var _ FrameworkClient = &frameworkClient{}

type FrameworkClient interface {
	Get(gCtx goctx.Context, key dynclient.ObjectKey, obj runtime.Object) error
	List(gCtx goctx.Context, opts *dynclient.ListOptions, list runtime.Object) error
	Create(gCtx goctx.Context, obj runtime.Object) error
	Delete(gCtx goctx.Context, obj runtime.Object, opts ...dynclient.DeleteOption) error
	Update(gCtx goctx.Context, obj runtime.Object) error
}

// Create uses the dynamic client to create an object and then adds a
// cleanup function to delete it when Cleanup is called. In addition to
// the standard controller-runtime client options
func (f *frameworkClient) Create(gCtx goctx.Context, obj runtime.Object) error {
	objCopy := obj.DeepCopyObject()
	err := f.Client.Create(gCtx, obj)
	if err != nil {
		return err
	}

	_, err1 := dynclient.ObjectKeyFromObject(objCopy)
	if err1 != nil {
		return err1
	}
	return nil
}

func (f *frameworkClient) Get(gCtx goctx.Context, key dynclient.ObjectKey, obj runtime.Object) error {
	return f.Client.Get(gCtx, key, obj)
}

func (f *frameworkClient) List(gCtx goctx.Context, opts *dynclient.ListOptions, list runtime.Object) error {
	return f.Client.List(gCtx, list, opts)
}

func (f *frameworkClient) Delete(gCtx goctx.Context, obj runtime.Object, opts ...dynclient.DeleteOption) error {
	return f.Client.Delete(gCtx, obj, opts...)
}

func (f *frameworkClient) Update(gCtx goctx.Context, obj runtime.Object) error {
	return f.Client.Update(gCtx, obj)
}

// KubeClient holds instances of interfaces for making requests to kubernetes client.
type KubeClient struct {
	Kube *kubernetes.Clientset
}

// Clients holds instances of interfaces for making requests to the Pipeline controllers.
type Clients struct {
	Client                 *frameworkClient
	KubeClient             *KubeClient
	KubeConfig             *rest.Config
	Scheme                 *runtime.Scheme
	Dynamic                dynamic.Interface
	Tekton                 versioned.Interface
	PipelineClient         v1beta1.PipelineInterface
	TaskClient             v1beta1.TaskInterface
	TaskRunClient          v1beta1.TaskRunInterface
	PipelineRunClient      v1beta1.PipelineRunInterface
	PipelineResourceClient resourcev1alpha1.PipelineResourceInterface
	ConditionClient        v1alpha1.ConditionInterface
	TriggersClient         triggersclientset.Interface
}

// NewKubeClient instantiates and returns several clientsets required for making request to the
// kube client specified by the combination of clusterName and configPath. Clients can make requests within namespace.
func NewKubeClient(configPath string, clusterName string) (*KubeClient, error) {
	cfg, err := BuildClientConfig(configPath, clusterName)
	if err != nil {
		return nil, err
	}

	k, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &KubeClient{Kube: k}, nil
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

// NewClients instantiates and returns several clientsets required for making requests to the
// Pipeline cluster specified by the combination of clusterName and configPath. Clients can
// make requests within namespace.
func NewClients(configPath, clusterName, namespace string) *Clients {

	var err error
	c := &Clients{}

	c.KubeClient, err = NewKubeClient(configPath, clusterName)
	if err != nil {
		testsuit.T.Errorf("failed to create kubeclient from config file at %s: %s", configPath, err)
	}

	c.KubeConfig, err = BuildClientConfig(configPath, clusterName)
	if err != nil {
		testsuit.T.Errorf("failed to create configuration obj from %s for cluster %s: %s", configPath, clusterName, err)
	}

	scheme := runtime.NewScheme()
	if err := cgoscheme.AddToScheme(scheme); err != nil {
		testsuit.T.Errorf("failed to add cgo scheme to runtime scheme: (%v)", err)
	}
	if err := extscheme.AddToScheme(scheme); err != nil {
		testsuit.T.Errorf("failed to add api extensions scheme to runtime scheme: (%v)", err)
	}
	cachedDiscoveryClient := cached.NewMemCacheClient(c.KubeClient.Kube.Discovery())
	restMapper = restmapper.NewDeferredDiscoveryRESTMapper(cachedDiscoveryClient)
	restMapper.Reset()
	dynClient, err := dynclient.New(c.KubeConfig, dynclient.Options{Scheme: scheme, Mapper: restMapper})
	if err != nil {
		testsuit.T.Errorf("failed to build the dynamic client: %v", err)
	}
	serializer.NewCodecFactory(scheme).UniversalDeserializer()
	c.Scheme = scheme
	c.Client = &frameworkClient{Client: dynClient}

	cs, err := versioned.NewForConfig(c.KubeConfig)
	if err != nil {
		testsuit.T.Errorf("failed to create pipeline clientset from config file at %s: %s", configPath, err)
	}
	c.Tekton = cs

	rcs, err := resourceversioned.NewForConfig(c.KubeConfig)
	if err != nil {
		testsuit.T.Errorf("Failed to create resource clientset from config file at %s: %s", configPath, err)
	}

	c.TriggersClient, err = triggersclientset.NewForConfig(c.KubeConfig)
	if err != nil {
		testsuit.T.Errorf("Failed to create triggers clientset from config file at %s: %s", configPath, err)
	}

	c.Dynamic, err = dynamic.NewForConfig(c.KubeConfig)
	if err != nil {
		testsuit.T.Errorf("Failed to create dynamic clients from config file at %s: %s", configPath, err)

	}

	c.PipelineClient = cs.TektonV1beta1().Pipelines(namespace)
	c.TaskClient = cs.TektonV1beta1().Tasks(namespace)
	c.TaskRunClient = cs.TektonV1beta1().TaskRuns(namespace)
	c.PipelineRunClient = cs.TektonV1beta1().PipelineRuns(namespace)
	c.PipelineResourceClient = rcs.TektonV1alpha1().PipelineResources(namespace)
	c.ConditionClient = cs.TektonV1alpha1().Conditions(namespace)
	c = initTestingFramework(c)
	return c
}

type addToSchemeFunc func(*runtime.Scheme) error

// AddToFrameworkScheme allows users to add the scheme for their custom resources
// to the framework's scheme for use with the dynamic client. The user provides
// the addToScheme function (located in the register.go file of their operator
// project) and the List struct for their custom resource. For example, for a
// memcached operator, the list stuct may look like:
// &MemcachedList{}
// The List object is needed because the CRD has not always been fully registered
// by the time this function is called. If the CRD takes more than 5 seconds to
// become ready, this function throws an error
func AddToFrameworkScheme(addToScheme addToSchemeFunc, obj runtime.Object, c *Clients) *Clients {
	mutex.Lock()
	defer mutex.Unlock()
	err := addToScheme(c.Scheme)
	if err != nil {
		return nil
	}
	restMapper.Reset()
	dynClient, err := dynclient.New(c.KubeConfig, dynclient.Options{Scheme: c.Scheme, Mapper: restMapper})
	if err != nil {
		return nil
	}
	err = wait.PollImmediate(time.Second, time.Second*10, func() (done bool, err error) {
		err = dynClient.List(goctx.TODO(), obj, &dynclient.ListOptions{Namespace: "default"})
		if err != nil {
			restMapper.Reset()
			return false, nil
		}
		c.Client = &frameworkClient{Client: dynClient}
		return true, nil
	})
	if err != nil {
		return nil
	}
	serializer.NewCodecFactory(c.Scheme).UniversalDeserializer()
	return c
}

func initTestingFramework(c *Clients) *Clients {
	apiVersion := "operator.tekton.dev/v1alpha1"
	kind := "Config"

	configList := &op.ConfigList{
		TypeMeta: metav1.TypeMeta{
			Kind:       kind,
			APIVersion: apiVersion,
		},
	}

	return AddToFrameworkScheme(apis.AddToScheme, configList, c)
}
