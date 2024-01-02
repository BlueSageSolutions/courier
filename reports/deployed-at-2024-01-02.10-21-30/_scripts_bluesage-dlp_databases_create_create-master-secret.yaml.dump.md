**Executed**:

```yaml
script: create-master-secret
description: ""
execution-context: null
path: /Users/jploughman/go/src/github.com/BlueSageSolutions/courier/scripts/bluesage-dlp/databases/create/create-master-secret.yaml
sources:
    mysql-config:
        transformations: []
        data: |
            [client]
            user = "<REPLACE_USER>"
            password = "<REPLACE_PASSWORD>"
            host = "<REPLACE_HOST>"
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
    - script-reference: ""
      executable: jq
      command: .SecretString
      description: parse the secret string
      sensitive: true
      source: create-master-secret:setup:step-1
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
      executable: jq
      command: .host
      description: parse the host
      sensitive: false
      source: create-master-secret:setup:step-2
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
      executable: jq
      command: .username
      description: parse the username
      sensitive: false
      source: create-master-secret:setup:step-2
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
      executable: jq
      command: .port
      description: parse the port
      sensitive: false
      source: create-master-secret:setup:step-2
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
      executable: jq
      command: .password
      description: parse the password
      sensitive: true
      source: create-master-secret:setup:step-2
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
      executable: cat
      command: ""
      description: build mysql config file
      sensitive: true
      source: mysql-config
      replacements:
        - match: <REPLACE_USER>
          replace-with: create-master-secret:setup:step-4
          replace-with-random: 0
        - match: <REPLACE_PASSWORD>
          replace-with: create-master-secret:setup:step-6
          replace-with-random: 0
        - match: <REPLACE_HOST>
          replace-with: create-master-secret:setup:step-3
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
run-main: false
main:
    - script-reference: ""
      executable: ""
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
          source: create-master-secret:setup:step-7
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
      executable: ""
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

```

