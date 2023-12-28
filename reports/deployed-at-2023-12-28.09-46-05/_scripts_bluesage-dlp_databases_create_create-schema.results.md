**Executing**: `/Users/jploughman/go/src/github.com/BlueSageSolutions/courier/scripts/bluesage-dlp/databases/create/create-schema.yaml`

**Command**: `/usr/local/bin/aws sts get-caller-identity`

**Description**: `check suthentication status`

**Script Reference**: `create-schema:setup:step-0`

**Output**:

```json
{
    "UserId": "AROAU4IXCK7SSFXS2KCSI:jploughman@bluesageusa.com",
    "Account": "335592708069",
    "Arn": "arn:aws:sts::335592708069:assumed-role/AWSReservedSSO_it-devops_dd6d43fab80b0112/jploughman@bluesageusa.com"
}

```

**Command**: `REDACTED: Command may contain sensitive data`

**Description**: `get the database credentials`

**Script Reference**: `create-schema:setup:step-1`

**Output**:

```json
REDACTED
```

**Command**: `REDACTED: Command may contain sensitive data`

**Description**: `parse the secret string`

**Script Reference**: `create-schema:setup:step-2`

**Output**:

```json
REDACTED
```

**Command**: `/bin/cat /tmp/create-schema_setup_step-2.json | /usr/local/bin/jq .host -r`

**Description**: `parse the host`

**Script Reference**: `create-schema:setup:step-3`

**Output**:

```json
bluedlp-lower.cluster-cs9l6nkpc8yl.us-east-1.rds.amazonaws.com

```

**Command**: `/bin/cat /tmp/create-schema_setup_step-2.json | /usr/local/bin/jq .username -r`

**Description**: `parse the username`

**Script Reference**: `create-schema:setup:step-4`

**Output**:

```json
bluedlp_lower_admin

```

**Command**: `/bin/cat /tmp/create-schema_setup_step-2.json | /usr/local/bin/jq .port -r`

**Description**: `parse the port`

**Script Reference**: `create-schema:setup:step-5`

**Output**:

```json
3306

```

**Command**: `REDACTED: Command may contain sensitive data`

**Description**: `parse the password`

**Script Reference**: `create-schema:setup:step-6`

**Output**:

```json
REDACTED
```

**Command**: `REDACTED: Command may contain sensitive data`

**Description**: `build mysql config file`

**Script Reference**: `create-schema:setup:step-7`

**Output**:

```json
REDACTED
```

**Command**: `/bin/cat /tmp/create-schema.test-connection-sql.json | /usr/local/bin/mysql --defaults-extra-file=/tmp/_mysql-config.json`

**Description**: `test connection`

**Script Reference**: `create-schema:setup:step-8`

**Output**:

```json
1
1

```

**Command**: `REDACTED: Command may contain sensitive data`

**Description**: `stash secret in ssm parameter`

**Script Reference**: `create-schema:main:step-0`

**Output**:

```json
REDACTED
```

**Command**: `/bin/cat /tmp/_create-schema-sql.json | /usr/local/bin/mysql --defaults-extra-file=/tmp/_mysql-config.json`

**Description**: `create schema`

**Script Reference**: `create-schema:main:step-1`

**Output**:

```json
{}
```

**Command**: `/usr/local/bin/aws ssm delete-parameter --name /systems/mysql/lower/defaults-extra-file`

**Description**: `delete parameter`

**Script Reference**: `create-schema:cleanup:step-0`

**Output**:

```json
{}
```

**Command**: `/bin/cat /tmp/_drop-schema-sql.json | /usr/local/bin/mysql --defaults-extra-file=/tmp/_mysql-config.json`

**Description**: `drop schema`

**Script Reference**: `create-schema:cleanup:step-1`

**Output**:

```json
{}
```

