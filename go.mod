module github.com/openshift-pipelines/release-tests

go 1.14

require (
	github.com/Netflix/go-expect v0.0.0-20201125194554-85d881c3777e
	github.com/getgauge-contrib/gauge-go v0.2.0
	github.com/getgauge/common v0.0.0-20200824023809-24587c106922 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/google/go-cmp v0.5.6
	github.com/openshift/api v0.0.0-20210910062324-a41d3573a3ba
	github.com/openshift/client-go v0.0.0-20210521082421-73d9475a9142
	github.com/operator-framework/api v0.10.3
	github.com/operator-framework/operator-lifecycle-manager v0.19.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.45.0 // indirect
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.45.0
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.30.0
	github.com/tektoncd/operator v0.50.2
	github.com/tektoncd/pipeline v0.27.3
	github.com/tektoncd/triggers v0.16.1
	go.opencensus.io v0.23.0
	gomodules.xyz/jsonpatch/v2 v2.2.0
	gotest.tools/v3 v3.0.3
	honnef.co/go/tools v0.0.1-2020.1.5 // indirect
	k8s.io/api v0.22.0
	k8s.io/apimachinery v0.22.0
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/pkg v0.0.0-20210827184538-2bd91f75571c
	knative.dev/test-infra v0.0.0-20200921012245-37f1a12adbd3 // indirect
)

// Pin k8s deps to 0.21.8
replace (
	k8s.io/api => k8s.io/api v0.21.8
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.8
	k8s.io/client-go => k8s.io/client-go v0.21.8
	k8s.io/code-generator => k8s.io/code-generator v0.21.8
)
