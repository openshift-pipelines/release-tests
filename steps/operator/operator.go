package olm

import (
	"fmt"
	"log"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Update TektonConfig CR to use param with name <paramName> and value <value> to <action> auto creation of <resourceType>", func(paramName, value, action, resourceType string) {
	patchData := fmt.Sprintf("{\"spec\":{\"params\":[{\"name\":\"%s\",\"value\":\"%s\"}]}}", paramName, value)
	log.Println(action, "auto creation of", resourceType)
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "patch", "TektonConfig", "config", "--type=merge", "-p", patchData).Stdout())
})

var _ = gauge.Step("Verify RBAC resources disabled successfully", func() {
	operator.ValidateRBACAfterDisable(store.Clients(), store.GetCRNames())
})

var _ = gauge.Step("Verify RBAC resources are auto created successfully", func() {
	operator.ValidateRBAC(store.Clients(), store.GetCRNames())
})

var _ = gauge.Step("Verify CA Bundle ConfigMaps are auto created successfully", func() {
	operator.ValidateCABundleConfigMaps(store.Clients(), store.GetCRNames())
})

var _ = gauge.Step("Verify CA Bundle ConfigMaps still exist", func() {
	operator.ValidateCABundleConfigMaps(store.Clients(), store.GetCRNames())
})
