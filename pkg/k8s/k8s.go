package k8s

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/oc"
	"github.com/openshift-pipelines/release-tests/pkg/openshift"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	secv1 "github.com/openshift/api/security/v1"
	secclient "github.com/openshift/client-go/security/clientset/versioned/typed/security/v1"
	"github.com/tektoncd/pipeline/pkg/names"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
)

// NewClientSet is a setup function which helps you
// 1. to create clientSet instance to `client`
// 2. to create a namespace with random name
func NewClientSet() (*clients.Clients, string, func()) {
	// TODO: fix this; method is in k8s but returns client.Clients
	ns := names.SimpleNameGenerator.RestrictLengthWithRandomSuffix("releasetest")
	cs, err := clients.NewClients(config.Flags.Kubeconfig, config.Flags.Cluster, ns)
	if err != nil {
		testsuit.T.Fail(err)
	}

	oc.CreateNewProject(ns)
	return cs, ns, func() {
		oc.DeleteProjectIgnoreErors(ns)
	}
}

// WaitForDeploymentDeletion checks to see if a given deployment is deleted
// the function returns an error if the given deployment is not deleted within the timeout
func WaitForDeploymentDeletion(cs *clients.Clients, namespace, name string) error {
	err := wait.PollUntilContextTimeout(cs.Ctx, config.APIRetry, config.APITimeout, false, func(context.Context) (bool, error) {
		kc := cs.KubeClient.Kube
		_, err := kc.AppsV1().Deployments(namespace).Get(cs.Ctx, name, metav1.GetOptions{})
		if err != nil {
			if errors.IsGone(err) || errors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		log.Printf("Waiting for deletion of %s deployment\n", name)
		return false, nil
	})
	if err != nil {
		testsuit.T.Errorf("failed to delete deployment %s \n %v", name, err)
	}
	return err
}

// WaitForServiceAccount checks if service account created
func WaitForServiceAccount(cs *clients.Clients, ns, targetSA string) *corev1.ServiceAccount {
	ret := &corev1.ServiceAccount{}
	err := wait.PollUntilContextTimeout(cs.Ctx, config.APIRetry, config.APITimeout, false, func(context.Context) (bool, error) {
		saList, err := cs.KubeClient.Kube.CoreV1().ServiceAccounts(ns).List(cs.Ctx, metav1.ListOptions{})
		for _, sa := range saList.Items {
			if sa.Name == targetSA {
				ret = &sa
				return true, nil
			}
		}
		return false, err
	})
	if err != nil {
		testsuit.T.Errorf("service account %s not found in namespace %s", targetSA, ns)
	}
	return ret
}

func ValidateSCCAdded(cs *clients.Clients, ns, sa string) {
	err := wait.PollUntilContextTimeout(cs.Ctx, config.APIRetry, config.APITimeout, false, func(context.Context) (bool, error) {
		privileged, err := GetPrivilegedSCC(cs)
		if err != nil {
			log.Printf("failed to get privileged scc: %s \n", err)
			return false, err
		}
		log.Printf("... looking at %v", privileged.Users)

		ctrlSA := fmt.Sprintf("system:serviceaccount:%s:%s", ns, sa)
		return inList(privileged.Users, ctrlSA), nil
	})
	if err != nil {
		testsuit.T.Errorf("failed to add privileged SCC to the service account %s \n %v", sa, err)
	}
}

func ValidateSCCRemoved(cs *clients.Clients, ns, sa string) {
	err := wait.PollUntilContextTimeout(cs.Ctx, config.APIRetry, config.APITimeout, false, func(context.Context) (bool, error) {
		privileged, err := GetPrivilegedSCC(cs)
		if err != nil {
			log.Printf("failed to get privileged scc: %s \n", err)
			return false, err
		}
		ctrlSA := fmt.Sprintf("system:serviceaccount:%s:%s", ns, sa)
		return !inList(privileged.Users, ctrlSA), nil
	})
	if err != nil {
		testsuit.T.Errorf("failed to remove privileged SCC from service account %s \n %v", sa, err)
	}
}

func inList(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func ValidateDeployments(cs *clients.Clients, ns string, deployments ...string) {
	kc := cs.KubeClient.Kube
	for _, d := range deployments {
		err := WaitForDeployment(cs.Ctx, kc, ns,
			d,
			1,
			config.APIRetry,
			config.APITimeout,
		)
		if err != nil {
			testsuit.T.Errorf("failed to create deployment %+v \n %v", d, err)
		}
	}
}

func GetPrivilegedSCC(cs *clients.Clients) (*secv1.SecurityContextConstraints, error) {
	sec, err := secclient.NewForConfig(cs.KubeConfig)
	if err != nil {
		return nil, err
	}
	return sec.SecurityContextConstraints().Get(cs.Ctx, "privileged", metav1.GetOptions{})
}

func ValidateDeploymentDeletion(cs *clients.Clients, ns string, deployments ...string) {
	for _, d := range deployments {
		err := WaitForDeploymentDeletion(cs, ns, d)
		if err != nil {
			testsuit.T.Errorf("failed to delete deployment %+v \n %v", d, err)
		}
	}
}

func WaitForDeployment(ctx context.Context, kc kubernetes.Interface, namespace, name string, replicas int, retryInterval, timeout time.Duration) error {
	err := wait.PollUntilContextTimeout(ctx, retryInterval, timeout, false, func(context.Context) (done bool, err error) {
		deployment, err := kc.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.Printf("Waiting for availability of %s deployment\n", name)
				return false, nil
			}
			return false, err
		}

		if int(deployment.Status.AvailableReplicas) == replicas {
			return true, nil
		}
		log.Printf("Waiting for full availability of deployment %s (%d/%d)\n", name, deployment.Status.AvailableReplicas, replicas)
		return false, nil
	})
	return err
}

