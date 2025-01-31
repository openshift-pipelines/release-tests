package statefulset

import (
	"context"
	"fmt"
	"log"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func ValidateStatefulSetDeployment(cs *clients.Clients, deploymentName string) {
	labelSelector := "app.kubernetes.io/part-of=tekton-pipelines"
	listOptions := metav1.ListOptions{LabelSelector: labelSelector}

	log.Printf("Starting validation for StatefulSet deployment: %s in namespace: %s", deploymentName, config.TargetNamespace)

	waitErr := wait.PollUntilContextTimeout(context.TODO(), config.APIRetry, config.APITimeout, true, func(ctx context.Context) (bool, error) {
		stsList, err := cs.KubeClient.Kube.AppsV1().StatefulSets(config.TargetNamespace).List(context.TODO(), listOptions)
		if err != nil {
			log.Printf("Error listing StatefulSets: %v", err)
			return false, fmt.Errorf("failed to list StatefulSets: %v", err)
		}

		log.Printf("Found %d StatefulSets in namespace %s", len(stsList.Items), config.TargetNamespace)

		for _, sts := range stsList.Items {
			if sts.Name == deploymentName {
				log.Printf("Found StatefulSet: %s", sts.Name)
				isAvailable, err := IsStatefulSetAvailable(&sts)
				if err != nil {
					log.Printf("Error checking availability of StatefulSet %s: %v", sts.Name, err)
					return false, err
				}
				if isAvailable {
					log.Printf("StatefulSet %s is available and ready", sts.Name)
					return true, nil
				} else {
					log.Printf("StatefulSet %s is not ready yet", sts.Name)
					return false, nil
				}
			}
		}
		log.Printf("StatefulSet %s not found yet. Continuing to wait...", deploymentName)
		return false, nil
	})

	if waitErr != nil {
		testsuit.T.Fatalf("StatefulSet %s was not found or not available within 5 minutes in the namespace %q: %v",
			deploymentName, config.TargetNamespace, waitErr)
	}
}

func IsStatefulSetAvailable(sts *appsv1.StatefulSet) (bool, error) {
	if sts.Status.ReadyReplicas == *sts.Spec.Replicas {
		return true, nil
	}
	return false, nil
}
