module github.com/openshift-pipelines/release-tests

go 1.14

require (
	github.com/Netflix/go-expect v0.0.0-20201125194554-85d881c3777e
	github.com/getgauge-contrib/gauge-go v0.2.0
	github.com/getgauge/common v0.0.0-20200824023809-24587c106922 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/google/go-cmp v0.5.4
	github.com/openshift/api v0.0.0-20200331152225-585af27e34fd
	github.com/openshift/client-go v0.0.0-20200326155132-2a6cd50aedd0
	github.com/operator-framework/api v0.3.20
	github.com/operator-framework/operator-lifecycle-manager v0.17.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.45.0 // indirect
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.45.0
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/common v0.15.0
	github.com/tektoncd/operator v0.15.2-1.0.20201218101805-8934fc40c87c
	github.com/tektoncd/pipeline v0.19.0
	github.com/tektoncd/triggers v0.10.2
	go.opencensus.io v0.22.5
	gomodules.xyz/jsonpatch/v2 v2.1.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.18.12
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v12.0.0+incompatible
	knative.dev/pkg v0.0.0-20201218185703-e41409af6cff
)

// Pin k8s deps to 0.18.8
replace (
	k8s.io/api => k8s.io/api v0.18.8
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.8
	k8s.io/client-go => k8s.io/client-go v0.18.8
	k8s.io/code-generator => k8s.io/code-generator v0.18.8
)
