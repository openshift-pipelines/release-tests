# release-tests
Validation of OpenShift Pipeline releases


## Setup

### ***Prerequisite*** 

#### OCP cluster (4.3.*)

#### Installing `oc`
Download `oc` latest binary executable for your operating system

#### Installing Gauge on Linux

* *Install using DNF Package Manager*

```sudo dnf install gauge```

* *Install using Curl*

Install Gauge to /usr/local/bin by running

```curl -SsL https://downloads.gauge.org/stable | sh```

Or install Gauge to a [custom path] using

```curl -SsL https://downloads.gauge.org/stable | sh -s -- --location-[custom path]```
### Alternative Gauge Installation methods

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

## Running olm install tests 
```
gauge run --env "test" --log-level=debug --verbose   specs/install.spec
```

## Run openshift-pipeline tests

```
gauge run --env "test" --log-level=debug  --verbose specs/
```

## Organisation

`specs` directory contains only specification / BDD. Any validation/automation
 of a particular feature/component will need to be in the `pkg` directory


### Spec directory


`specs` directory is divided into the following
  -  features :  contains specs related to the features tekton offers like (pipelines, cli, triggers, catalog, operator)git 
  -  olm : containse sepcs related to olm
       *  install: contains specs related to olm install operator
       *  uninstall: contains specs related to olm uninstall operator
       *  upgrade: contains specs related to olm upgrade operator