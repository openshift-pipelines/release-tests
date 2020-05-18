package assert

import (
	"github.com/getgauge-contrib/gauge-go/testsuit"
)

// NoError confirms the error returned is null
func NoError(err error, description string) {
	if err != nil {
		testsuit.T.Errorf("%s, \n err: %v", description, err)
	}
}
