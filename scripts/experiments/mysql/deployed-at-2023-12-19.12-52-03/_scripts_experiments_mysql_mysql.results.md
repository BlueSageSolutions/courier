**Executing**: `/Users/jploughman/go/src/github.com/BlueSageSolutions/courier/scripts/experiments/mysql/mysql.yaml`

**Command**: `/usr/local/bin/aws sts get-caller-identity`

**Script Reference**: `mysql:setup:step-0`

**Output**:

```json
{
    "UserId": "AROAU4IXCK7SSFXS2KCSI:jploughman@bluesageusa.com",
    "Account": "335592708069",
    "Arn": "arn:aws:sts::335592708069:assumed-role/AWSReservedSSO_it-devops_dd6d43fab80b0112/jploughman@bluesageusa.com"
}

```

**Command**: `/usr/local/bin/aws secretsmanager get-secret-value --secret-id bluedlp-lower-db`

**Script Reference**: `mysql:main:step-0`

**Output**:

```json
REDACTED
```

**Command**: `/bin/cat /tmp/mysql_main_step-0.json | /usr/local/bin/jq .SecretString -r`

**Script Reference**: `mysql:main:step-1`

**Output**:

```json
REDACTED
```

**Command**: `/bin/cat /tmp/mysql_main_step-1.json | /usr/local/bin/jq .host -r`

**Script Reference**: `mysql:main:step-2`

**Output**:

```json
bluedlp-lower.cluster-cs9l6nkpc8yl.us-east-1.rds.amazonaws.com

```

**Command**: `/bin/cat /tmp/mysql_main_step-1.json | /usr/local/bin/jq .username -r`

**Script Reference**: `mysql:main:step-3`

**Output**:

```json
bluedlp_lower_admin

```

**Command**: `/bin/cat /tmp/mysql_main_step-1.json | /usr/local/bin/jq .port -r`

**Script Reference**: `mysql:main:step-4`

**Output**:

```json
3306

```

**Command**: `/bin/cat /tmp/mysql_main_step-1.json | /usr/local/bin/jq .password -r`

**Script Reference**: `mysql:main:step-5`

**Output**:

```json
REDACTED
```

**Command**: `/bin/cat /tmp/_mysql-config.json | /bin/cat`

**Script Reference**: `mysql:main:step-6`

**Output**:

```json
REDACTED
```

**Command**: `/bin/cat /tmp/mysql.test-connection.json | /usr/local/bin/mysql --defaults-extra-file=/tmp/_mysql-config.json`

**Script Reference**: `mysql:main:step-7`

**Output**:

```json
1
1

```

**Command**: `/bin/cat /tmp/_create-schema.json | /usr/local/bin/mysql --defaults-extra-file=/tmp/_mysql-config.json`

**Script Reference**: `mysql:main:step-8`

**Output**:

```json
{}
```

**Command**: `/bin/cat /tmp/_drop-schema.json | /usr/local/bin/mysql --defaults-extra-file=/tmp/_mysql-config.json`

**Script Reference**: `mysql:cleanup:step-0`

**Output**:

```json
{}
```

