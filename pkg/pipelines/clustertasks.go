package pipelines

import (
	"fmt"
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AssertClustertaskPresent(c *clients.Clients, clusterTaskName string) {
	log.Printf("Verifying if the clustertask %v is present", clusterTaskName)
	_, err := c.ClustertaskClient.Get(c.Ctx, clusterTaskName, v1.GetOptions{})
	if err != nil {
		assert.FailOnError(fmt.Errorf("Clustertasks %v Expected: Present, Actual: Not Present", clusterTaskName))
	}
	log.Printf("Clustertask %v is present", clusterTaskName)
}

func AssertClustertaskNotPresent(c *clients.Clients, clusterTaskName string) {
	log.Printf("Verifying if the clustertask %v is present", clusterTaskName)
	_, err := c.ClustertaskClient.Get(c.Ctx, clusterTaskName, v1.GetOptions{})
	if err == nil {
		assert.FailOnError(fmt.Errorf("Clustertasks %v Expected: Not Present, Actual: Present", clusterTaskName))
	}
	log.Printf("Clustertask %v is not present", clusterTaskName)
}
