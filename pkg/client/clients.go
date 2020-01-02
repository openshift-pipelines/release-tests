package client

import (
	"log"
	"sync"
	"time"

	goctx "context"

	"github.com/tektoncd/operator/pkg/apis"
	op "github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1alpha1"
	extscheme "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/discovery/cached"
	cgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	knativetest "knative.dev/pkg/test"
	dynclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	// Global framework struct
	Global *Clients
	// mutex for AddToFrameworkScheme
	mutex = sync.Mutex{}
	// decoder used by createFromYaml
	dynamicDecoder runtime.Decoder
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
	Delete(gCtx goctx.Context, obj runtime.Object, opts ...dynclient.DeleteOptionFunc) error
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
	return f.Client.List(gCtx, opts, list)
}

func (f *frameworkClient) Delete(gCtx goctx.Context, obj runtime.Object, opts ...dynclient.DeleteOptionFunc) error {
	return f.Client.Delete(gCtx, obj, opts...)
}

func (f *frameworkClient) Update(gCtx goctx.Context, obj runtime.Object) error {
	return f.Client.Update(gCtx, obj)
}

type Clients struct {
	Client                 *frameworkClient
	KubeClient             *knativetest.KubeClient
	KubeConfig             *rest.Config
	Scheme                 *runtime.Scheme
	PipelineClient         v1alpha1.PipelineInterface
	TaskClient             v1alpha1.TaskInterface
	TaskRunClient          v1alpha1.TaskRunInterface
	PipelineRunClient      v1alpha1.PipelineRunInterface
	PipelineResourceClient v1alpha1.PipelineResourceInterface
	ConditionClient        v1alpha1.ConditionInterface
}

// newClients instantiates and returns several clientsets required for making requests to the
// Pipeline cluster specified by the combination of clusterName and configPath. Clients can
// make requests within namespace.
func NewClients(configPath, clusterName, namespace string) *Clients {

	var err error
	c := &Clients{}

	c.KubeClient, err = knativetest.NewKubeClient(configPath, clusterName)
	if err != nil {
		log.Fatalf("failed to create kubeclient from config file at %s: %s", configPath, err)
	}

	c.KubeConfig, err = knativetest.BuildClientConfig(configPath, clusterName)
	if err != nil {
		log.Fatalf("failed to create configuration obj from %s for cluster %s: %s", configPath, clusterName, err)
	}

	scheme := runtime.NewScheme()
	if err := cgoscheme.AddToScheme(scheme); err != nil {
		log.Fatalf("failed to add cgo scheme to runtime scheme: (%v)", err)
	}
	if err := extscheme.AddToScheme(scheme); err != nil {
		log.Fatalf("failed to add api extensions scheme to runtime scheme: (%v)", err)
	}
	cachedDiscoveryClient := cached.NewMemCacheClient(c.KubeClient.Kube.Discovery())
	restMapper = restmapper.NewDeferredDiscoveryRESTMapper(cachedDiscoveryClient)
	restMapper.Reset()
	dynClient, err := dynclient.New(c.KubeConfig, dynclient.Options{Scheme: scheme, Mapper: restMapper})
	if err != nil {
		log.Fatalf("failed to build the dynamic client: %v", err)
	}
	serializer.NewCodecFactory(scheme).UniversalDeserializer()
	c.Scheme = scheme
	c.Client = &frameworkClient{Client: dynClient}

	cs, err := versioned.NewForConfig(c.KubeConfig)
	if err != nil {
		log.Fatalf("failed to create pipeline clientset from config file at %s: %s", configPath, err)
	}
	c.PipelineClient = cs.TektonV1alpha1().Pipelines(namespace)
	c.TaskClient = cs.TektonV1alpha1().Tasks(namespace)
	c.TaskRunClient = cs.TektonV1alpha1().TaskRuns(namespace)
	c.PipelineRunClient = cs.TektonV1alpha1().PipelineRuns(namespace)
	c.PipelineResourceClient = cs.TektonV1alpha1().PipelineResources(namespace)
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
		err = dynClient.List(goctx.TODO(), &dynclient.ListOptions{Namespace: "default"}, obj)
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
	dynamicDecoder = serializer.NewCodecFactory(c.Scheme).UniversalDeserializer()
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
