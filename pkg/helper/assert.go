package helper

import (
	"github.com/getgauge-contrib/gauge-go/testsuit"
)

// AssertNoError confirms the error returned is null
func AssertNoError(err error, description string) {
	//Expect(err).ShouldNot(HaveOccurred(), description)

	if err != nil {
		testsuit.T.Errorf("%s, \n err:%s", description, err)
	}
}
