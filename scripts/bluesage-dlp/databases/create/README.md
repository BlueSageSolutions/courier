# database.create

The deployment scripts here will create the schema, application users, grant permissions and place application credentials into AWS SSM Parameter Store.

## SCRIPTS

### [create-master-secret.yaml](create-master-secret.yaml)

The terraform for the BlueSage DLP environment will create a superuser for RDS that is used to create schemas, users, and grant permissions. The credentials for this superuser are stashed in AWS Secrets Manager as part of a JSON parcel.

Execution of the `create-master-secret` script will harvest that credential and stash it on the local file system as a mySQL config file.

```sh
courier run -s ./scripts/bluesage-dlp/databases/create/create-master-secret.yaml -c 2024client -e dev
```

### [create-schema.yaml](create-schema.yaml)

The `create-schema` script will use the mysql CLI to create the new client's database:

```sh
 courier run -s ./scripts/bluesage-dlp/databases/create/create-schema.yaml -c 2024client -e dev
```


### [create-users.yaml](create-schema.yaml)

The `create-users` script will use the mysql CLI to create an application user for the client's database. It will also place the password for this user in AWS SSM parameter store at `/clients/<BSDLP_CLIENT>/<BSDLP_ENVIRONMENT>/<BSDLP_CLIENT>_<BSDLP_ENVIRONMENT>`:

```sh
 courier run -s ./scripts/bluesage-dlp/databases/create/create-users.yaml -c 2024client -e dev
```