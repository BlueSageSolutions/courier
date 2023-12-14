# Executed

The results of: `/Users/jploughman/go/src/github.com/BlueSageSolutions/courier/scripts/exp/cloudquery-ecs.yaml`

**Executing**: `cloudquery-ecs`

**Command**: `/bin/cat /tmp/cloudquery-ecs.cloudquery-config.json | /usr/bin/base64`

**Script Reference**: `cloudquery-ecs:setup:step-0`

**Output**:

```json
a2luZDogc291cmNlIHNwZWM6CiAgbmFtZTogYXdzCiAgcGF0aDogImNsb3VkcXVlcnkvYXdzIgogIHJlZ2lzdHJ5OiAiY2xvdWRxdWVyeSIKICB2ZXJzaW9uOiAidjIyLjE5LjIiCiAgdGFibGVzOiBbImF3c19zM19idWNrZXRzIl0KICBkZXN0aW5hdGlvbnM6IFsiczMiXSAKa2luZDogZGVzdGluYXRpb24gc3BlYzoKICBuYW1lOiAiczMiCiAgcGF0aDogImNsb3VkcXVlcnkvczMiCiAgcmVnaXN0cnk6ICJjbG91ZHF1ZXJ5IgogIHZlcnNpb246ICJ2NC44LjMiCiAgd3JpdGVfbW9kZTogImFwcGVuZCIKICBzcGVjOgogICAgYnVja2V0OiBibHVlc2FnZS1kZXZvcHMtY2xvdWRxdWVyeS1idWNrZXQKICAgIHBhdGg6ICJ7e1RBQkxFfX0ve3tVVUlEfX0ucGFycXVldCIKICAgIGZvcm1hdDogInBhcnF1ZXQiCiAgICBhdGhlbmE6IHRydWUgICAgCg==

```

**Command**: `/usr/local/bin/aws sso login`

**Script Reference**: `cloudquery-ecs:setup:step-1`

**Output**:

```json
Attempting to automatically open the SSO authorization page in your default browser.
If the browser does not open or you wish to use a different device to authorize this request, open the following URL:

https://device.sso.us-east-1.amazonaws.com/

Then enter the code:

CDVK-MWXG
Successfully logged into Start URL: https://blue-sage.awsapps.com/start

```

**Command**: `/usr/local/bin/aws sts get-caller-identity`

**Script Reference**: `cloudquery-ecs:setup:step-2`

**Output**:

```json
{
    "UserId": "AROAY3WRRNAFLJJDUHKUQ:jploughman@bluesageusa.com",
    "Account": "609251059722",
    "Arn": "arn:aws:sts::609251059722:assumed-role/AWSReservedSSO_it-devops_556d65574656c287/jploughman@bluesageusa.com"
}

```

**Executed**:

