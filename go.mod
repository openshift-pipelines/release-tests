module github.com/openshift-pipelines/release-tests

go 1.18

require (
	github.com/Netflix/go-expect v0.0.0-20220104043353-73e0943537d2
	github.com/getgauge-contrib/gauge-go v0.3.1
	github.com/google/go-cmp v0.5.9
	github.com/openshift/api v0.0.0-20220525145417-ee5b62754c68
	github.com/openshift/client-go v0.0.0-20220603133046-984ee5ebedcf
	github.com/operator-framework/api v0.16.0
	github.com/operator-framework/operator-lifecycle-manager v0.22.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.15.0
	github.com/prometheus/common v0.42.0
	github.com/tektoncd/operator v0.67.0
	github.com/tektoncd/pipeline v0.47.0
	github.com/tektoncd/triggers v0.24.0
	go.opencensus.io v0.24.0
	gotest.tools/v3 v3.4.0
	k8s.io/api v0.26.3
	k8s.io/apimachinery v0.26.4
	k8s.io/client-go v1.5.2
	knative.dev/pkg v0.0.0-20230320014357-4c84b1b51ee8
)

replace knative.dev/pkg => knative.dev/pkg v0.0.0-20230221145627-8efb3485adcf

replace k8s.io/api => k8s.io/api v0.25.9

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.25.9

replace k8s.io/apimachinery => k8s.io/apimachinery v0.25.9

replace k8s.io/client-go => k8s.io/client-go v0.25.9

replace k8s.io/code-generator => k8s.io/code-generator v0.25.9

require (
	contrib.go.opencensus.io/exporter/ocagent v0.7.1-0.20200907061046-05415f1de66d // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.4.2 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr v1.4.10 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/blendle/zapdriver v1.3.1 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cloudevents/sdk-go/v2 v2.14.0 // indirect
	github.com/creack/pty v1.1.17 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dmotylev/goproperties v0.0.0-20140630191356-7cbffbaada47 // indirect
	github.com/emicklei/go-restful/v3 v3.10.2 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/getgauge/common v0.0.0-20200824023809-24587c106922 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/cel-go v0.13.0 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-containerregistry v0.14.1-0.20230409045903-ed5c185df419 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.14 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/manifestival/manifestival v0.7.2 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/openshift-pipelines/pipelines-as-code v0.19.1 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/prometheus/statsd_exporter v0.23.1 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stoewer/go-strcase v1.2.1 // indirect
	github.com/tidwall/gjson v1.12.1 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/sjson v1.2.4 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/oauth2 v0.7.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/term v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/api v0.119.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.55.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.26.3 // indirect
	k8s.io/klog/v2 v2.90.1 // indirect
	k8s.io/kube-openapi v0.0.0-20230308215209-15aac26d736a // indirect
	k8s.io/utils v0.0.0-20230313181309-38a27ef9d749 // indirect
	knative.dev/eventing v0.36.7 // indirect
	sigs.k8s.io/controller-runtime v0.14.6 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
