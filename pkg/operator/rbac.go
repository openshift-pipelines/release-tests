package operator

import (
	"context"
	"fmt"

	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func AssertServiceAccount(clients *clients.Clients, ns, targetSA string) {

	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		saList, err := clients.KubeClient.Kube.CoreV1().ServiceAccounts(ns).List(context.TODO(), metav1.ListOptions{})
		for _, item := range saList.Items {
			if item.Name == targetSA {
				return true, nil
			}
		}
		return false, err
	})
	if err != nil {
		assert.FailOnError(fmt.Errorf("could not find serviceaccount %s/%s: %q", ns, targetSA, err))
	}
}
func AssertRoleBinding(clients *clients.Clients, ns, roleBindingName string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		rbList, err := clients.KubeClient.Kube.RbacV1().RoleBindings(ns).List(context.TODO(), metav1.ListOptions{})
		for _, item := range rbList.Items {
			if item.Name == roleBindingName {
				return true, nil
			}
		}
		return false, err
	})
	if err != nil {
		assert.FailOnError(fmt.Errorf("could not find Rolebinding %s/%s: %q", ns, roleBindingName, err))
	}
}

func AssertConfigMap(clients *clients.Clients, ns, configMapName string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		rbList, err := clients.KubeClient.Kube.CoreV1().ConfigMaps(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, item := range rbList.Items {
			if item.Name == configMapName {
				return true, nil
			}
		}
		return false, err
	})
	if err != nil {
		assert.FailOnError(fmt.Errorf("could not find ConfigMap %s/%s: %q", ns, configMapName, err))
	}
}

func AssertClusterRole(clients *clients.Clients, clusterRoleName string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		rbList, err := clients.KubeClient.Kube.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, item := range rbList.Items {
			if item.Name == clusterRoleName {
				return true, nil
			}
		}
		return false, err
	})
	if err != nil {
		assert.FailOnError(fmt.Errorf("could not find ClusterRole %s: %q", clusterRoleName, err))
	}
}
