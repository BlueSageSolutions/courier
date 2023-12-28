**Executed**:

```yaml
script: create-users
description: ""
execution-context: null
path: /Users/jploughman/go/src/github.com/BlueSageSolutions/courier/scripts/bluesage-dlp/databases/create/create-users.yaml
sources:
    create-user-sql:
        transformations: []
        data: |
            CREATE USER <BSDLP_CLIENT>_<BSDLP_ENVIRONMENT> IDENTIFIED BY '<REPLACE_PASSWORD>';
    drop-user-sql:
        transformations: []
        data: |
            DROP USER <BSDLP_CLIENT>_<BSDLP_ENVIRONMENT>;
    grant-superuser-privileges-sql:
        transformations: []
        data: |
            GRANT ALL PRIVILEGES ON <BSDLP_CLIENT>_<BSDLP_ENVIRONMENT>.* TO dlp_fly;
    grant-user-privileges-sql:
        transformations: []
        data: |
            GRANT ALL PRIVILEGES ON <BSDLP_CLIENT>_<BSDLP_ENVIRONMENT>.* TO <BSDLP_CLIENT>_<BSDLP_ENVIRONMENT>;
    password-stub:
        transformations: []
        data: |
            <REPLACE_PASSWORD>
    revoke-superuser-privileges-sql:
        transformations: []
        data: |
            REVOKE ALL PRIVILEGES, GRANT OPTION FROM dlp_fly; FLUSH PRIVILEGES;
    revoke-user-privileges-sql:
        transformations: []
        data: |
            REVOKE ALL PRIVILEGES, GRANT OPTION FROM <BSDLP_CLIENT>_<BSDLP_ENVIRONMENT>; FLUSH PRIVILEGES;
    test-connection-sql:
        transformations: []
        data: |
            SELECT 1;
setup:
    - script-reference: ""
      executable: ""
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
    - script-reference: ""
      executable: echo
      command: ""
      description: create password
      sensitive: false
      source: password-stub
      replacements:
        - match: <REPLACE_PASSWORD>
          replace-with: ""
          replace-with-random: 12
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
    - script-reference: ""
      executable: cat
      command: /tmp/_password-stub.json
      description: ""
      sensitive: false
      source: ""
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
run-main: false
main:
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
          source: create-users:setup:step-2
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
      description: create user
      sensitive: false
      source: create-user-sql
      replacements:
        - match: <BSDLP_CLIENT>
          replace-with: __CLIENT__
          replace-with-random: 0
        - match: <BSDLP_ENVIRONMENT>
          replace-with: __ENVIRONMENT__
          replace-with-random: 0
        - match: <REPLACE_PASSWORD>
          replace-with: create-users:setup:step-4
          replace-with-random: 0
      environment: []
      directory: ""
      sub-command: ""
      arguments:
        - name: use-config
          description: ""
          randomize: 0
          value: --defaults-extra-file=/tmp/create-users_main_step-0.json
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
    - script-reference: ""
      executable: mysql
      command: ""
      description: grant user all privileges
      sensitive: false
      source: grant-user-privileges-sql
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
          value: --defaults-extra-file=/tmp/create-users_main_step-0.json
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
    - script-reference: ""
      executable: mysql
      command: ""
      description: grant superuser all privileges
      sensitive: false
      source: grant-superuser-privileges-sql
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
          value: --defaults-extra-file=/tmp/create-users_main_step-0.json
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
    - script-reference: ""
      executable: mysql
      command: ""
      description: drop user
      sensitive: false
      source: drop-user-sql
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
          value: --defaults-extra-file=/tmp/create-users_main_step-0.json
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

