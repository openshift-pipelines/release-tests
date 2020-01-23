# release-tests
Validation of OpenShift Pipeline releases


## Setup

### ***Prerequisite*** 

#### OCP cluster (4.2.*)

#### Installing `oc`
Download `oc` latest binary executable for your operating system

### Build tkn binary
gauge tests uses tkn Binary created by using below command
```
 make download-tkn TKN_VERSION=0.6.0
 ```

### Install guauge
* [Gauge](https://docs.gauge.org/getting_started/installing-gauge.html)
* Gauge Go plugin
  * can be installed using 
  ```
  gauge install go
  ```
* Gauge html plugin

  * can be installed using 
  ```
  gauge install html-report
  ```  

  * Telemetry should be off
  ```
  gauge telemetry off
  ```
## Running olm related tests 
```
gauge run --env "test" --log-level=debug --verbose   specs/olm
```

## Run openshift-pipeline tests

```
gauge run --env "test" --log-level=debug  --verbose specs/features
```

## Organisation

`specs` directory contains only specification / BDD. Any validation/automation
 of a particular feature/component will need to be in the `pkg` directory


### Spec directory


`specs` directory is divided into the following
  -  features :  contains specs related to the features tekton offers like (pipelines, cli, triggers, catalog, operator)
  -  olm : containse sepcs related to olm
       *  install: contains specs related to olm install operator
       *  uninstall: contains specs related to olm uninstall operator
       *  upgrade: contains specs related to olm upgrade operator