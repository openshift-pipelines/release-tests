package helper

import (
	"context"
	"fmt"
	"log"

	"time"

	. "github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/client"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	secv1 "github.com/openshift/api/security/v1"
	secclient "github.com/openshift/client-go/security/clientset/versioned/typed/security/v1"
	op "github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/names"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

// AssertNoError confirms the error returned is null
func AssertNoError(err error, description string) {
	//Expect(err).ShouldNot(HaveOccurred(), description)

	if err != nil {
		T.Errorf("%s, \n err:%s", description, err)
	}
}

// NewClientSet is a setup function which helps you
// 	1. to creates clientSet instance to `client`
// 	2. creates Random namespace
func NewClientSet() (*client.Clients, string, func()) {
	ns := names.SimpleNameGenerator.RestrictLengthWithRandomSuffix("releasetest")
	cs := client.NewClients(config.Flags.Kubeconfig, config.Flags.Cluster, ns)
	CreateNamespace(cs.KubeClient, ns)

	return cs, ns, func() {
		DeleteNamespace(cs.KubeClient, ns)
	}
}

// WaitForClusterCR waits for cluster CR to be created
// the function returns an error if Cluster CR is not created within timeout
func WaitForClusterCR(cs *client.Clients, name string) *op.Config {

	objKey := types.NamespacedName{Name: name}
	cr := &op.Config{}

	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		err := cs.Client.Get(context.TODO(), objKey, cr)
		if err != nil {
			if apierrors.IsNotFound(err) {
				log.Printf("Waiting for availability of %s cr\n", name)
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	AssertNoError(err, fmt.Sprintf("CR: %s is not avaialble\n", name))
	return cr
}

// WaitForDeploymentDeletion checks to see if a given deployment is deleted
// the function returns an error if the given deployment is not deleted within the timeout
func WaitForDeploymentDeletion(cs *client.Clients, namespace, name string) error {
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
	AssertNoError(err, fmt.Sprintf("%s Deployment deletion failed\n", name))
	return err
}

func WaitForServiceAccount(cs *client.Clients, ns, targetSA string) *corev1.ServiceAccount {

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
	AssertNoError(err, fmt.Sprintf("ServiceAccount: %s, not found in namespace %s\n", targetSA, ns))
	return ret
}

func DeleteClusterCR(cs *client.Clients, name string) {
	var err error
	// ensure object exists before deletion
	objKey := types.NamespacedName{Name: name}
	cr := &op.Config{}
	err = cs.Client.Get(context.TODO(), objKey, cr)

	AssertNoError(err, fmt.Sprintf("Failed to find cluster CR: %s : %s\n", name, err))

	err = wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		err := cs.Client.Delete(context.TODO(), cr)
		if err != nil {
			log.Printf("Deletion of CR %s failed %s \n", name, err)
			return false, err
		}

		return true, nil
	})

	AssertNoError(err, fmt.Sprintf("%s cluster CR deletion failed\n", name))
}

func ValidateSCCAdded(cs *client.Clients, ns, sa string) {
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
	AssertNoError(err, fmt.Sprintf("failed to Add privilaged scc: %s\n", sa))
}

func ValidateSCCRemoved(cs *client.Clients, ns, sa string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		privileged, err := GetPrivilegedSCC(cs)
		if err != nil {
			log.Printf("failed to get privileged scc: %s \n", err)
			return false, err
		}

		ctrlSA := fmt.Sprintf("system:serviceaccount:%s:%s", ns, sa)
		return !inList(privileged.Users, ctrlSA), nil
	})
	AssertNoError(err, fmt.Sprintf("failed to Remove privilaged scc: %s\n", sa))
}

func inList(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func ValidateDeployments(cs *client.Clients, ns string, deployments ...string) {

	kc := cs.KubeClient.Kube
	for _, d := range deployments {
		err := WaitForDeployment(kc, ns,
			d,
			1,
			config.APIRetry,
			config.APITimeout,
		)
		AssertNoError(err, fmt.Sprintf("Deployments: %+v, failed to create\n", deployments))
	}

}

func GetPrivilegedSCC(cs *client.Clients) (*secv1.SecurityContextConstraints, error) {
	sec, err := secclient.NewForConfig(cs.KubeConfig)
	if err != nil {
		return nil, err
	}
	return sec.SecurityContextConstraints().Get("privileged", metav1.GetOptions{})
}

func ValidateDeploymentDeletion(cs *client.Clients, ns string, deployments ...string) {

	for _, d := range deployments {
		err := WaitForDeploymentDeletion(cs, ns, d)
		AssertNoError(err, fmt.Sprintf("Deployments: %+v, failed to delete\n", deployments))
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

func CreateNamespace(kc *client.KubeClient, ns string) {
	log.Printf("Create namespace %s ", ns)
	_, err := kc.Kube.CoreV1().Namespaces().Create(
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: ns},
		})
	AssertNoError(err, fmt.Sprintf("Failed to Created namespace: %s \n", ns))
}

func DeleteNamespace(kc *client.KubeClient, namespace string) {
	log.Printf("Deleting namespace %s", namespace)
	if err := kc.Kube.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{}); err != nil {
		log.Printf("Failed to delete namespace %s: %s", namespace, err)
	}
}

func VerifyServiceAccountExists(kc *client.KubeClient, namespace string) {
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
