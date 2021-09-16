package utility

import (
	"log"
	"strconv"
	"time"

	"github.com/getgauge-contrib/gauge-go/gauge"
)

var _ = gauge.Step("Sleep for <numberOfSeconds> seconds", func(numberOfSeconds string) {
	log.Printf("Sleeping for %v seconds", numberOfSeconds)
	numberOfSecondsInt, _ := strconv.Atoi(numberOfSeconds)
	time.Sleep(time.Duration(numberOfSecondsInt) * time.Second)
})
