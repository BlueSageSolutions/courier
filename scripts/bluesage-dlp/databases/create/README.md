BlueSageDLP Database Users and Secrets

All of the following assumes the BlueSageDLP architecture using the T2 infrastructure model.

Terms:

BSDLP_ENVIRONMENT: dev, test, prod
BSDLP_CLIENT: client1
BSDLP_TIER: lower, prod

Proposal:

Use AWS Parameter Store as the gold source for DB credentials. 

Design:

Simple as possible design for effective RBAC. Secrets will be organized as follows:

/systems/system-name/tier-name/parameter-name
/clients/client-name/environment-name/parameter-name

aws ssm put-parameter --name "/systems/flyway/lower/dlp_fly" --value "your_value" --type "SecureString"

aws ssm put-parameter --name "/clients/client1/dev/client1_dev" --value "your_value" --type "SecureString"

Policy for Systems:

Example of lower tier access for flyway
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ssm:GetParameter",
        "ssm:GetParameters",
        "ssm:GetParametersByPath"
      ],
      "Resource": "arn:aws:ssm:*:*:parameter/systems/flyway/lower/*"
    }
  ]
}
This policy allows read access to any parameter under /systems/flyway/lower/.

Policy for Clients:

{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ssm:GetParameter",
        "ssm:GetParameters",
        "ssm:GetParametersByPath"
      ],
      "Resource": "arn:aws:ssm:*:*:parameter/clients/client1/dev/*"
    }
  ]
}
This policy allows read access to any parameter under /clients/client1/dev/

System Credentials

Currently, the master database user and password are stored in AWS Secrets Manager using the naming convention: bluedlp-${BSDLP_TIER}-db.

e.g., bluedlp-lower-db

In that secret, we find a user named: bluedlp_${BSDLP_TIER}_admin

e.g., bluedlp_lower_admin

We use this user to create client schemas named:

${BSDLP_CLIENT}_${BSDLP_ENVIRONMENT} for all the environments in that tier

e.g., client1_test, client1_dev

We create users who can access from any host named:

${BSDLP_CLIENT}_${BSDLP_ENVIRONMENT}@% 

e.g., client1_test, client1_dev

We grant all privileges (permissions) to each user's corresponding database:

GRANT ALL PRIVILEGES ON ${BSDLP_CLIENT}_${BSDLP_ENVIRONMENT}.* TO '${BSDLP_CLIENT}_${BSDLP_ENVIRONMENT}'@'%';

e.g., GRANT ALL PRIVILEGES ON client1_dev.* TO 'client1_dev'@'%';

GRANT USAGE ON *.* TO '${BSDLP_CLIENT}_${BSDLP_ENVIRONMENT}'@'%';

e.g., GRANT USAGE ON *.* TO 'client1_dev'@'%';

We grant all privileges (permissions) to dlp_fly - which already exists - for the corresponding databases in both dev/test:

GRANT ALL PRIVILEGES ON ${BSDLP_CLIENT}_${BSDLP_ENVIRONMENT}.* TO 'dlp_fly'@'%';

e.g., GRANT ALL PRIVILEGES ON client1_dev.* TO 'dlp_fly'@'%';

These secrets will be stored in AWS Parameter Store in the BlueSageDLP account at the following coordinates:

/systems/flyway/lower/dlp_fly
/clients/client1/dev/client_dev
/clients/client1/test/client_test

All database user credentials are also persisted as groovy files on the host. Eventually, these will be pulled from AWS Parameter Store at runtime.