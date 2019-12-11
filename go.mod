module github.com/openshift-pipelines/release-tests

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.12.8 // indirect
	github.com/AlecAivazis/survey/v2 v2.0.5
	github.com/Netflix/go-expect v0.0.0-20190729225929-0e00d9168667
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/openshift/api v3.9.1-0.20190911180052-9f80b7806f58+incompatible
	github.com/openshift/client-go v3.9.0+incompatible
	github.com/smartystreets/goconvey v1.6.4
	github.com/tektoncd/operator v0.0.0-20191212145541-b93217f690fe
	github.com/tektoncd/pipeline v0.9.2
	go.opencensus.io v0.22.1
	k8s.io/api v0.0.0-20191004102255-dacd7df5a50b
	k8s.io/apiextensions-apiserver v0.0.0-20190918161926-8f644eb6e783
	k8s.io/apimachinery v0.0.0-20191004074956-01f8b7d1121a
	k8s.io/client-go v11.0.0+incompatible
	knative.dev/pkg v0.0.0-20190909195211-528ad1c1dd62
	sigs.k8s.io/controller-runtime v0.4.0
)

//Pinned to kubernetes-1.13.4
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
