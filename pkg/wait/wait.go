package wait

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"go.opencensus.io/trace"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"knative.dev/pkg/apis"
)

// ConditionAccessorFn is a condition function used polling functions
type ConditionAccessorFn func(ca apis.ConditionAccessor) (bool, error)

// WaitForTaskRunState polls the status of the TaskRun called name from client every
// interval until inState returns `true` indicating it is done, returns an
// error or timeout. desc will be used to name the metric that is emitted to
// track how long it took for name to get into the state checked by inState.
func WaitForTaskRunState(c *clients.Clients, name string, inState ConditionAccessorFn, desc string) error {
	metricName := fmt.Sprintf("WaitForTaskRunState/%s/%s", name, desc)
	_, span := trace.StartSpan(context.Background(), metricName)
	defer span.End()

	return wait.PollImmediate(config.Interval, config.Timeout, func() (bool, error) {
		r, err := c.TaskRunClient.Get(name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		return inState(&r.Status)
	})
}

// WaitForDeploymentState polls the status of the Deployment called name
// from client every interval until inState returns `true` indicating it is done,
// returns an  error or timeout. desc will be used to name the metric that is emitted to
// track how long it took for name to get into the state checked by inState.
func WaitForDeploymentState(c *clients.Clients, name string, namespace string, inState func(d *appsv1.Deployment) (bool, error), desc string) error {
	metricName := fmt.Sprintf("WaitForDeploymentState/%s/%s", name, desc)
	_, span := trace.StartSpan(context.Background(), metricName)
	defer span.End()

	return wait.PollImmediate(config.Interval, config.Timeout, func() (bool, error) {
		d, err := c.KubeClient.Kube.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		return inState(d)
	})
}

// WaitForPodState polls the status of the Pod called name from client every
// interval until inState returns `true` indicating it is done, returns an
// error or timeout. desc will be used to name the metric that is emitted to
// track how long it took for name to get into the state checked by inState.
func WaitForPodState(c *clients.Clients, name string, namespace string, inState func(r *corev1.Pod) (bool, error), desc string) error {
	metricName := fmt.Sprintf("WaitForPodState/%s/%s", name, desc)
	_, span := trace.StartSpan(context.Background(), metricName)
	defer span.End()

	return wait.PollImmediate(config.Interval, config.Timeout, func() (bool, error) {
		r, err := c.KubeClient.Kube.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		return inState(r)
	})
}

// WaitForPipelineRunState polls the status of the PipelineRun called name from client every
// interval until inState returns `true` indicating it is done, returns an
// error or timeout. desc will be used to name the metric that is emitted to
// track how long it took for name to get into the state checked by inState.
func WaitForPipelineRunState(c *clients.Clients, name string, inState ConditionAccessorFn, desc string) error {
	metricName := fmt.Sprintf("WaitForPipelineRunState/%s/%s", name, desc)
	_, span := trace.StartSpan(context.Background(), metricName)
	defer span.End()

	return wait.PollImmediate(config.Interval, config.Timeout, func() (bool, error) {
		r, err := c.PipelineRunClient.Get(name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		return inState(&r.Status)
	})
}

// WaitForServiceExternalIPState polls the status of the a k8s Service called name from client every
// interval until an external ip is assigned indicating it is done, returns an
// error or timeout. desc will be used to name the metric that is emitted to
// track how long it took for name to get into the state checked by inState.
func WaitForServiceExternalIPState(c *clients.Clients, namespace, name string, inState func(s *corev1.Service) (bool, error), desc string) error {
	metricName := fmt.Sprintf("WaitForServiceExternalIPState/%s/%s", name, desc)
	_, span := trace.StartSpan(context.Background(), metricName)
	defer span.End()

	return wait.PollImmediate(config.Interval, config.Timeout, func() (bool, error) {
		r, err := c.KubeClient.Kube.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		return inState(r)
	})
}

// Succeed provides a poll condition function that checks if the ConditionAccessor
// resource has successfully completed or not.
func Succeed(name string) ConditionAccessorFn {
	return func(ca apis.ConditionAccessor) (bool, error) {
		c := ca.GetCondition(apis.ConditionSucceeded)
		if c != nil {
			if c.Status == corev1.ConditionTrue {
				return true, nil
			} else if c.Status == corev1.ConditionFalse {
				return true, fmt.Errorf("%q failed", name)
			}
		}
		return false, nil
	}
}

// Failed provides a poll condition function that checks if the ConditionAccessor
// resource has failed or not.
func Failed(name string) ConditionAccessorFn {
	return func(ca apis.ConditionAccessor) (bool, error) {
		c := ca.GetCondition(apis.ConditionSucceeded)
		if c != nil {
			if c.Status == corev1.ConditionTrue {
				return true, fmt.Errorf("%q succeeded", name)
			} else if c.Status == corev1.ConditionFalse {
				return true, nil
			}
		}
		return false, nil
	}
}

// FailedWithReason provides a poll function that checks if the ConditionAccessor
// resource has failed with the TimeoudOut reason
func FailedWithReason(reason, name string) ConditionAccessorFn {
	return func(ca apis.ConditionAccessor) (bool, error) {
		c := ca.GetCondition(apis.ConditionSucceeded)
		if c != nil {
			if c.Status == corev1.ConditionFalse {
				if c.Reason == reason {
					return true, nil
				}
				return true, fmt.Errorf("%q completed with the wrong reason: %s", name, c.Reason)
			} else if c.Status == corev1.ConditionTrue {
				return true, fmt.Errorf("%q completed successfully, should have been failed with reason %q", name, reason)
			}
		}
		return false, nil
	}
}

// FailedWithMessage provides a poll function that checks if the ConditionAccessor
// resource has failed with the TimeoudOut reason
func FailedWithMessage(message, name string) ConditionAccessorFn {
	return func(ca apis.ConditionAccessor) (bool, error) {
		c := ca.GetCondition(apis.ConditionSucceeded)
		if c != nil {
			if c.Status == corev1.ConditionFalse {
				if strings.Contains(c.Message, message) {
					return true, nil
				}
				return true, fmt.Errorf("%q completed with the wrong message: %s", name, c.Message)
			} else if c.Status == corev1.ConditionTrue {
				return true, fmt.Errorf("%q completed successfully, should have been failed with message %q", name, message)
			}
		}
		return false, nil
	}
}

// Running provides a poll condition function that checks if the ConditionAccessor
// resource is currently running.
func Running(name string) ConditionAccessorFn {
	return func(ca apis.ConditionAccessor) (bool, error) {
		c := ca.GetCondition(apis.ConditionSucceeded)
		if c != nil {
			if c.Status == corev1.ConditionTrue || c.Status == corev1.ConditionFalse {
				return true, fmt.Errorf(`%q already finished`, name)
			} else if c.Status == corev1.ConditionUnknown && (c.Reason == "Running" || c.Reason == "Pending") {
				return true, nil
			}
		}
		return false, nil
	}
}

// TaskRunSucceed provides a poll condition function that checks if the TaskRun
// has successfully completed.
func TaskRunSucceed(name string) ConditionAccessorFn {
	return Succeed(name)
}

// TaskRunFailed provides a poll condition function that checks if the TaskRun
// has failed.
func TaskRunFailed(name string) ConditionAccessorFn {
	return Failed(name)
}

// PipelineRunSucceed provides a poll condition function that checks if the PipelineRun
// has successfully completed.
func PipelineRunSucceed(name string) ConditionAccessorFn {
	return Succeed(name)
}

// PipelineRunFailed provides a poll condition function that checks if the PipelineRun
// has failed.
func PipelineRunFailed(name string) ConditionAccessorFn {
	return Failed(name)
}

// ============================== Triggers Wait ==============================================

// WaitFor waits for the specified ConditionFunc every internal until the timeout.
func WaitFor(waitFunc wait.ConditionFunc) error {
	return wait.PollImmediate(config.Interval, config.Timeout, waitFunc)
}

// EventListenerReady returns a function that checks if all conditions on the
// specified EventListener are true and that the deployment available condition
// is within this set
func EventListenerReady(c *clients.Clients, namespace, name string) wait.ConditionFunc {
	return func() (bool, error) {
		el, err := c.TriggersClient.TriggersV1alpha1().EventListeners(namespace).Get(name, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			log.Printf("EventListener not found")
			return false, nil
		}
		log.Printf("EventListenerStatus: %+v", el.Status)
		// No conditions have been set yet
		if len(el.Status.Conditions) == 0 {
			return false, nil
		}
		if el.Status.GetCondition(apis.ConditionType(appsv1.DeploymentAvailable)) == nil {
			return false, nil
		}
		for _, cond := range el.Status.Conditions {
			if cond.Status != corev1.ConditionTrue {
				return false, nil
			}
		}
		return true, nil
	}
}

func WaitForPodsWithLabels(c *clients.Clients, namespace, labels string) wait.ConditionFunc {
	lastKnownPodNumber := -1
	return func() (bool, error) {
		listOpts := metav1.ListOptions{LabelSelector: labels}
		pods, err := c.KubeClient.Kube.CoreV1().Pods(namespace).List(listOpts)
		if err != nil {
			log.Printf("[apiclient] Error getting Pods with label selector %q [%v]\n", labels, err)
			return false, nil
		}

		if lastKnownPodNumber != len(pods.Items) {
			log.Printf("[apiclient] Found %d Pods for label selector %s\n", len(pods.Items), labels)
			lastKnownPodNumber = len(pods.Items)
		}

		if len(pods.Items) == 0 {
			return false, nil
		}

		for _, pod := range pods.Items {
			if pod.Status.Phase != corev1.PodRunning {
				return false, nil
			}
		}

		return true, nil
	}
}

// DeploymentNotExist returns a function that checks if the specified Deployment does not exist
func DeploymentNotExist(c *clients.Clients, namespace, name string) wait.ConditionFunc {
	return func() (bool, error) {
		_, err := c.KubeClient.Kube.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return true, nil
		}
		return false, nil
	}
}

// ServiceNotExist returns a function that checks if the specified Service does not exist
func ServiceNotExist(c *clients.Clients, namespace, name string) wait.ConditionFunc {
	return func() (bool, error) {
		_, err := c.KubeClient.Kube.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return true, nil
		}
		return false, nil
	}
}

// PipelineResourceExist returns a function that checks if the specified PipelineResource exists
func PipelineResourceExist(c *clients.Clients, name string) wait.ConditionFunc {
	return func() (bool, error) {
		_, err := c.PipelineResourceClient.Get(name, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, nil
		}
		return true, err
	}
}

// PipelineRunExist returns a function that checks if the specified PipelineRun exists
func PipelineRunExist(c *clients.Clients, name string) wait.ConditionFunc {
	return func() (bool, error) {
		_, err := c.PipelineRunClient.Get(name, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, nil
		}
		return true, err
	}
}
