package wait

import (
	"context"
	"fmt"

	"github.com/openshift-pipelines/release-tests/pkg/client"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"go.opencensus.io/trace"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"knative.dev/pkg/apis"
)

// TaskRunStateFn is a condition function on TaskRun used polling functions
type TaskRunStateFn func(r *v1alpha1.TaskRun) (bool, error)

// PipelineRunStateFn is a condition function on TaskRun used polling functions
type PipelineRunStateFn func(pr *v1alpha1.PipelineRun) (bool, error)

// WaitForTaskRunState polls the status of the TaskRun called name from client every
// interval until inState returns `true` indicating it is done, returns an
// error or timeout. desc will be used to name the metric that is emitted to
// track how long it took for name to get into the state checked by inState.
func WaitForTaskRunState(c *client.Clients, name string, inState TaskRunStateFn, desc string) error {
	metricName := fmt.Sprintf("WaitForTaskRunState/%s/%s", name, desc)
	_, span := trace.StartSpan(context.Background(), metricName)
	defer span.End()

	return wait.PollImmediate(config.Interval, config.Timeout, func() (bool, error) {
		r, err := c.TaskRunClient.Get(name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		return inState(r)
	})
}

// WaitForDeploymentState polls the status of the Deployment called name
// from client every interval until inState returns `true` indicating it is done,
// returns an  error or timeout. desc will be used to name the metric that is emitted to
// track how long it took for name to get into the state checked by inState.
func WaitForDeploymentState(c *client.Clients, name string, namespace string, inState func(d *appsv1.Deployment) (bool, error), desc string) error {
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
func WaitForPodState(c *client.Clients, name string, namespace string, inState func(r *corev1.Pod) (bool, error), desc string) error {
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
func WaitForPipelineRunState(c *client.Clients, name string, inState PipelineRunStateFn, desc string) error {
	metricName := fmt.Sprintf("WaitForPipelineRunState/%s/%s", name, desc)
	_, span := trace.StartSpan(context.Background(), metricName)
	defer span.End()

	return wait.PollImmediate(config.Interval, config.Timeout, func() (bool, error) {
		r, err := c.PipelineRunClient.Get(name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		return inState(r)
	})
}

// TaskRunSucceed provides a poll condition function that checks if the TaskRun
// has successfully completed.
func TaskRunSucceed(name string) TaskRunStateFn {
	return func(tr *v1alpha1.TaskRun) (bool, error) {
		c := tr.Status.GetCondition(apis.ConditionSucceeded)
		if c != nil {
			if c.Status == corev1.ConditionTrue {
				return true, nil
			} else if c.Status == corev1.ConditionFalse {
				return true, fmt.Errorf("task run %q failed", name)
			}
		}
		return false, nil
	}
}

// TaskRunFailed provides a poll condition function that checks if the TaskRun
// has failed.
func TaskRunFailed(name string) TaskRunStateFn {
	return func(tr *v1alpha1.TaskRun) (bool, error) {
		c := tr.Status.GetCondition(apis.ConditionSucceeded)
		if c != nil {
			if c.Status == corev1.ConditionTrue {
				return true, fmt.Errorf("task run %q succeeded", name)
			} else if c.Status == corev1.ConditionFalse {
				return true, nil
			}
		}
		return false, nil
	}
}

// PipelineRunSucceed provides a poll condition function that checks if the PipelineRun
// has successfully completed.
func PipelineRunSucceed(name string) PipelineRunStateFn {
	return func(pr *v1alpha1.PipelineRun) (bool, error) {
		c := pr.Status.GetCondition(apis.ConditionSucceeded)
		if c != nil {
			if c.Status == corev1.ConditionTrue {
				return true, nil
			} else if c.Status == corev1.ConditionFalse {
				return true, fmt.Errorf("pipeline run %q failed", name)
			}
		}
		return false, nil
	}
}

// PipelineRunFailed provides a poll condition function that checks if the PipelineRun
// has failed.
func PipelineRunFailed(name string) PipelineRunStateFn {
	return func(tr *v1alpha1.PipelineRun) (bool, error) {
		c := tr.Status.GetCondition(apis.ConditionSucceeded)
		if c != nil {
			if c.Status == corev1.ConditionTrue {
				return true, fmt.Errorf("task run %q succeeded", name)
			} else if c.Status == corev1.ConditionFalse {
				return true, nil
			}
		}
		return false, nil
	}
}
