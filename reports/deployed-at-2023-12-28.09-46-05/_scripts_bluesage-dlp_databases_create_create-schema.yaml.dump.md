**Executed**:

```yaml
script: create-schema
description: ""
execution-context: null
path: /Users/jploughman/go/src/github.com/BlueSageSolutions/courier/scripts/bluesage-dlp/databases/create/create-schema.yaml
sources:
    create-schema-sql:
        transformations: []
        data: |
            CREATE DATABASE <BSDLP_CLIENT>_<BSDLP_ENVIRONMENT> CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    drop-schema-sql:
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
    test-connection-sql:
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
          randomize: 0
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
      source: create-schema:setup:step-1
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
      source: create-schema:setup:step-2
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
      source: create-schema:setup:step-2
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
      source: create-schema:setup:step-2
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
      source: create-schema:setup:step-2
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
          replace-with: create-schema:setup:step-4
          replace-with-random: 0
        - match: <REPLACE_PASSWORD>
          replace-with: create-schema:setup:step-6
          replace-with-random: 0
        - match: <REPLACE_HOST>
          replace-with: create-schema:setup:step-3
          replace-with-random: 0
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
      source: test-connection-sql
      replacements: []
      environment: []
      directory: ""
      sub-command: ""
      arguments:
        - name: use-config
          description: ""
          randomize: 0
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
run-main: false
main:
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
          randomize: 0
          value: /systems/mysql/lower/defaults-extra-file
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: type
          description: ""
          randomize: 0
          value: SecureString
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: overwrite
          description: ""
          randomize: 0
          value: ""
          style: ""
          quote-type: ""
          source-type: ""
          source: ""
          interpolation: null
        - name: value
          description: ""
          randomize: 0
          value: ""
          style: ""
          quote-type: ""
          source-type: text
          source: create-schema:setup:step-7
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
      source: create-schema-sql
      replacements:
        - match: <BSDLP_CLIENT>
          replace-with: __CLIENT__
          replace-with-random: 0
        - match: <BSDLP_ENVIRONMENT>
          replace-with: __ENVIRONMENT__
          replace-with-random: 0
      environment: []
      directory: ""
      sub-command: ""
      arguments:
        - name: use-config
          description: ""
          randomize: 0
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
          randomize: 0
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
      source: drop-schema-sql
      replacements:
        - match: <BSDLP_CLIENT>
          replace-with: __CLIENT__
          replace-with-random: 0
        - match: <BSDLP_ENVIRONMENT>
          replace-with: __ENVIRONMENT__
          replace-with-random: 0
      environment: []
      directory: ""
      sub-command: ""
      arguments:
        - name: use-config
          description: ""
          randomize: 0
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

