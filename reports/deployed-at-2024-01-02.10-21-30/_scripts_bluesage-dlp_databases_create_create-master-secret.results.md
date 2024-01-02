**Executing**: `/Users/jploughman/go/src/github.com/BlueSageSolutions/courier/scripts/bluesage-dlp/databases/create/create-master-secret.yaml`

**Command**: `/usr/local/bin/aws sts get-caller-identity`

**Description**: `check aws authentication status`

**Script Reference**: `create-master-secret:setup:step-0`

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

**Script Reference**: `create-master-secret:setup:step-1`

**Output**:

```json
REDACTED
```

**Command**: `REDACTED: Command may contain sensitive data`

**Description**: `parse the secret string`

**Script Reference**: `create-master-secret:setup:step-2`

**Output**:

```json
REDACTED
```

**Command**: `/bin/cat /tmp/create-master-secret_setup_step-2.json | /usr/local/bin/jq .host -r`

**Description**: `parse the host`

**Script Reference**: `create-master-secret:setup:step-3`

**Output**:

```json
bluedlp-lower.cluster-cs9l6nkpc8yl.us-east-1.rds.amazonaws.com

```

**Command**: `/bin/cat /tmp/create-master-secret_setup_step-2.json | /usr/local/bin/jq .username -r`

**Description**: `parse the username`

**Script Reference**: `create-master-secret:setup:step-4`

**Output**:

```json
bluedlp_lower_admin

```

**Command**: `/bin/cat /tmp/create-master-secret_setup_step-2.json | /usr/local/bin/jq .port -r`

**Description**: `parse the port`

**Script Reference**: `create-master-secret:setup:step-5`

**Output**:

```json
3306

```

**Command**: `REDACTED: Command may contain sensitive data`

**Description**: `parse the password`

**Script Reference**: `create-master-secret:setup:step-6`

**Output**:

```json
REDACTED
```

**Command**: `REDACTED: Command may contain sensitive data`

**Description**: `build mysql config file`

**Script Reference**: `create-master-secret:setup:step-7`

**Output**:

```json
REDACTED
```

**Command**: `REDACTED: Command may contain sensitive data`

**Description**: `stash secret in ssm parameter`

**Script Reference**: `create-master-secret:main:step-0`

**Output**:

```json
REDACTED
```

