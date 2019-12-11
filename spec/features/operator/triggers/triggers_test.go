package triggers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTriggers(t *testing.T) {
	Convey("Given the Operator is installed", t, func() {
		Convey("When I create a new Triggers SA with Required RBAC", nil)
		Convey("Then I create event listner instance with Trigger SA created", func() {
			Convey("Then it should create service to event listner", nil)
			Convey("I should be able expose service externally (Route/Ingress)", nil)
			Convey("I should be able to Configure webhook URL to (github/gitLab)", func() {
				Convey("Then verify event listners is able receive an event", func() {
					Convey("Validate whether Trigger Binding and Trigger Template are resolved properly", func() {
						Convey("Then Validate right resources (defined under RBAC), are created when an event occured(Eg: PipelineResources, PipelineRun)", nil)
					})
					Convey("Check whether can an interceptor intercepts the request & deny the further progress of event for a particular combination of TriggerBindings & TriggerTemplates", func() {
						Convey("If yes, It should create Pipeline Run or associated Resources under triggers Template ", nil)
					})
				})
			})
		})
	})
}
