# release-tests
Validation of OpenShift Pipeline releases


## Setup

Install `goconvey`

## Running tests

```
goconvey
```

### Running tests on commandline

```
go test -v
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

