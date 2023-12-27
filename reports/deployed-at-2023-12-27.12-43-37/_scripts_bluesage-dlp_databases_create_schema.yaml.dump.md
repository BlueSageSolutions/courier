**Executed**:

```yaml
script: create-users
description: ""
execution-context: null
path: /Users/jploughman/go/src/github.com/BlueSageSolutions/courier/scripts/bluesage-dlp/databases/create/schema.yaml
sources:
    create-schema:
        transformations: []
        data: |
            CREATE DATABASE <BSDLP_CLIENT>_<BSDLP_ENVIRONMENT> CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    drop-schema:
        transformations: []
        data: |
            DROP DATABASE <BSDLP_CLIENT>_<BSDLP_ENVIRONMENT>;
    mysql-config:
        transformations: []
        data: |
            [client]
            user = "<REPLACE_USER>"
            password = "<REPLACE_PASSWORD>"
            host = "<REPLACE_HOST>"
    test-connection:
        transformations: []
        data: |
            SELECT 1;
setup:
    - executable: ""
      command: sts
      description: check suthentication status
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
      description: parse the secret string
      sensitive: true
      source: create-users:setup:step-1
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
      source: create-users:setup:step-2
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
      source: create-users:setup:step-2
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
      source: create-users:setup:step-2
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
      source: create-users:setup:step-2
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
          replace-with: create-users:setup:step-4
        - match: <REPLACE_PASSWORD>
          replace-with: create-users:setup:step-6
        - match: <REPLACE_HOST>
          replace-with: create-users:setup:step-3
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
    - executable: ""
      command: ssm
      description: stash secret in ssm parameter
      sensitive: true
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: put-parameter
      arguments:
        - name: name
          description: ""
          value: /systems/mysql/lower/defaults-extra-file
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: type
          description: ""
          value: SecureString
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: overwrite
          description: ""
          value: ""
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: value
          description: ""
          value: ""
          style: ""
          quote-type: ""
          source-type: text
          source: create-users:setup:step-7
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
run-main: false
main:
    - executable: mysql
      command: ""
      description: create schema
      sensitive: false
      source: create-schema
      replacements:
        - match: <BSDLP_CLIENT>
          replace-with: __CLIENT__
        - match: <BSDLP_ENVIRONMENT>
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
run-cleanup: true
cleanup:
    - executable: ""
      command: ssm
      description: delete parameter
      sensitive: false
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: delete-parameter
      arguments:
        - name: name
          description: ""
          value: /systems/mysql/lower/defaults-extra-file
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
    - executable: mysql
      command: ""
      description: drop schema
      sensitive: false
      source: drop-schema
      replacements:
        - match: <BSDLP_CLIENT>
          replace-with: __CLIENT__
        - match: <BSDLP_ENVIRONMENT>
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

