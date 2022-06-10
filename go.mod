module github.com/openshift-pipelines/release-tests

go 1.16

require (
	github.com/Netflix/go-expect v0.0.0-20201125194554-85d881c3777e
	github.com/getgauge-contrib/gauge-go v0.2.0
	github.com/getgauge/common v0.0.0-20200824023809-24587c106922 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/google/go-cmp v0.5.6
	github.com/openshift/api v0.0.0-20210910062324-a41d3573a3ba
	github.com/openshift/client-go v0.0.0-20210521082421-73d9475a9142
	github.com/operator-framework/api v0.10.3
	github.com/operator-framework/operator-lifecycle-manager v0.18.3
	github.com/pkg/errors v0.9.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.45.0 // indirect
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.45.0
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.32.1
	github.com/tektoncd/operator v0.55.1
	github.com/tektoncd/pipeline v0.31.0
	github.com/tektoncd/triggers v0.18.0
	go.opencensus.io v0.23.0
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.21.8
	k8s.io/apimachinery v0.21.8
	k8s.io/client-go v0.21.8
	knative.dev/pkg v0.0.0-20211206113427-18589ac7627e
)
