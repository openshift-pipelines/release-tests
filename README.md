# release-tests
Validation of OpenShift Pipeline releases using [Gauge](https://docs.gauge.org/getting_started/installing-gauge.html)


### ***Prerequisite***

* [Go](https://golang.org/)

* [Gauge](https://docs.gauge.org/getting_started/installing-gauge.html?os=linux&language=python&ide=vscode)

* Clone this repository into [GOPATH](https://github.com/golang/go/wiki/GOPATH).

* Need [OCP](https://gitlab.cee.redhat.com/tekton/plumbing/) cluster (4.4 and above)

* Download [OC](https://mirror.openshift.com/pub/openshift-v4/clients/oc/latest/) latest binary executable for your operating system

#### Installation Instructions

* *Install using DNF Package Manager*

```> sudo dnf install gauge```

* *Install using Curl*

Install Gauge to `/usr/local/bin` by running

```> curl -SsL https://downloads.gauge.org/stable | sh```

Or Install Gauge to a `<custom path>` using

```> curl -SsL https://downloads.gauge.org/stable | sh -s -- --location-[custom path]```

### Alternative Installation Methods

* Refer [Doc.](https://docs.gauge.org/getting_started/installing-gauge.html)

* Follow the steps to add the Gauge VS Code plugin from the IDE

  * Install Gauge extension for [VS Code](https://marketplace.visualstudio.com/items?itemName=getgauge.gauge).

#### Install Plugins

* Install go plugin

```
> gauge install go
```

* Install html plugin

```
> gauge install html-report
```

* Install xml-report

```
> gauge install xml-report
```

* (optional) Install reportportal

```
> gauge install reportportal
```

## Run a specification

* You can run a Gauge specification by using the gauge run command. When this command is run, Gauge scans the directories and sub-directories at `<project_root>` (location at which the Gauge project is created) and picks up valid specification files.

```
> gauge run [args] [flags]
```
   * `<project_root>` - location at which a Gauge project is created
   * `[args]` - directories in which specifications are stored, location of specification files and scenarios
   * `[flags]` - options that can be used with this command such as --tags, -e, -f, and so on

> Note:
Gauge specifications can also be run from within the IDE (VS Code)

* Run multiple specifications
```
> gauge run <path_to_spec1> <path_to_spec2> <path_to_spec3>
```

* Run multiple directories
```
> gauge run specs test_specs
```
* Filter Specifications by `Tags`

```
> gauge run --tags "search" specs
```

* Tag Expressions (Useful stuff):


Tags                         | Selects specs/scenarios that
-----------------------------|-------------------------------
!TagA                        |do not have TagA
TagA & TagB (or) TagA,TagB   |have both TagA and TagB.
TagA & !TagB                 |have TagA and not TagB.
Tag\|TagB                    |have either TagA or TagB.
(TagA & TagB) \| TagC        |have either TagC or both TagA and TagB
!(TagA & TagB) \| TagC       |have either TagC or do not have both TagA and TagB
(TagA \| TagB) & TagC        |have either [TagA and TagC] or [TagB and TagC]

## Gauge Commands for reference
- Synopsis
  - Gauge is a light-weight cross-platform test automation tool with the ability to author test cases in the business language.

- `gauge <command> [flags] [args]`

- Examples:
```
 > gauge run specs/
 > gauge run --parallel specs/
```

  Short hand notation       | Description
----------------------------|-------------------------------
  -d, --dir string          | Set the working directory for the current command, accepts a path relative to current directory (default ".")
  -h, --help                |  help for gauge
  -l, --log-level string    | Set level of logging to debug, info, warning, error or critical (default "info")
  -m, --machine-readable    | Prints output in JSON format
  -v, --version             | Print Gauge and plugin versions

## Run openshift-pipeline tests

```
> gauge run --env "default, test" --log-level=debug --verbose specs/pipelines specs/triggers
```

## Run pipelines tests

```
> gauge run --env "default, test" --log-level=debug --verbose specs/pipelines/
```

## Run openshift-pipelines monitoring acceptance tests

```
> gauge run --env "default, test" --log-level=debug --tags "e2e" --verbose specs/metrics/
```

## Run olm tests

### Fresh installation
```
> CATALOG_SOURCE=pre-stage-operators CHANNEL=preview gauge run --env "default, test" --tags "install" --log-level=debug --verbose specs/olm.spec
```

> Notes: 
> - set `CATALOG_SOURCE` eg: `pre-stage-operators`
> - set `CHANNEL` env variable Eg: `stable`,
> - helps user to install operator by subscribing to `CHANNEL` (Assumption: pipelines operator shouldn't be installed) for `redhat-operators` or user defined catalog sources

### Upgrade operator
```
> CATALOG_SOURCE=$CATALOG_SOURCE CHANNEL=$CHANNEL gauge run --env "default, test" --tags "upgrade" --log-level=debug --verbose specs/olm.spec
```
> Notes:
> - helps user to upgrade operator by updating subscription to latest `CHANNEL` (Assumption: cluster should have pipelines operator installed)

### Uninstall Operator
```
> gauge run --env "default, test" --tags "uninstall" --log-level=debug --verbose specs/olm.spec
```
 
## Package structures

- `specs` directory contains only specification written `Markdown` syntax.

- Any validation/automation
 of a particular feature/component will need to be in the `pkg` directory

- `env` Directory where we store gauge/framework related configurations.

- `logs` directory where logs gets stored on each execution

- `reports` directory contains reports generated on each execution
- `specs` directory is divided into the following
  -  `pipelines` :  contains specs related to the component pipeline
  -  `triggers` :  contains specs related to the component triggers
  - `metrics` : contains specs related to the openshift-pipelines metrics

  -  `olm` : containse sepcs related to olm
       *  install: contains specs related to olm install operator
       *  uninstall: contains specs related to olm uninstall operator
       *  upgrade: contains specs related to olm upgrade operator

## Docker Build

- Build docker file

```
> podman build . -t <username>/release-tests:v0.1
```
> Note:
> 1. Docker file have pre-requisites like go, gauge installed to run release-tests

- It's decoupled from code changes, try to use `gitvolumes` or `pipeline-resources(git)` to pass `release-tests` code

- Needs `OC` login before we execute `gauge run` instructions

## Workaround to run docker image locally

```
> podman run --rm --it <username>/release-tests:v0.1 /bin/sh
```

```
> git clone <release-tests> repo to GOPATH
```

```
> oc login
```

```
> gauge run --env "test" --log-level=debug --verbose specs/pipelines/
```

## Dogfooding own product

- See Real usecase of `release-tests` in [CI](https://gitlab.cee.redhat.com/tekton/plumbing/ci) system
