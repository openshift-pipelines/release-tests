# release-tests
Validation of OpenShift Pipeline releases


## Setup

### Prerequisite 

#### OCP cluster (4.2.*)

#### Installing `oc`
Download `oc` latest binary executable for your operating system


### Environment Varaiable

```
export TKN_VERSION=0.6.0
```

### Build tkn binary
Ginkgo tests uses tkn Binary created by using below command
```
 make download-tkn
 ```

### Getting Ginkgo (optional)

```
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega/...
```

ref: http://onsi.github.io/ginkgo/#getting-ginkgo

### Run tests using `ginkgo`

### Run Single Test case 
```
ginkgo -v ./spec/... --timeout 20m
```

### Run Single Test case 
```
ginkgo -v ./spec/features/pipelines -focus "Run sample pipeline" --count=1 --timeout 10m
```

### Run Olm Test Suite
```
ginkgo -v ./spec/olm/...  --count=1
```

### Run Pipeline Test Suite
```
ginkgo -v ./spec/features/pipelines --count=1 
```

### lint Test

```
make lint
```

## Organisation

`spec` directory contains Features to Tests/Test Suites against Tekton project in BDD style. Any validation/automation
 of a particular feature/component will need to be in the `pkg` directory


### Spec directory


`spec` directory is divided into the following
  -  olm: contains specs/Tests Related to olm features (install, upgrade & uninstall) operator
  -  features:  contains specs related to each of the features added to the cluster

### Pkg directory
`pkg` directory is divided into the following
  -  operator   : contains validation code of operator feature defined under `spec`
  -  olm        : contains olm related validation code      
  -  pipelines  : contains validation code of pipelines feature defined under `spec`
  -  rbac       : contains validation code for `RBAC` 

### Config Directory
`config` directory includes subscription yaml file
