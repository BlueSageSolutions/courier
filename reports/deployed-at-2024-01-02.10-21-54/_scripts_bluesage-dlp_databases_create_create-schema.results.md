**Executing**: `/Users/jploughman/go/src/github.com/BlueSageSolutions/courier/scripts/bluesage-dlp/databases/create/create-schema.yaml`

**Command**: `/usr/local/bin/aws sts get-caller-identity`

**Description**: `check aws authentication status`

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

**Description**: `grab secret and handle whitespace`

**Script Reference**: `create-schema:setup:step-2`

**Output**:

```json
REDACTED
```

**Command**: `REDACTED: Command may contain sensitive data`

**Description**: `put credentials in /tmp`

**Script Reference**: `create-schema:setup:step-3`

**Output**:

```json
REDACTED
```

**Command**: `/bin/cat /tmp/create-schema.test-connection-sql.json | /usr/local/bin/mysql --defaults-extra-file=/tmp/create-schema_setup_step-3.json`

**Description**: `test connection`

**Script Reference**: `create-schema:setup:step-4`

**Output**:

```json
1
1

```

**Command**: `/bin/cat /tmp/_create-schema-sql.json | /usr/local/bin/mysql --defaults-extra-file=/tmp/create-schema_setup_step-3.json`

**Description**: `create schema`

**Script Reference**: `create-schema:main:step-0`

**Output**:

```json
{}
```

