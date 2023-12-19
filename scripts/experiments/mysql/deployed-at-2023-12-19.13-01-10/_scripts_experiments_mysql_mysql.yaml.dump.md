**Executed**:

```yaml
script: mysql
description: ""
execution-context: null
path: /Users/jploughman/go/src/github.com/BlueSageSolutions/courier/scripts/experiments/mysql/mysql.yaml
sources:
    create-schema:
        transformation: []
        data: |
            CREATE DATABASE <REPLACE_CLIENT>_<REPLACE_ENVIRONMENT> CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    drop-schema:
        transformation: []
        data: |
            DROP DATABASE <REPLACE_CLIENT>_<REPLACE_ENVIRONMENT>;
    mysql-config:
        transformation: []
        data: |
            [client]
            user = "<REPLACE_USER>"
            password = "<REPLACE_PASSWORD>"
            host = "<REPLACE_HOST>"
    test-connection:
        transformation: []
        data: |
            SELECT 1;
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
run-main: false
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
          value: bluedlp-lower-db
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
      source: test-connection
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
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - executable: mysql
      command: ""
      description: create schema
      sensitive: false
      source: create-schema
      replacements:
        - match: <REPLACE_CLIENT>
          replace-with: __CLIENT__
        - match: <REPLACE_ENVIRONMENT>
          replace-with: __ENVIRONMENT__
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
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
run-cleanup: false
cleanup:
    - executable: mysql
      command: ""
      description: drop schema
      sensitive: false
      source: drop-schema
      replacements:
        - match: <REPLACE_CLIENT>
          replace-with: __CLIENT__
        - match: <REPLACE_ENVIRONMENT>
          replace-with: __ENVIRONMENT__
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
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""

```

