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
    test-connection-sql:
        transformations: []
        data: |
            SELECT 1;
setup:
    - script-reference: ""
      executable: ""
      command: sts
      description: check aws authentication status
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
    - script-reference: ""
      executable: ""
      command: ssm
      description: get the database credentials
      sensitive: true
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: get-parameter
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
        - name: with-decryption
          description: ""
          randomize: 0
          value: ""
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
    - script-reference: ""
      executable: jq
      command: .Parameter.Value
      description: grab secret and handle whitespace
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
    - script-reference: ""
      executable: persist-arguments
      command: ""
      description: put credentials in /tmp
      sensitive: true
      source: ""
      replacements: []
      environment: []
      directory: ""
      sub-command: ""
      arguments:
        - name: password
          description: ""
          randomize: 0
          value: ""
          style: plain
          quote-type: ""
          source-type: text
          source: create-schema:setup:step-2
          interpolation: null
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
    - script-reference: ""
      executable: mysql
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
          value: --defaults-extra-file=/tmp/create-schema_setup_step-3.json
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
    - script-reference: ""
      executable: mysql
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
          value: --defaults-extra-file=/tmp/create-schema_setup_step-3.json
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
    - script-reference: ""
      executable: mysql
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
          value: --defaults-extra-file=/tmp/create-schema_setup_step-3.json
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

