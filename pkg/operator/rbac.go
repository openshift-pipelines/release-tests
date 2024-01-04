package operator

import (
	"context"
	"fmt"
	"log"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	scc "github.com/openshift/client-go/security/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func AssertServiceAccountPresent(clients *clients.Clients, ns, targetSA string) {

	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		log.Printf("Verifying that service account %s exists\n", targetSA)
		saList, err := clients.KubeClient.Kube.CoreV1().ServiceAccounts(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, item := range saList.Items {
			if item.Name == targetSA {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Expected: Service account %v present in the namespace %v, Actual: Service account %v not present in the namespace %v, Error: %v", targetSA, ns, targetSA, ns, err))
	}
}
func AssertRoleBindingPresent(clients *clients.Clients, ns, roleBindingName string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		log.Printf("Verifying that role binding %s exists\n", roleBindingName)
		rbList, err := clients.KubeClient.Kube.RbacV1().RoleBindings(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, item := range rbList.Items {
			if item.Name == roleBindingName {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Expected: Rolebinding %v present in the namespace %v, Actual: Rolebinding %v not present in the namespace %v, Error: %v", roleBindingName, ns, roleBindingName, ns, err))
	}
}

func AssertConfigMapPresent(clients *clients.Clients, ns, configMapName string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		log.Printf("Verifying that config map %s exists\n", configMapName)
		rbList, err := clients.KubeClient.Kube.CoreV1().ConfigMaps(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, item := range rbList.Items {
			if item.Name == configMapName {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Expected: Configmap %v present in the namespace %v, Actual: Configmap %v not present in the namespace %v, Error: %v", configMapName, ns, configMapName, ns, err))
	}
}

func AssertClusterRolePresent(clients *clients.Clients, clusterRoleName string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		log.Printf("Verifying that cluster role %s exists\n", clusterRoleName)
		rbList, err := clients.KubeClient.Kube.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, item := range rbList.Items {
			if item.Name == clusterRoleName {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Expected: Clusterrole %v present, Actual: Clusterrole %v not present, Error: %v", clusterRoleName, clusterRoleName, err))
	}
}

func AssertServiceAccountNotPresent(clients *clients.Clients, ns, targetSA string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		log.Printf("Verifying that service account %s doesn't exist\n", targetSA)
		saList, err := clients.KubeClient.Kube.CoreV1().ServiceAccounts(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, item := range saList.Items {
			if item.Name == targetSA {
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Expected: Service account %v not present in the namespace %v, Actual: Service account %v is present in the namespace %v, Error: %v", targetSA, ns, targetSA, ns, err))
	}
}

func AssertRoleBindingNotPresent(clients *clients.Clients, ns, roleBindingName string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		log.Printf("Verifying that role binding %s doesn't exist\n", roleBindingName)
		rbList, err := clients.KubeClient.Kube.RbacV1().RoleBindings(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, item := range rbList.Items {
			if item.Name == roleBindingName {
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Expected: Rolebinding %v not present in the namespace %v, Actual: Rolebinding %v present in the namespace %v, Error: %v", roleBindingName, ns, roleBindingName, ns, err))
	}
}

func AssertConfigMapNotPresent(clients *clients.Clients, ns, configMapName string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		log.Printf("Verifying that config map %s doesn't exist\n", configMapName)
		cmList, err := clients.KubeClient.Kube.CoreV1().ConfigMaps(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, item := range cmList.Items {
			if item.Name == configMapName {
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Expected: Configmap %v not present in the namespace %v, Expected: Configmap %v present in the namespace %v, Error: %v", configMapName, ns, configMapName, ns, err))
	}
}

func AssertClusterRoleNotPresent(clients *clients.Clients, clusterRoleName string) {
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		log.Printf("Verifying that cluster role %s doesn't exist\n", clusterRoleName)
		rbList, err := clients.KubeClient.Kube.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, item := range rbList.Items {
			if item.Name == clusterRoleName {
				return false, nil
			}
		}
		return true, err
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Expected, Clusterrole %v not present, Actual: Clusterrole %v present, Error: %v", clusterRoleName, clusterRoleName, err))
	}
}

func AssertSCCPresent(clients *clients.Clients, sccName string) {
	s := scc.NewForConfigOrDie(clients.KubeConfig)
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		log.Printf("Verifying that security context constraint %s exists\n", sccName)
		sccList, err := s.SecurityV1().SecurityContextConstraints().List(clients.Ctx, metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, item := range sccList.Items {
			if item.Name == sccName {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Expected: security context constraint %q present, Actual: security context constraint %q not present , Error: %v", sccName, sccName, err))
	}
}

func AssertSCCNotPresent(clients *clients.Clients, sccName string) {
	s := scc.NewForConfigOrDie(clients.KubeConfig)
	err := wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		log.Printf("Verifying that security context constraint %s doesn't exist\n", sccName)
		sccList, err := s.SecurityV1().SecurityContextConstraints().List(clients.Ctx, metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, item := range sccList.Items {
			if item.Name == sccName {
				return false, nil
			}
		}
		return true, err
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Expected: secuirty context constraint %q not present, Actual: security context constraint %q present, Error: %v", sccName, sccName, err))
	}
}
