# Executed

The results of: `/Users/cypherhat/go/src/github.com/BlueSageSolutions/courier/scripts/test.yaml`

**Executing**: `test`

**Command**: `/opt/homebrew/bin/aws sts get-caller-identity`

**Script Reference**: `test:setup:step-0`

**Output**:

```json
{
    "UserId": "AIDA5V4OYF3XNFKBEOO2N",
    "Account": "940360347374",
    "Arn": "arn:aws:iam::940360347374:user/experimental"
}

```

**Command**: `/bin/echo arn:aws:iam::940360347374:user/experimental`

**Script Reference**: `test:test:step-0`

**Output**:

```json
arn:aws:iam::940360347374:user/experimental

```

**Executed**:

```yaml
script: test
sources: {}
setup:
    - executable: aws
      command: sts
      sensitive: false
      source: ""
      sub-command: get-caller-identity
      arguments: []
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
main:
    - executable: echo
      command: arn.txt
      sensitive: false
      source: ""
      sub-command: ""
      arguments:
        - name: arn
          value: ""
          style: plain
          quote-type: ""
          source-type: ""
          source: test:setup:step-0
          interpolation:
            translation: null
            ephemeral: false
            name: arn
            type: jq
            parameters:
                - name: jq-query
                  value: ""
                  values:
                    - .Arn
                    - .input.Arn
                    - .outputs."Arn".[]
      sleep:
        timeout: 0
        before: 0
        after: 0
        after-message: ""
        before-message: ""
cleanup: []

```

