**Executing**: `/Users/jploughman/go/src/github.com/BlueSageSolutions/courier/scripts/bluesage-dlp/databases/create/create-users.yaml`

**Command**: `/usr/local/bin/aws sts get-caller-identity`

**Description**: ``

**Script Reference**: `create-users:setup:step-0`

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

**Script Reference**: `create-users:setup:step-1`

**Output**:

```json
REDACTED
```

**Command**: `REDACTED: Command may contain sensitive data`

**Description**: `grab secret and handle whitespace`

**Script Reference**: `create-users:setup:step-2`

**Output**:

```json
REDACTED
```

**Command**: `/bin/cat /tmp/_password-stub.json | /bin/echo`

**Description**: `create password`

**Script Reference**: `create-users:setup:step-3`

**Output**:

```json


```

**Command**: `/bin/cat /tmp/_password-stub.json`

**Description**: ``

**Script Reference**: `create-users:setup:step-4`

**Output**:

```json
glnNCqIThyiU

```

**Command**: `REDACTED: Command may contain sensitive data`

**Description**: `put credentials in /tmp`

**Script Reference**: `create-users:main:step-0`

**Output**:

```json
REDACTED
```

**Command**: `/bin/cat /tmp/_create-user-sql.json | /usr/local/bin/mysql --defaults-extra-file=/tmp/create-users_main_step-0.json`

**Description**: `create user`

**Script Reference**: `create-users:main:step-1`

**Output**:

```json
{}
```

**Command**: `/bin/cat /tmp/_grant-user-privileges-sql.json | /usr/local/bin/mysql --defaults-extra-file=/tmp/create-users_main_step-0.json`

**Description**: `grant user all privileges`

**Script Reference**: `create-users:main:step-2`

**Output**:

```json
{}
```

**Command**: `/bin/cat /tmp/_grant-superuser-privileges-sql.json | /usr/local/bin/mysql --defaults-extra-file=/tmp/create-users_main_step-0.json`

**Description**: `grant superuser all privileges`

**Script Reference**: `create-users:main:step-3`

**Output**:

```json
{}
```

**Command**: `/bin/cat /tmp/_drop-user-sql.json | /usr/local/bin/mysql --defaults-extra-file=/tmp/create-users_main_step-0.json`

**Description**: `drop user`

**Script Reference**: `create-users:cleanup:step-0`

**Output**:

```json
{}
```