func VerifyNoServiceAccount(ctx context.Context, kc *clients.KubeClient, sa, ns string) {
	log.Printf("Verify SA %q is absent in namespace %q", sa, ns)
	if err := wait.PollUntilContextTimeout(ctx, config.APIRetry, config.APITimeout, true, func(context.Context) (done bool, err error) {
		_, err = kc.Kube.CoreV1().ServiceAccounts(ns).Get(ctx, sa, metav1.GetOptions{})
		if err == nil || !errors.IsNotFound(err) {
			return false, fmt.Errorf("sa %q exists in namespace %q", sa, ns)
		}
		return true, nil
	}); err != nil {
		testsuit.T.Errorf("Fail: SA %q exists in namespace %q, err: %s", sa, ns, err)
	}
}

func VerifyServiceAccountExists(ctx context.Context, kc *clients.KubeClient, sa, ns string) {
	log.Printf("Verify SA %q is created in namespace %q", sa, ns)

	if err := wait.PollUntilContextTimeout(ctx, config.APIRetry, config.APITimeout, true, func(context.Context) (done bool, err error) {
		_, err = kc.Kube.CoreV1().ServiceAccounts(ns).Get(ctx, sa, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, nil
		}
		return true, err
	}); err != nil {
		testsuit.T.Errorf("failed to get SA %s in namespace %s for tests: %v", sa, ns, err)
	}
}

func VerifyNamespaceExists(ctx context.Context, kc *clients.KubeClient, ns string) {
	log.Printf("Verify namespace %q exists", ns)
	if err := wait.PollUntilContextTimeout(ctx, config.APIRetry, config.APITimeout, true, func(context.Context) (done bool, err error) {
		_, err = kc.Kube.CoreV1().Namespaces().Get(ctx, ns, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, nil
		}
		return true, err
	}); err != nil {
		testsuit.T.Errorf("failed to get namespace %s for tests: %v", ns, err)
	}
}

func CreateCronJob(c *clients.Clients, args []string, schedule, namespace string) {
	cronjob := &batchv1.CronJob{
		TypeMeta: metav1.TypeMeta{APIVersion: batchv1.SchemeGroupVersion.String(), Kind: "CronJob"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "hello",
		},
		Spec: batchv1.CronJobSpec{
			Schedule: schedule,
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "hello",
				},
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "hello",
									Image: "image-registry.openshift-image-registry.svc:5000/openshift/golang",
									Args:  args,
								},
							},
							RestartPolicy: corev1.RestartPolicy("Never"),
						},
					},
				},
			},
		},
	}
	cj, err := c.KubeClient.Kube.BatchV1().CronJobs(namespace).Create(c.Ctx, cronjob, metav1.CreateOptions{})
	if err != nil {
		testsuit.T.Errorf("failed to create cron job %s \n %v", cj.Name, err)
	}
	log.Printf("Cronjob: %s created in namespace: %s", cj.Name, namespace)
	store.PutScenarioData("cronjob", cj.Name)
}

