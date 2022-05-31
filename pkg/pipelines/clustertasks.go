package pipelines

import (
	"fmt"
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetClusterTask(c *clients.Clients, clusterTaskName, status string) {
	log.Printf("Verifying if the clustertask %v is %v", clusterTaskName, status)
	_, err := c.ClustertaskClient.Get(c.Ctx, clusterTaskName, v1.GetOptions{})
	if status == "present" {
		if err == nil {
			log.Printf("Clustertask %v- Expected: Present, Actual: Present", clusterTaskName)
		} else {
			assert.FailOnError(fmt.Errorf("Clustertask %v- Expected: Present, Actual: Not Present, Error: %v", clusterTaskName, err))
		}

	} else {
		if err == nil {
			assert.FailOnError(fmt.Errorf("Clustertask %v- Expected: Not Present, Actual: Present", clusterTaskName))
		} else {
			log.Printf("Clustertask %v- Expected: Not Present, Actual: Not Present", clusterTaskName)
		}
	}
}
