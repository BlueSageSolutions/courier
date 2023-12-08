# courier: a proxied command controller
- [courier: a proxied command controller](#courier-a-proxied-command-controller)
  - [The Deployment Script](#the-deployment-script)
    - [3 Stages: Setup, Main, Cleanup](#3-stages-setup-main-cleanup)
  - [Pipelining](#pipelining)
    - [Example](#example)
      - [Derive from template](#derive-from-template)
  - [Installation](#installation)
  - [Running](#running)

What is a *proxied command controller*? It is basically a script runner that allows BlueSage DevSecOps opinions to inform the execution of scripts and commands both prior to and after the execution of said scripts/commands.

What kinds of opinions?
* Where should secrets be sourced from and how should they be delivered to the environment?
* Are environment variables sourced properly?
* What scripts/commands need to run before or after others?
* What scripts/commands can be run in parallel?

We do this by creating a *deployment script*: a `yaml` file that describes the pipeline for execution of scripts/commands. The *deployment script* is interpretted by `courier` which allows the proxying and intermediation of command execution.

`courier` is a very simple tool that can run anywhere: there is no setup, there is no configuration. `courier` allows us to change opinions over time and gradually automate manual processes. `courier` generates a markdown file that describes the results of each run in human-readable form. These *reports* can be archived and reviewed to see exactly what each deployment did.

Since `courier` is a security tool, it can redact any sensitive data from reports and the environment.

## The Deployment Script

A few examples might serve until more comprehensive documentation is developed.

### 3 Stages: Setup, Main, Cleanup

Here is an example of the `setup` section of a deployment script. Note: one of the *opinions* is that AWS is the default executable. Any executable can be specified, but if one is not, then `courier` assumes you mean the `awscli`:

```yaml
  setup:
    - command: sts
      sub-command: get-caller-identity
```

In this case, the purpose of `setup` is to short circuit the script if there are no valid credentials.

Here is an example of the `main` section of a script. This will create a bucket, setup ECS, and create some IAM roles:

```yaml
  main:
    - command: s3api
      sub-command: create-bucket
      arguments:
        - name: bucket
          value: bluesage-devops-cloudquery-bucket
        - name: region
          value: us-east-2
        - name: create-bucket-configuration
          value: LocationConstraint=us-east-2
    - command: ecs
      sub-command: create-cluster
      arguments:
        - name: cluster-name
          value: bluesage-devops-cloudquery-ecs
        - name: region
          value: us-east-2
    - command: logs
      sub-command: create-log-group
      arguments:
        - name: log-group-name
          value: bluesage-devops-cloudquery-logs
        - name: region
          value: us-east-2
    - command: logs
      sub-command: put-retention-policy
      arguments:
        - name: log-group-name
          value: bluesage-devops-cloudquery-logs
        - name: retention-in-days
          value: 14
    - command: iam
      sub-command: create-role
      arguments:
        - name: role-name
          value: bluesage-devops-cloudquery-task-role
        - name: assume-role-policy-document
          source: task-role-trust-policy
          source-type: file
    - command: iam
      sub-command: put-role-policy
      arguments:
        - name: role-name
          value: bluesage-devops-cloudquery-task-role
        - name: policy-name
          value: bluesage-devops-cloudquery-task-policy
        - name: policy-document
          source: data-access-policy
          source-type: file
    - command: iam
      sub-command: attach-role-policy
      arguments:
        - name: role-name
          value: bluesage-devops-cloudquery-task-role
        - name: policy-arn
          value: arn:aws:iam::aws:policy/ReadOnlyAccess
    - command: iam
      sub-command: attach-role-policy
      arguments:
        - name: role-name
          value: bluesage-devops-cloudquery-task-role
        - name: policy-arn
          value: arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy
```

Here is an example of the `cleanup` section of the `main` script above. It will destroy/detach all the things created/attached in reverse order.

```yaml
  cleanup:
    - command: iam
      sub-command: detach-role-policy
      arguments:
        - name: role-name
          value: bluesage-devops-cloudquery-task-role
        - name: policy-arn
          value: arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy
    - command: iam
      sub-command: detach-role-policy
      arguments:
        - name: role-name
          value: bluesage-devops-cloudquery-task-role
        - name: policy-arn
          value: arn:aws:iam::aws:policy/ReadOnlyAccess
    - command: iam
      sub-command: delete-role-policy
      arguments:
        - name: role-name
          value: bluesage-devops-cloudquery-task-role
        - name: policy-name
          value: bluesage-devops-cloudquery-task-policy
    - command: iam
      sub-command: delete-role
      arguments:
        - name: role-name
          value: bluesage-devops-cloudquery-task-role
    - command: logs
      sub-command: delete-retention-policy
      arguments:
        - name: log-group-name
          value: bluesage-devops-cloudquery-logs
    - command: logs
      sub-command: delete-log-group
      arguments:
        - name: log-group-name
          value: bluesage-devops-cloudquery-logs
    - command: ecs
      sub-command: delete-cluster
      arguments:
        - name: cluster
          value: bluesage-devops-cloudquery-ecs
    - command: s3api
      sub-command: delete-bucket
      arguments:
        - name: bucket
          value: bluesage-devops-cloudquery-bucket

```

## Pipelining

Deployment often consists of retrieving a deployment specific-value (e.g., a string like "UAT" or an identifier of a newly created resource) and propagating it downstream. This "propagation" is typically a match-and-replace operation on file names and contents. `courier` facilitates this by allowing the results of a previous command to be used to replace data derived along the way.

### Example

Follows is a simple example:

`courier` allows you to define `sources` of data that drive a deployment. These sources can be templates. Consider the ECS task definition:

```yaml
- script: cloudquery-ecs
  sources:
    task-definition:
      data: >
        {
          "containerDefinitions": [
            {
              "name": "ScheduledWorker",
              "image": "ghcr.io/cloudquery/cloudquery:4.2.0",
              "command": [
                "/bin/sh",
                "-c",
                "echo $CQ_CONFIG| base64 -d  > ./file.yml;/app/cloudquery sync ./file.yml --log-console --log-format json"
              ],
              "environment": [
                { "name": "CQ_CONFIG", "value": "<REPLACE_WITH_CQ_BASE64_ENCODED_CONFIG>" }
              ],
              "logConfiguration": {
                "logDriver": "awslogs",
                "options": {
                  "awslogs-group": "bluesage-devops-cloudquery-logs",
                  "awslogs-region": "us-east-2",
                  "awslogs-stream-prefix": "cloudquery-"
                }
              },
              "entryPoint": [""]
            }
          ],
          "family": "cloudquery",
          "requiresCompatibilities": ["FARGATE"],
          "cpu": "1024",
          "memory": "2048",
          "networkMode": "awsvpc",
          "taskRoleArn": "<REPLACE_WITH_TASK_ROLE_ARN>",
          "executionRoleArn": "<REPLACE_WITH_TASK_ROLE_ARN>"
        }
```

There are a few things that need to be replaced:

* `<REPLACE_WITH_CQ_BASE64_ENCODED_CONFIG>`
* `<REPLACE_WITH_TASK_ROLE_ARN>`

The ARN isn't known for a bit:

```yaml
    - command: iam
      sub-command: create-role
      arguments:
        - name: role-name
          value: bluesage-devops-cloudquery-task-role
        - name: assume-role-policy-document
          source: task-role-trust-policy
          source-type: file
```
Since we need the `"` trimmed from the ARN we need a tiny step (`step-8` of `main`)

```yaml
    - executable: jq
      sub-command: -r
      command: .Role.Arn
      source: cloudquery-ecs:main:step-4
```

And the base64 encoded config needs to be computed at setup. This is `step-0` of the `setup` section of the deployment script. (We start at `0` because we code.):

```yaml
  setup:
    - executable: base64
      source: cloudquery-config
```

#### Derive from template

Putting it together, we use the trimmed ARN and the base64 encoded config:

```yaml
    - executable: cat
      source: task-definition
      replacements:
        - match: <REPLACE_WITH_CQ_BASE64_ENCODED_CONFIG>
          replace-with: cloudquery-ecs:setup:step-0
        - match: <REPLACE_WITH_TASK_ROLE_ARN>
          replace-with: cloudquery-ecs:main:step-8
```

## Installation

This is a golang utility, so you need to install `go` to build it.

`brew install go` (If you are mac)

Then, to install `courier`:

`go install github.com/BlueSageSolutions/courier`

## Running

There is [a sample deployment script](./scripts/exp/cloudquery-ecs.yaml) that you can use. It will create and delete AWS resources in `us-east-2`, so be aware.

The following command will run all the scripts in the `./scripts/exp` asynchronously (in parallel):

```sh
courier run -s ./scripts/exp -m -c
```

[Sample results are here](./scripts/exp/deployed-at-2023-12-08.10-39-19/README.md).

