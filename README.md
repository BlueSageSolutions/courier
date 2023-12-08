# courier: a proxied command controller

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