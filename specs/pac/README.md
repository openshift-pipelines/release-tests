# PAC E2E tests 
## _Pipelines-as-code_
Pipelines as code is a project allowing you to define your CI/CD using Tekton PipelineRuns and Tasks in a file located in your source control management (SCM) system, such as GitHub or GitLab. This file is then used to automatically create a pipeline for a Pull Request or a Push to a branch.

## _Settingup PAC in Gitlab_

- Create a New project in gitlab.com
- Change the visibility of the project to Public
- Set the branches as unprotected branch
- Copy the project ID by clicking on three dots in project root directory and`export GITLAB_PROJECT_ID=<ProjectID>`
- Open `User Settings --> Access tokens`
- Create a New Personal Access Token and `export GITLAB_TOKEN=<Token>`
- Create a new Group in Gitlab and Copy the URL name to `export GITLAB_GROUP_NAMESPACE=<GroupURLName>`

## Running PAC E2E tests
Export the following Env Variables
```
export GITLAB_TOKEN=<Token>
export GITLAB_PROJECT_ID=<ProjectID>
export GITLAB_GROUP_NAMESPACE=<GroupURLName>
```

To run pac e2e tests...

```
gauge run --log-level=debug --verbose  specs/pac/pac-gitlab.spec
```