```yaml
script: cloudquery-ecs
description: ""
sources:
    cloudquery-config:
        transformation: []
        data: "kind: source spec:\n  name: aws\n  path: \"cloudquery/aws\"\n  registry: \"cloudquery\"\n  version: \"v22.19.2\"\n  tables: [\"aws_s3_buckets\"]\n  destinations: [\"s3\"] \nkind: destination spec:\n  name: \"s3\"\n  path: \"cloudquery/s3\"\n  registry: \"cloudquery\"\n  version: \"v4.8.3\"\n  write_mode: \"append\"\n  spec:\n    bucket: bluesage-devops-cloudquery-bucket\n    path: \"{{TABLE}}/{{UUID}}.parquet\"\n    format: \"parquet\"\n    athena: true    \n"
    data-access-policy:
        transformation: []
        data: "{\n    \"Version\": \"2012-10-17\",\n    \"Statement\": [\n        {\n            \"Action\": [\n                \"s3:PutObject\"\n            ],\n            \"Resource\": [\n                \"arn:aws:s3:::bluesage-devops-cloudquery-bucket/*\"\n            ],\n            \"Effect\": \"Allow\"\n        },\n        {\n            \"Action\": [\n                \"s3:GetObject\"\n            ],\n            \"Effect\": \"Deny\",\n            \"NotResource\": [\n                \"arn:aws:s3:::bluesage-devops-cloudquery-bucket/*\"\n            ]\n        },\n        {\n            \"Action\": [\n                \"cloudformation:GetTemplate\",\n                \"dynamodb:GetItem\",\n                \"dynamodb:BatchGetItem\",\n                \"dynamodb:Query\",\n                \"dynamodb:Scan\",\n                \"ec2:GetConsoleOutput\",\n                \"ec2:GetConsoleScreenshot\",\n                \"ecr:BatchGetImage\",\n                \"ecr:GetAuthorizationToken\",\n                \"ecr:GetDownloadUrlForLayer\",\n                \"kinesis:Get*\",\n                \"lambda:GetFunction\",\n                \"logs:GetLogEvents\",\n                \"sdb:Select*\",\n                \"sqs:ReceiveMessage\"\n            ],\n            \"Resource\": \"*\",\n            \"Effect\": \"Deny\"\n        }\n    ]\n}      \n"
    task-definition:
        transformation: []
        data: |
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
    task-role-trust-policy:
        transformation: []
        data: "{\n  \"Version\": \"2012-10-17\",\n  \"Statement\": [\n    {\n      \"Sid\": \"\",\n      \"Effect\": \"Allow\",\n      \"Principal\": {\n        \"Service\": \"ecs-tasks.amazonaws.com\"\n      },\n      \"Action\": \"sts:AssumeRole\"\n    }\n  ]\n}      \n"
setup:
    - executable: base64
      command: ""
      description: ""
      sensitive: false
      source: cloudquery-config
      replacements: []
      environment: []
      directory: ""
      sub-command: ""
      arguments: []
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: sso
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: login
      arguments: []
      sleep:
        timeout: 0
        before: 0
        after: 10
        after-message: ""
        before-message: ""
    - executable: ""
      command: sts
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: get-caller-identity
      arguments: []
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
run-main: false
main:
    - executable: ""
      command: s3api
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: create-bucket
      arguments:
        - name: bucket
          description: ""
          value: bluesage-devops-cloudquery-bucket
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: region
          description: ""
          value: us-east-2
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: create-bucket-configuration
          description: ""
          value: LocationConstraint=us-east-2
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: ecs
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: create-cluster
      arguments:
        - name: cluster-name
          description: ""
          value: bluesage-devops-cloudquery-ecs
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: region
          description: ""
          value: us-east-2
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: logs
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: create-log-group
      arguments:
        - name: log-group-name
          description: ""
          value: bluesage-devops-cloudquery-logs
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: region
          description: ""
          value: us-east-2
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: logs
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: put-retention-policy
      arguments:
        - name: log-group-name
          description: ""
          value: bluesage-devops-cloudquery-logs
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: retention-in-days
          description: ""
          value: "14"
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: iam
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: create-role
      arguments:
        - name: role-name
          description: ""
          value: bluesage-devops-cloudquery-task-role
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: assume-role-policy-document
          description: ""
          value: ""
          style: ""
          quote-type: ""
          source-type: file
          source: task-role-trust-policy
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: iam
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: put-role-policy
      arguments:
        - name: role-name
          description: ""
          value: bluesage-devops-cloudquery-task-role
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: policy-name
          description: ""
          value: bluesage-devops-cloudquery-task-policy
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: policy-document
          description: ""
          value: ""
          style: ""
          quote-type: ""
          source-type: file
          source: data-access-policy
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: iam
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: attach-role-policy
      arguments:
        - name: role-name
          description: ""
          value: bluesage-devops-cloudquery-task-role
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: policy-arn
          description: ""
          value: arn:aws:iam::aws:policy/ReadOnlyAccess
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: iam
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: attach-role-policy
      arguments:
        - name: role-name
          description: ""
          value: bluesage-devops-cloudquery-task-role
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: policy-arn
          description: ""
          value: arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: jq
      command: .Role.Arn
      description: ""
      sensitive: false
      source: cloudquery-ecs:main:step-4
      replacements: []
      environment: []
      directory: ""
      sub-command: -r
      arguments: []
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: cat
      command: ""
      description: ""
      sensitive: false
      source: task-definition
      replacements:
        - match: <REPLACE_WITH_CQ_BASE64_ENCODED_CONFIG>
          replace-with: cloudquery-ecs:setup:step-0
        - match: <REPLACE_WITH_TASK_ROLE_ARN>
          replace-with: cloudquery-ecs:main:step-8
      environment: []
      directory: ""
      sub-command: ""
      arguments: []
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
run-cleanup: false
cleanup:
    - executable: ""
      command: iam
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: detach-role-policy
      arguments:
        - name: role-name
          description: ""
          value: bluesage-devops-cloudquery-task-role
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: policy-arn
          description: ""
          value: arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: iam
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: detach-role-policy
      arguments:
        - name: role-name
          description: ""
          value: bluesage-devops-cloudquery-task-role
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: policy-arn
          description: ""
          value: arn:aws:iam::aws:policy/ReadOnlyAccess
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: iam
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: delete-role-policy
      arguments:
        - name: role-name
          description: ""
          value: bluesage-devops-cloudquery-task-role
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: policy-name
          description: ""
          value: bluesage-devops-cloudquery-task-policy
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: iam
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: delete-role
      arguments:
        - name: role-name
          description: ""
          value: bluesage-devops-cloudquery-task-role
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: logs
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: delete-retention-policy
      arguments:
        - name: log-group-name
          description: ""
          value: bluesage-devops-cloudquery-logs
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: logs
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: delete-log-group
      arguments:
        - name: log-group-name
          description: ""
          value: bluesage-devops-cloudquery-logs
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: ecs
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: delete-cluster
      arguments:
        - name: cluster
          description: ""
          value: bluesage-devops-cloudquery-ecs
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: ""
      command: s3api
      description: ""
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: delete-bucket
      arguments:
        - name: bucket
          description: ""
          value: bluesage-devops-cloudquery-bucket
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""

```

