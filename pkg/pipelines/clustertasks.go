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

func AssertClustertaskPresent(c *clients.Clients, clusterTaskName string) {
	err := wait.PollUntilContextTimeout(c.Ctx, config.APIRetry, config.ResourceTimeout, false, func(context.Context) (bool, error) {
		log.Printf("Verifying if the clustertask %v is present", clusterTaskName)
		_, err := c.ClustertaskClient.Get(c.Ctx, clusterTaskName, v1.GetOptions{})
		if err == nil {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Clustertasks %v Expected: Present, Actual: Not Present, Error: %v", clusterTaskName, err))
	} else {
		log.Printf("Clustertask %v is present", clusterTaskName)
	}
}

func AssertClustertaskNotPresent(c *clients.Clients, clusterTaskName string) {
	err := wait.PollUntilContextTimeout(c.Ctx, config.APIRetry, config.ResourceTimeout, false, func(context.Context) (bool, error) {
		log.Printf("Verifying if the clustertask %v is not present", clusterTaskName)
		_, err := c.ClustertaskClient.Get(c.Ctx, clusterTaskName, v1.GetOptions{})
		if err == nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Clustertasks %v Expected: Not Present, Actual: Present, Error: %v", clusterTaskName, err))
	} else {
		log.Printf("Clustertask %v is not present", clusterTaskName)
	}
}