func DeleteCronJob(c *clients.Clients, name, ns string) error {
	propagationPolicy := metav1.DeletePropagationBackground // Also delete jobs and pods related to cronjob
	return c.KubeClient.Kube.BatchV1().CronJobs(ns).Delete(c.Ctx, name, metav1.DeleteOptions{PropagationPolicy: &propagationPolicy})
}

func Get(ctx context.Context, gr schema.GroupVersionResource, clients *clients.Clients, objname, ns string, op metav1.GetOptions) (*unstructured.Unstructured, error) {
	gvr, err := GetGroupVersionResource(gr, clients.Tekton.Discovery())
	if err != nil {
		return nil, err
	}

	obj, err := clients.Dynamic.Resource(*gvr).Namespace(ns).Get(ctx, objname, op)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// Watch func helps you to watch on dynamic resources
func Watch(ctx context.Context, gr schema.GroupVersionResource, clients *clients.Clients, ns string, op metav1.ListOptions) (watch.Interface, error) {
	gvr, err := GetGroupVersionResource(gr, clients.Tekton.Discovery())
	if err != nil {
		return nil, err
	}
	watch, err := clients.Dynamic.Resource(*gvr).Namespace(ns).Watch(ctx, op)
	if err != nil {
		return nil, err
	}
	return watch, nil
}

func GetGroupVersionResource(gr schema.GroupVersionResource, discovery discovery.DiscoveryInterface) (*schema.GroupVersionResource, error) {
	apiGroupRes, err := restmapper.GetAPIGroupResources(discovery)
	if err != nil {
		return nil, err
	}
	rm := restmapper.NewDiscoveryRESTMapper(apiGroupRes)
	gvr, err := rm.ResourceFor(gr)
	if err != nil {
		return nil, err
	}
	return &gvr, nil
}

func AssertIfDefaultCronjobExists(c *clients.Clients, namespace string) {
	cronJobs, err := c.KubeClient.Kube.BatchV1().CronJobs(namespace).List(c.Ctx, metav1.ListOptions{})
	if err != nil {
		testsuit.T.Errorf("failed to get default cron job in namespace %s \n %v", namespace, err)
	}
	if len(cronJobs.Items) == 0 {
		testsuit.T.Errorf("No cronjobs present in the namespace %v", namespace)
	}
	present := false
	for _, cj := range cronJobs.Items {
		if cj.Spec.Schedule == config.PrunerSchedule {
			if strings.Contains(cj.Name, config.PrunerNamePrefix) {
				present = true
				log.Printf("Cronjob with schedule %v and with name prefix %v is present", config.PrunerSchedule, config.PrunerNamePrefix)
				break
			}
		}
	}
	if !present {
		testsuit.T.Errorf("No cronjobs with schedule %v and with prefix %v is not present", config.PrunerSchedule, config.PrunerNamePrefix)
	}
}

func GetCronjobNameWithSchedule(c *clients.Clients, namespace, schedule string) string {
	name := ""
	cronJobs, err := c.KubeClient.Kube.BatchV1().CronJobs(namespace).List(c.Ctx, metav1.ListOptions{})
	if err != nil {
		testsuit.T.Errorf("failed to get cron job from namespace %s \n %v", namespace, err)
	}
	if len(cronJobs.Items) == 0 {
		testsuit.T.Errorf("No cronjobs present in the namespace %v", namespace)
	}
	for _, cj := range cronJobs.Items {
		if cj.Spec.Schedule == schedule {
			if strings.Contains(cj.Name, "tekton-resource-pruner-") {
				name = cj.Name
			}
		}
	}
	return name
}

func AssertPrunerCronjobWithContainer(c *clients.Clients, namespace, num string) {
	log.Printf("Verifying if the cronjob with prefix tekton-resource-pruner in namespace %v contains %v number of containers", namespace, num)
	cronJobs, err := c.KubeClient.Kube.BatchV1().CronJobs(namespace).List(c.Ctx, metav1.ListOptions{})
	if err != nil {
		testsuit.T.Errorf("Error while getting cronjobs %v", err)
	}
	jobFound := false
	for _, cr := range cronJobs.Items {
		if strings.Contains(cr.Name, "tekton-resource-pruner") {
			jobFound = true
			containers := cr.Spec.JobTemplate.Spec.Template.Spec.Containers
			numInt, _ := strconv.Atoi(num)
			if len(containers) != numInt {
				testsuit.T.Errorf("Expected: %v containers in cronjob spec, Actual: %v containers in cronjob spec", numInt, len(containers))
			}
			log.Printf("%v containers found in the cronjob spec", numInt)
			break
		}
	}
	if !jobFound {
		testsuit.T.Errorf("Cronjob with prefix tekton-resource-pruner not found in %v namespace", namespace)
	}
}

func AssertCronjobPresent(c *clients.Clients, cronJobName, namespace string) {
	err := wait.PollUntilContextTimeout(c.Ctx, config.APIRetry, config.ResourceTimeout, false, func(context.Context) (bool, error) {
		log.Printf("Verifying if cronjob with prefix %v is present in namespace %v", cronJobName, namespace)
		cjs, err := c.KubeClient.Kube.BatchV1().CronJobs(namespace).List(c.Ctx, metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, cj := range cjs.Items {
			if strings.Contains(cj.Name, cronJobName) {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Expected: cronjob with prefix %v present in namespace %v, Actual: cronjob with prefix %v not present in namespace %v", cronJobName, namespace, cronJobName, namespace))
	}
	log.Printf("Cronjob with prefix %v is present in namespace %v", cronJobName, namespace)
}

func AssertCronjobNotPresent(c *clients.Clients, cronJobName, namespace string) {
	err := wait.PollUntilContextTimeout(c.Ctx, config.APIRetry, config.ResourceTimeout, false, func(context.Context) (bool, error) {
		log.Printf("Verifying if cronjob with prefix %v is present in namespace %v", cronJobName, namespace)
		cjs, err := c.KubeClient.Kube.BatchV1().CronJobs(namespace).List(c.Ctx, metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, cj := range cjs.Items {
			if strings.Contains(cj.Name, cronJobName) {
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Expected: cronjob with prefix %v present in namespace %v, Actual: cronjob with prefix %v not present in namespace %v", cronJobName, namespace, cronJobName, namespace))
	}
	log.Printf("Cronjob with prefix %v is present in namespace %v", cronJobName, namespace)
}

func ValidateTektonInstallersetStatus(c *clients.Clients) {
	tis, err := c.Operator.TektonInstallerSets().List(c.Ctx, metav1.ListOptions{})
	failedInstallersets := make([]string, 0)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Error getting tektoninstallersets: %v", err))
	}

	for _, is := range tis.Items {
		log.Printf("Verifying if the installerset %s is in ready state", is.Name)
		if !is.Status.IsReady() {
			failedInstallersets = append(failedInstallersets, is.Name)
		}
	}

	if len(failedInstallersets) > 0 {
		testsuit.T.Fail(fmt.Errorf("The installersets %s is/are not in ready status", strings.Join(failedInstallersets, ",")))
	}
	log.Print("All the installersets are in ready state")
}

func ValidateTektonInstallersetNames(c *clients.Clients) {
	tis, err := c.Operator.TektonInstallerSets().List(c.Ctx, metav1.ListOptions{})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Error getting tektoninstallersets: %v", err))
	}
	missingInstallersets := make([]string, 0)
	for _, isp := range config.TektonInstallersetNamePrefixes {
		if !openshift.IsCapabilityEnabled(c, "Console") &&
			(isp == "addon-custom-consolecli" || isp == "addon-custom-openshiftconsole") {
			log.Printf("OpenShift Console is not enabled, skipping validation of installer set %s", isp)
			continue
		}

		if config.Flags.IsDisconnected && isp == "addon-custom-communityclustertask" {
			log.Printf("Testing on a disconnected cluster, skipping validation of installer set %s", isp)
			continue
		}

		log.Printf("Verifying if the installerset with prefix %s is present\n", isp)
		found := false
		for _, is := range tis.Items {
			if strings.HasPrefix(is.Name, isp) {
				found = true
				log.Printf("Installerset with prefix %s is present\n", isp)
				break
			}
		}

		if !found {
			missingInstallersets = append(missingInstallersets, isp)
		}
	}

	if len(missingInstallersets) > 0 {
		testsuit.T.Fail(fmt.Errorf("Installersets with prefix %s is not found", strings.Join(missingInstallersets, ",")))
	}
}

func GetWarningEvents(c *clients.Clients, namespace string) (string, error) {
	var eventString string
	var eventSlice []string
	events, err := c.KubeClient.Kube.CoreV1().Events(namespace).List(c.Ctx, metav1.ListOptions{FieldSelector: "type=Warning"})
	if err != nil {
		return eventString, err
	}
	for _, item := range events.Items {
		eventSlice = append(eventSlice, item.Message)
	}
	return strings.Join(eventSlice, "\n"), nil
}
