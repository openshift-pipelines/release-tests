package olm

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOperatorExists(t *testing.T) {
	Convey("Given new cluster", t, func() {

		Convey("When I want to install Pipelines", func() {
			Convey("I should be able to find Pipelines operator in OLM", nil)
			Convey("The version of the Pipelines operator should be 0.8", nil)
		})

	})
}
