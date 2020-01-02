# release-tests
Validation of OpenShift Pipeline releases


## Setup

### Prerequisite 

#### OCP cluster (4.2.*)

#### Installing `tkn`
Download the latest binary executable for your operating system:

[tkn install](https://github.com/tektoncd/cli/blob/master/README.md#installing-tkn)

#### Installing `oc`
Download `oc` latest binary executable for your operating system

#### Environment varaibales
```
Eg: export OPERATOR_VERSION=v0.9.1

To Specify operator version, which you want install from canary channel
```


### Install goconvey

```
go get github.com/smartystreets/goconvey
```

ref: https://github.com/smartystreets/goconvey#installation

## Running tests

```
goconvey
```

### Running tests on commandline

```
go test -v ./spec
```

### Running Operator installation Test
```
go test -v --count=1 ./...  --run ^TestFreshInstall$
```

### Running Pipeline Test
```
go test -v --count=1 ./...  --run ^TestSamplePipelineRun$
```

## Organisation

`spec` directory contains only specification / BDD. Any validation/automation
 of a particular feature/component will need to be in the `pkg` directory


### Spec directory


`spec` directory is divided into the following
  -  install:  contains installation related specs
  -  upgrade:  contains upgradation related specs
  -  uninstall: contains unistallation/cleanup related specs
  -  features:  contains specs related to each of the features added to the
     cluster

### Pkg directory
`pkg` directory is divided into the following
  -  operator:  contains validation code of operator feature defined under `spec`
  -  pipelines:  contains validation code of pipelines feature defined under `spec`
  -  rbac: contains validation code for `RBAC` 

### Config Directory
`config` directory includes subscription yaml file
 