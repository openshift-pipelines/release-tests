package olm

import (
	"fmt"
	"log"
	"strings"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/models"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Update TektonConfig CR to use param with name createRbacResource and value <value> to <action> auto creation of RBAC resources", func(value, action string) {
	patchData := fmt.Sprintf("{\"spec\":{\"params\":[{\"name\":\"createRbacResource\",\"value\":\"%s\"}]}}", value)
	log.Println(action, "auto creation of RBAC resources")
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "patch", "TektonConfig", "config", "--type=merge", "-p", patchData).Stdout())
})

var _ = gauge.Step("Verify RBAC resources disabled successfully", func() {
	operator.ValidateRBACAfterDisable(store.Clients(), store.GetCRNames())
})

var _ = gauge.Step("Verify RBAC resources are auto created successfully", func() {
	operator.ValidateRBAC(store.Clients(), store.GetCRNames())
})

var _ = gauge.Step("Verify the roles are present in <namespace> namespace: <table>", func(namespace string, rolesTable *models.Table) {
	gauge.GetScenarioStore()["rolesTable"] = rolesTable
	for _, row := range rolesTable.Rows {
		role := row.Cells[0]
		operator.VerifyRolesArePresent(store.Clients(), role, namespace)
	}
})

var _ = gauge.Step("Verify the total number of roles in <namespace> namespace matches the table", func(namespace string) {
	rolesTable, ok := gauge.GetScenarioStore()["rolesTable"].(*models.Table)
	if !ok {
		testsuit.T.Errorf("Could not get rolesTable from scenario store")
	}
	fullOutput := cmd.MustSucceed("oc", "get", "role", "-n", namespace, "-o", "name").Stdout()
	lines := strings.Split(fullOutput, "\n")
	actualCount := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			actualCount++
		}
	}
	expectedCount := len(rolesTable.Rows)
	if actualCount != expectedCount {
		testsuit.T.Errorf("Mismatch in number of roles in namespace %s. Expected: %d (from table), Actual: %d (from oc get role)\nFull output of 'oc get role -n %s':\n%s", namespace, expectedCount, actualCount, namespace, fullOutput)
	}
})
