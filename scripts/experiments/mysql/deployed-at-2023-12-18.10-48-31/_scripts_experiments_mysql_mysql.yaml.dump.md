**Executed**:

```yaml
script: mysql
description: ""
path: /Users/jploughman/go/src/github.com/BlueSageSolutions/courier/scripts/experiments/mysql/mysql.yaml
sources:
    mysql-config:
        transformation: []
        data: |
            [client]
            user = "<REPLACE_USER>"
            password = "<REPLACE_PASSWORD>"
            host = "<REPLACE_HOST>"
setup:
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
run-main: true
main:
    - executable: ""
      command: secretsmanager
      description: get the database credentials
      sensitive: true
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: get-secret-value
      arguments:
        - name: secret-id
          description: ""
          value: bluedlp-sandbox-db
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
      command: .SecretString
      description: parse the secret string - step-1
      sensitive: true
      source: mysql:main:step-0
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
    - executable: jq
      command: .host
      description: parse the host
      sensitive: false
      source: mysql:main:step-1
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
    - executable: jq
      command: .username
      description: parse the username
      sensitive: false
      source: mysql:main:step-1
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
    - executable: jq
      command: .port
      description: parse the port
      sensitive: false
      source: mysql:main:step-1
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
    - executable: jq
      command: .password
      description: parse the password
      sensitive: true
      source: mysql:main:step-1
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
      description: build mysql config file
      sensitive: true
      source: mysql-config
      replacements:
        - match: <REPLACE_USER>
          replace-with: mysql:main:step-3
        - match: <REPLACE_PASSWORD>
          replace-with: mysql:main:step-5
        - match: <REPLACE_HOST>
          replace-with: mysql:main:step-2
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
    - executable: mysql
      command: ""
      description: test connection
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: ""
      arguments:
        - name: use-config
          description: ""
          value: --defaults-extra-file=/tmp/_mysql-config.json
          style: plain
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: e
          description: ""
          value: SELECT 1;
          style: short
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
run-cleanup: true
cleanup: []

```

