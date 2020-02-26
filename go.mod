module github.com/openshift-pipelines/release-tests

require (
	contrib.go.opencensus.io/exporter/ocagent v0.6.0 // indirect
	github.com/Netflix/go-expect v0.0.0-20190729225929-0e00d9168667
	github.com/armon/go-metrics v0.0.0-20190430140413-ec5e00d3c878 // indirect
	github.com/dgryski/go-sip13 v0.0.0-20190329191031-25c5027a8c7b // indirect
	github.com/edsrzf/mmap-go v1.0.0 // indirect
	github.com/getgauge-contrib/gauge-go v0.1.3
	github.com/getgauge/common v0.0.0-20191206062403-08d97644169b // indirect
	github.com/go-openapi/analysis v0.19.4 // indirect
	github.com/go-openapi/runtime v0.19.3 // indirect
	github.com/go-openapi/strfmt v0.19.2 // indirect
	github.com/gophercloud/gophercloud v0.3.0 // indirect
	github.com/hashicorp/consul/api v1.1.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.1.0 // indirect
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/hashicorp/go-rootcerts v1.0.1 // indirect
	github.com/hashicorp/memberlist v0.1.4 // indirect
	github.com/hashicorp/serf v0.8.3 // indirect
	github.com/influxdata/influxdb v1.7.7 // indirect
	github.com/jpillora/backoff v0.0.0-20180909062703-3050d21c67d7 // indirect
	github.com/miekg/dns v1.1.15 // indirect
	github.com/openshift/api v3.9.1-0.20190911180052-9f80b7806f58+incompatible
	github.com/openshift/client-go v3.9.0+incompatible
	github.com/opentracing-contrib/go-stdlib v0.0.0-20190519235532-cf7a6c988dc9 // indirect
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/prometheus/alertmanager v0.18.0 // indirect
	github.com/prometheus/client_golang v1.2.0 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20190810000440-0ceca61e4d75 // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/tektoncd/operator v0.0.0-20200124060450-57c720be11c6
	github.com/tektoncd/pipeline v0.10.1
	go.mongodb.org/mongo-driver v1.0.4 // indirect
	go.opencensus.io v0.22.1
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect

	gotest.tools/v3 v3.0.0
	k8s.io/api v0.17.0
	k8s.io/apiextensions-apiserver v0.0.0-20190918161926-8f644eb6e783
	k8s.io/apimachinery v0.17.1
	k8s.io/client-go v11.0.0+incompatible
	knative.dev/pkg v0.0.0-20191111150521-6d806b998379
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
