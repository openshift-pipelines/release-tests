PIPELINES-18
# Verify running buildah as unprivileged

Pre condition:
  * Validate Operator should be installed


## Run buildah with the userid 1000
Tags: e2e, buildah
Component: Pipelines
Level: Integration
Type: Functional
Important: Critical

Steps:
  * Create
      | S.NO | resource_dir                                            |
      |------|---------------------------------------------------------|
      | 1    | testdata/v1beta1/taskrun/buildah-as-user-1000.rbac.yaml |
      | 2    | testdata/v1beta1/taskrun/buildah-as-user-1000.yaml      |

  * Verify taskrun
      | S.NO | taskrun_run_name     | status     | check_lable_propagation |
      |------|----------------------|------------|-------------------------|
      | 1    | buildah-as-user-1000 | successful | no                      |

## Run buildah as root in container, user on the host
