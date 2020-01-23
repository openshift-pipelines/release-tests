package helper

import (
	"bytes"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"runtime"
	"text/template"
)

// RootDir returns you the root directory of this package as string
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func CreateSubscriptionYaml(channel, installPlan, csv string) {
	var err error
	var data = struct {
		Channel     string
		InstallPlan string
		CSV         string
	}{
		Channel:     channel,
		InstallPlan: installPlan,
		CSV:         csv,
	}

	var tmplBytes bytes.Buffer

	b, err := ioutil.ReadFile(filepath.Join(RootDir(), "../config/subscription.yaml.tmp")) // just pass the file name
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("subscription").Parse(string(b))
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(&tmplBytes, data)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filepath.Join(RootDir(), "../config/subscription.yaml"), tmplBytes.Bytes(), 0777)
	// handle this error
	if err != nil {
		// print it out
		log.Fatal(err)
	}

}
