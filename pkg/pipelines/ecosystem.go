package pipelines

import (
	"context"
	"fmt"
	"log"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func AssertTaskPresent(c *clients.Clients, namespace string, taskName string) {
	err := wait.PollUntilContextTimeout(c.Ctx, config.APIRetry, config.ResourceTimeout, false, func(context.Context) (bool, error) {
		log.Printf("Verifying if the task %v is present", taskName)
		_, err := c.Tekton.TektonV1().Tasks(namespace).Get(c.Ctx, taskName, v1.GetOptions{})
		if err == nil {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Tasks %v Expected: Present, Actual: Not Present, Error: %v", taskName, err))
	} else {
		log.Printf("Task %v is present", taskName)
	}
}

func AssertTaskNotPresent(c *clients.Clients, namespace string, taskName string) {
	err := wait.PollUntilContextTimeout(c.Ctx, config.APIRetry, config.ResourceTimeout, false, func(context.Context) (bool, error) {
		log.Printf("Verifying if the task %v is not present", taskName)
		_, err := c.Tekton.TektonV1().Tasks(namespace).Get(c.Ctx, taskName, v1.GetOptions{})
		if err == nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Tasks %v Expected: Not Present, Actual: Present, Error: %v", taskName, err))
	} else {
		log.Printf("Task %v is not present", taskName)
	}
}

func AssertStepActionPresent(c *clients.Clients, namespace string, stepActionName string) {
	err := wait.PollUntilContextTimeout(c.Ctx, config.APIRetry, config.ResourceTimeout, false, func(context.Context) (bool, error) {
		log.Printf("Verifying if the stepAction %v is present", stepActionName)
		_, err := c.Tekton.TektonV1beta1().StepActions(namespace).Get(c.Ctx, stepActionName, v1.GetOptions{})
		if err == nil {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("StepAction %v Expected: Present, Actual: Not Present, Error: %v", stepActionName, err))
	} else {
		log.Printf("StepAction %v is present", stepActionName)
	}
}

func AssertStepActionNotPresent(c *clients.Clients, namespace string, stepActionName string) {
	err := wait.PollUntilContextTimeout(c.Ctx, config.APIRetry, config.ResourceTimeout, false, func(context.Context) (bool, error) {
		log.Printf("Verifying if the stepAction %v is not present", stepActionName)
		_, err := c.Tekton.TektonV1beta1().StepActions(namespace).Get(c.Ctx, stepActionName, v1.GetOptions{})
		if err == nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("StepAction %v Expected: Not Present, Actual: Present, Error: %v", stepActionName, err))
	} else {
		log.Printf("StepAction %v is not present", stepActionName)
	}
}
