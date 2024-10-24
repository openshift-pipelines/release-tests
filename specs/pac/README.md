# PAC E2E tests 
## _Pipelines-as-code_
Pipelines as code is a project allowing you to define your CI/CD using Tekton PipelineRuns and Tasks in a file located in your source control management (SCM) system, such as GitHub or GitLab. This file is then used to automatically create a pipeline for a Pull Request or a Push to a branch.

## _Settingup PAC in Gitlab_

- Create a New project in gitlab.com
- Change the visibility of the project to Public
- Set the main branch as unprotected branch
- Copy the project ID by clicking on three dots in project root directory and`export GITLAB_PROJECT_ID=<ProjectID>`
- Click on your profile under `preferences` Under `User Settings --> Access tokens`
- Create a New Personal Access Token and `export GITLAB_TOKEN=<Token>`
- Create a new Public Group in Gitlab and Copy the only the Group name from URL e.g: From GitLab URL `https://gitlab.com/groups/test324345` Copy only the group name `test324345` and `export GITLAB_GROUP_NAMESPACE=<GroupName>`
- Enter any WebhookSecret to be used for gitlab webhook `export WEBHOOK_TOKEN=<WebhookSecret>`

## Running PAC E2E tests
Export the following Env Variables
```
export GITLAB_TOKEN=<Token>
export GITLAB_PROJECT_ID=<ProjectID>
export GITLAB_GROUP_NAMESPACE=<GroupName>
export WEBHOOK_TOKEN=<WebhookSecret>
```

To run pac e2e tests...

```
gauge run --log-level=debug --verbose --tags e2e specs/pac/pac-gitlab.spec
```
