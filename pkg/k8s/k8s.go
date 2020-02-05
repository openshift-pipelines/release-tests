package k8s

import (
	"fmt"
	"log"

	"time"

	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	secv1 "github.com/openshift/api/security/v1"
	secclient "github.com/openshift/client-go/security/clientset/versioned/typed/security/v1"
	"github.com/tektoncd/pipeline/pkg/names"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

// NewClientSet is a setup function which helps you
// 	1. to creates clientSet instance to `client`
// 	2. creates Random namespace
func NewClientSet() (*clients.Clients, string, func()) {
	// TODO: fix this; method is in k8s but returns client.Clients
	ns := names.SimpleNameGenerator.RestrictLengthWithRandomSuffix("releasetest")
	cs := clients.NewClients(config.Flags.Kubeconfig, config.Flags.Cluster, ns)
	CreateNamespace(cs.KubeClient, ns)

	return cs, ns, func() {
		DeleteNamespace(cs.KubeClient, ns)
	}
}

// WaitForDeploymentDeletion checks to see if a given deployment is deleted
// the function returns an error if the given deployment is not deleted within the timeout
func WaitForDeploymentDeletion(cs *clients.Clients, namespace, name string) error {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		kc := cs.KubeClient.Kube
		_, err := kc.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{IncludeUninitialized: true})
		if err != nil {
			if apierrors.IsGone(err) || apierrors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		log.Printf("Waiting for deletion of %s deployment\n", name)
		return false, nil
	})
	assert.NoError(err, fmt.Sprintf("%s Deployment deletion failed\n", name))
	return err
}

func WaitForServiceAccount(cs *clients.Clients, ns, targetSA string) *corev1.ServiceAccount {

	ret := &corev1.ServiceAccount{}

	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		saList, err := cs.KubeClient.Kube.CoreV1().ServiceAccounts(ns).List(metav1.ListOptions{})
		for _, sa := range saList.Items {
			if sa.Name == targetSA {
				ret = &sa
				return true, nil
			}
		}
		return false, err
	})
	assert.NoError(err, fmt.Sprintf("ServiceAccount: %s, not found in namespace %s\n", targetSA, ns))
	return ret
}

func ValidateSCCAdded(cs *clients.Clients, ns, sa string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		privileged, err := GetPrivilegedSCC(cs)
		if err != nil {
			log.Printf("failed to get privileged scc: %s \n", err)
			return false, err
		}
		log.Printf("... looking at %v", privileged.Users)

		ctrlSA := fmt.Sprintf("system:serviceaccount:%s:%s", ns, sa)
		return inList(privileged.Users, ctrlSA), nil
	})
	assert.NoError(err, fmt.Sprintf("failed to Add privilaged scc: %s\n", sa))
}

func ValidateSCCRemoved(cs *clients.Clients, ns, sa string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		privileged, err := GetPrivilegedSCC(cs)
		if err != nil {
			log.Printf("failed to get privileged scc: %s \n", err)
			return false, err
		}

		ctrlSA := fmt.Sprintf("system:serviceaccount:%s:%s", ns, sa)
		return !inList(privileged.Users, ctrlSA), nil
	})
	assert.NoError(err, fmt.Sprintf("failed to Remove privilaged scc: %s\n", sa))
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
		err := WaitForDeployment(kc, ns,
			d,
			1,
			config.APIRetry,
			config.APITimeout,
		)
		assert.NoError(err, fmt.Sprintf("Deployments: %+v, failed to create\n", deployments))
	}

}

func GetPrivilegedSCC(cs *clients.Clients) (*secv1.SecurityContextConstraints, error) {
	sec, err := secclient.NewForConfig(cs.KubeConfig)
	if err != nil {
		return nil, err
	}
	return sec.SecurityContextConstraints().Get("privileged", metav1.GetOptions{})
}

func ValidateDeploymentDeletion(cs *clients.Clients, ns string, deployments ...string) {

	for _, d := range deployments {
		err := WaitForDeploymentDeletion(cs, ns, d)
		assert.NoError(err, fmt.Sprintf("Deployments: %+v, failed to delete\n", deployments))
	}
}

func WaitForDeployment(kc kubernetes.Interface, namespace, name string, replicas int, retryInterval, timeout time.Duration) error {

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		deployment, err := kc.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{IncludeUninitialized: true})
		if err != nil {
			if apierrors.IsNotFound(err) {
				log.Printf("Waiting for availability of %s deployment\n", name)
				return false, nil
			}
			return false, err
		}

		if int(deployment.Status.AvailableReplicas) == replicas {
			return true, nil
		}
		log.Printf("Waiting for full availability of %s deployment (%d/%d)\n", name, deployment.Status.AvailableReplicas, replicas)
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func CreateNamespace(kc *clients.KubeClient, ns string) {
	log.Printf("Create namespace %s ", ns)
	_, err := kc.Kube.CoreV1().Namespaces().Create(
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: ns},
		})
	assert.NoError(err, fmt.Sprintf("Failed to Created namespace: %s \n", ns))
}

func DeleteNamespace(kc *clients.KubeClient, namespace string) {
	log.Printf("Deleting namespace %s", namespace)
	if err := kc.Kube.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{}); err != nil {
		log.Printf("Failed to delete namespace %s: %s", namespace, err)
	}
}

func VerifyNoServiceAccount(kc *clients.KubeClient, sa, ns string) {
	log.Printf("Verify SA %q is absent in namespace %q", sa, ns)

	if err := wait.PollImmediate(config.APIRetry, config.APITimeout, func() (bool, error) {
		_, err := kc.Kube.CoreV1().ServiceAccounts(ns).Get(sa, metav1.GetOptions{})
		if err == nil || errors.IsNotFound(err) {
			return false, fmt.Errorf("sa %q exists in namespace %q", sa, ns)
		}
		return true, nil
	}); err != nil {
		log.Printf("Fail: SA %q exists in namespace %q, err: %s", sa, ns, err)
	}
}

func VerifyServiceAccountExists(kc *clients.KubeClient, namespace string) {
	// TODO: shouldn't this recieve an arg?
	defaultSA := "pipeline"
	log.Printf("Verify SA %q is created in namespace %q", defaultSA, namespace)

	if err := wait.PollImmediate(config.APIRetry, config.APITimeout, func() (bool, error) {
		_, err := kc.Kube.CoreV1().ServiceAccounts(namespace).Get(defaultSA, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, nil
		}
		return true, err
	}); err != nil {
		log.Printf("Failed to get SA %q in namespace %q for tests: %s", defaultSA, namespace, err)
	}
}
