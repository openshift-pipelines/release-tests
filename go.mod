module github.com/openshift-pipelines/release-tests

require (
	contrib.go.opencensus.io/exporter/prometheus v0.1.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.12.8 // indirect
	github.com/openshift/api v3.9.1-0.20190911180052-9f80b7806f58+incompatible
	github.com/openshift/client-go v0.0.0-20190813201236-5a5508328169
	github.com/operator-framework/operator-sdk v0.10.1
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/pflag v1.0.3
	github.com/tektoncd/operator v0.0.0-20191212145541-b93217f690fe
	github.com/tektoncd/pipeline v0.9.2
	gotest.tools v2.2.0+incompatible
	gotest.tools/v3 v3.0.0
	k8s.io/api v0.0.0-20190612125737-db0771252981
	k8s.io/apiextensions-apiserver v0.0.0-20190820104113-47893d27d7f7
	k8s.io/apimachinery v0.0.0-20190612125636-6a5db36e93ad
	k8s.io/client-go v11.0.0+incompatible
	knative.dev/pkg v0.0.0-20191216211902-b26ddf762bc9
	sigs.k8s.io/controller-runtime v0.1.12
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190228180357-d002e88f6236
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190228174230-b40b2a5939e4
)

replace (
	k8s.io/kube-state-metrics => k8s.io/kube-state-metrics v1.6.0
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.12
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.1.11-0.20190411181648-9d55346c2bde

)

go 1.13
