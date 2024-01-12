# OpenShift Pipelines E2E tests

Validation of OpenShift Pipeline releases using [Gauge](https://docs.gauge.org/getting_started/installing-gauge.html)

### Prerequisites

* [Go](https://golang.org/)

* [Gauge](https://docs.gauge.org/getting_started/installing-gauge.html?os=linux&language=python&ide=vscode)

* Clone this repository

* Need [OpenShift](https://gitlab.cee.redhat.com/tekton/plumbing/) cluster

* Download latest [OpenShift Client](https://mirror.openshift.com/pub/openshift-v4/clients/oc/latest/) for your operating system

### Installation instructions

Install Gauge to `/usr/local/bin` by running

```curl -SsL https://downloads.gauge.org/stable | sh```

or install Gauge to a `<custom path>` using

```curl -SsL https://downloads.gauge.org/stable | sh -s -- --location-[custom path]```

For other installation methods, refer to the [Gauge documentation.](https://docs.gauge.org/getting_started/installing-gauge.html)

### Install plugins

```
GO111MODULE=off gauge install go
GO111MODULE=off gauge install html-report
GO111MODULE=off gauge install screenshot
GO111MODULE=off gauge install xml-report
```

(optional)

```
GO111MODULE=off gauge install reportportal
```

## Run a specification

Refer to the [Gauge documentation](https://docs.gauge.org/execution.html) for general information about how to run specifications (aka `specs`).

## Run tests for OpenShift Pipelines

Majority of tests assume that they run on an OpenShift cluster with already installed OpenShift Pipelines operator (e.g. latest release of a nightly build).

### Operator installation, upgrade and uninstallation

Operator installation tests have to run as `admin` user

```
CATALOG_SOURCE=custom-operators CHANNEL=latest gauge run --log-level=debug --verbose --tags install specs/olm.spec
CATALOG_SOURCE=custom-operators CHANNEL=latest gauge run --log-level=debug --verbose --tags upgrade specs/olm.spec
gauge run --log-level=debug --verbose --tags uninstall specs/olm.spec
```

> Notes: 
> - `CATALOG_SOURCE` - catalog source name, `redhat-operators` for released versions, `custom-operators` for nightly builds
> - `CHANNEL` - channel to which the installation test is supposed to subscribe, e.g. `latest` or `pipelines-1.9`

### Most common test sub-suites

The following tests have to run as `admin` user

```
gauge run --log-level=debug --verbose --tags e2e specs/metrics
gauge run --log-level=debug --verbose --tags 'e2e & tls' specs/triggers/eventlistener.spec
gauge run --log-level=debug --verbose --tags e2e specs/clustertasks/clustertask.spec
gauge run --log-level=debug --verbose --tags 'e2e & linux/amd64' specs/clustertasks/clustertask-multiarch.spec
gauge run --log-level=debug --verbose --tags e2e specs/operator/rbac.spec
gauge run --log-level=debug --verbose --tags e2e specs/operator/auto-prune.spec
gauge run --log-level=debug --verbose --tags e2e specs/operator/addon.spec
```

The following tests can run as both `admin` and regular/non-admin user

```
gauge run --log-level=debug --verbose --tags e2e specs/pipelines
gauge run --log-level=debug --verbose --tags 'e2e & !tls' specs/triggers
gauge run --log-level=debug --verbose --tags disconnected-e2e specs/clustertasks/clustertask.spec
gauge run --log-level=debug --verbose --tags 'e2e & !skip_linux/amd64' specs/clustertasks/clustertask-s2i.spec
```

## Authoring a new test specification

1. Create or update a spec file in `specs` directory using `Markdown` syntax.
2. If necessary, create steps in a new or appropriate existing `Go` file in `steps` directory.
3. If necessary, create test resources in `YAML` in `testdata` directory.
4. If necessary, implement new steps using `Go` in new or appropriate existing file in `pkg` directory.

## Running tests in a container

CI system is running these tests inside a container using image [quay.io/openshift-pipeline/ci](https://quay.io/repository/openshift-pipeline/ci?tab=tags&tag=latest) built using a Dockerfile named [Dockerfile.CI](Dockerfile.CI) hosted in a this repository. 

```
cd <path_with_content_of_this_repo>
podman run --rm -it -v $KUBECONFIG:/root/.kube/config:z -v .:/root/release-tests:z -w /root/release-tests quay.io/openshift-pipeline/ci /bin/bash
gauge run ...
```

## Containerised tests

[Dockerfile](Dockerfile) is added to execute Openshift Pipelines tests in a framework that requires that setup and tests be performed from a container, e.g. interop testing executed in OpenShift CI.
