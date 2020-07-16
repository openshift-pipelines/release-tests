package pipelines

import (
	"bytes"
	"fmt"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// collectMatchingEvents collects list of events under 5 seconds that match
// 1. matchKinds which is a map of Kind of Object with name of objects
// 2. reason which is the expected reason of event
func collectMatchingEvents(c *clients.Clients, namespace string, kinds map[string][]string, reason string) ([]*corev1.Event, error) {
	var events []*corev1.Event

	watchEvents, err := c.KubeClient.Kube.CoreV1().Events(namespace).Watch(metav1.ListOptions{})
	// close watchEvents channel
	defer watchEvents.Stop()
	if err != nil {
		return events, err
	}

	// create timer to not wait for events longer than 5 seconds
	timer := time.NewTimer(5 * time.Second)

	for {
		select {
		case wevent := <-watchEvents.ResultChan():
			event := wevent.Object.(*corev1.Event)
			if val, ok := kinds[event.InvolvedObject.Kind]; ok {
				for _, expectedName := range val {
					if event.InvolvedObject.Name == expectedName && event.Reason == reason {
						events = append(events, event)
					}
				}
			}
		case <-timer.C:
			return events, nil
		}
	}
}

// checkLabelPropagation checks that labels are correctly propagating from
// Pipelines, PipelineRuns, and Tasks to TaskRuns and Pods.
func checkLabelPropagation(c *clients.Clients, namespace string, pipelineRunName string, tr *v1beta1.TaskRun) {
	// Our controllers add 4 labels automatically. If custom labels are set on
	// the Pipeline, PipelineRun, or Task then the map will have to be resized.
	labels := make(map[string]string, 4)

	// Check label propagation to PipelineRuns.
	pr, err := c.PipelineRunClient.Get(pipelineRunName, metav1.GetOptions{})
	assert.NoError(err, fmt.Sprintf("Couldn't get expected PipelineRun for %s: %s", tr.Name, err))

	p, err := c.PipelineClient.Get(pr.Spec.PipelineRef.Name, metav1.GetOptions{})
	assert.NoError(err, fmt.Sprintf("Couldn't get expected Pipeline for %s: %s", pr.Name, err))

	for key, val := range p.ObjectMeta.Labels {
		labels[key] = val
	}
	// This label is added to every PipelineRun by the PipelineRun controller
	labels[pipeline.GroupName+pipeline.PipelineLabelKey] = p.Name
	AssertLabelsMatch(labels, pr.ObjectMeta.Labels)

	// Check label propagation to TaskRuns.
	for key, val := range pr.ObjectMeta.Labels {
		labels[key] = val
	}
	// This label is added to every TaskRun by the PipelineRun controller
	labels[pipeline.GroupName+pipeline.PipelineRunLabelKey] = pr.Name
	if tr.Spec.TaskRef != nil {
		task, err := c.TaskClient.Get(tr.Spec.TaskRef.Name, metav1.GetOptions{})
		assert.NoError(err, fmt.Sprintf("Couldn't get expected Task for %s: %s", tr.Name, err))
		for key, val := range task.ObjectMeta.Labels {
			labels[key] = val
		}
		// This label is added to TaskRuns that reference a Task by the TaskRun controller
		labels[pipeline.GroupName+pipeline.TaskLabelKey] = task.Name
	}
	AssertLabelsMatch(labels, tr.ObjectMeta.Labels)

	// PodName is "" iff a retry happened and pod is deleted
	// This label is added to every Pod by the TaskRun controller
	if tr.Status.PodName != "" {
		// Check label propagation to Pods.
		pod := GetPodForTaskRun(c, namespace, tr)
		// This label is added to every Pod by the TaskRun controller
		labels[pipeline.GroupName+pipeline.TaskRunLabelKey] = tr.Name
		AssertLabelsMatch(labels, pod.ObjectMeta.Labels)
	}
}

// checkAnnotationPropagation checks that annotations are correctly propagating from
// Pipelines, PipelineRuns, and Tasks to TaskRuns and Pods.
func checkAnnotationPropagation(c *clients.Clients, namespace string, pipelineRunName string, tr *v1beta1.TaskRun) {
	annotations := make(map[string]string)

	// Check annotation propagation to PipelineRuns.
	pr, err := c.PipelineRunClient.Get(pipelineRunName, metav1.GetOptions{})
	assert.NoError(err, fmt.Sprintf("Couldn't get expected PipelineRun for %s: %s", tr.Name, err))
	p, err := c.PipelineClient.Get(pr.Spec.PipelineRef.Name, metav1.GetOptions{})
	assert.NoError(err, fmt.Sprintf("Couldn't get expected Pipeline for %s: %s", pr.Name, err))
	for key, val := range p.ObjectMeta.Annotations {
		annotations[key] = val
	}
	AssertAnnotationsMatch(annotations, pr.ObjectMeta.Annotations)

	// Check annotation propagation to TaskRuns.
	for key, val := range pr.ObjectMeta.Annotations {
		annotations[key] = val
	}
	if tr.Spec.TaskRef != nil {
		task, err := c.TaskClient.Get(tr.Spec.TaskRef.Name, metav1.GetOptions{})
		assert.NoError(err, fmt.Sprintf("Couldn't get expected Task for %s: %s", tr.Name, err))
		for key, val := range task.ObjectMeta.Annotations {
			annotations[key] = val
		}
	}
	AssertAnnotationsMatch(annotations, tr.ObjectMeta.Annotations)

	// Check annotation propagation to Pods.
	pod := GetPodForTaskRun(c, namespace, tr)
	AssertAnnotationsMatch(annotations, pod.ObjectMeta.Annotations)
}

func GetPodForTaskRun(c *clients.Clients, namespace string, tr *v1beta1.TaskRun) *corev1.Pod {
	// The Pod name has a random suffix, so we filter by label to find the one we care about.
	pods, err := c.KubeClient.Kube.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: pipeline.GroupName + pipeline.TaskRunLabelKey + " = " + tr.Name,
	})
	assert.NoError(err, fmt.Sprintf("Couldn't get expected Pod for %s: %s", tr.Name, err))
	if numPods := len(pods.Items); numPods != 1 {
		testsuit.T.Errorf(fmt.Sprintf("Expected 1 Pod for %s, but got %d Pods", tr.Name, numPods))
	}
	return &pods.Items[0]
}

func AssertLabelsMatch(expectedLabels, actualLabels map[string]string) {
	for key, expectedVal := range expectedLabels {
		if actualVal := actualLabels[key]; actualVal != expectedVal {
			testsuit.T.Errorf(fmt.Sprintf("Expected labels containing %s=%s but labels were %v", key, expectedVal, actualLabels))
		}
	}
}

func AssertAnnotationsMatch(expectedAnnotations, actualAnnotations map[string]string) {
	for key, expectedVal := range expectedAnnotations {
		if actualVal := actualAnnotations[key]; actualVal != expectedVal {
			testsuit.T.Errorf(fmt.Sprintf("Expected annotations containing %s=%s but annotations were %v", key, expectedVal, actualAnnotations))
		}
	}
}

func Cast2pipelinerun(obj runtime.Object) (*v1beta1.PipelineRun, error) {
	var run *v1beta1.PipelineRun
	unstruct, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstruct, &run); err != nil {
		return nil, err
	}
	return run, nil
}

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s = %s\n", key, value)
	}
	return b.String()
}
