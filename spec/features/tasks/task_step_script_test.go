package tasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//Simplifies execution of scripts inside container, Allows you to execute bash scripts
//Eg: script: |
//#!/usr/bin/env bash
//echo "Hello from Bash!"
func TestTaskRunWithStepAsScript(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("I should create a Task with bash script", func() {
					Convey("Then I should Run Task successfully", nil)
				})
			})
		})
	})
}
