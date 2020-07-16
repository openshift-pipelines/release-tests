module github.com/openshift-pipelines/release-tests

go 1.13

require (
	github.com/Netflix/go-expect v0.0.0-20190729225929-0e00d9168667
	github.com/getgauge-contrib/gauge-go v0.1.4
	github.com/getgauge/common v0.0.0-20200429105102-5b0a7c1a1bd6 // indirect
	github.com/openshift/api v3.9.1-0.20190924102528-32369d4db2ad+incompatible
	github.com/openshift/client-go v0.0.0-20190923180330-3b6373338c9b
	github.com/tektoncd/operator v0.0.0-20200505103736-ab3f9da795f4
	github.com/tektoncd/pipeline v0.12.1
	github.com/tektoncd/triggers v0.5.0
	go.opencensus.io v0.22.2
	gomodules.xyz/jsonpatch/v2 v2.1.0
	gotest.tools/v3 v3.0.2
	k8s.io/api v0.18.2
	k8s.io/apiextensions-apiserver v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	knative.dev/pkg v0.0.0-20200306230727-a56a6ea3fa56
	sigs.k8s.io/controller-runtime v0.5.2
)

replace (
	k8s.io/api => k8s.io/api v0.16.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.5
	k8s.io/client-go => k8s.io/client-go v0.16.5
)
