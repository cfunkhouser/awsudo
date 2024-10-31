# awsudo

Execute a command locally with AWS assumed role credentials.

## Usage

`awsudo` requires a valid AWS configuration containing credentials capable of
assuming the specified role.

```console
$ awsudo
awsudo - execute commands with AWS assumed role credentials
Usage:
  -p, --profile string   The AWS configuration profile to load. Empty for default.
  -u, --role string      ARN for the AWS Role to assume for execution.
  -S, --session string   The role session name to be used during role assumption. (default "awsudo")
```

To test that it's working, get the current caller ID from the AWS CLI:

```console
$ awsudo -u arn:aws:iam::867530900042:role/SomethingAccess -- \
    aws sts get-caller-identity
{
    "UserId": "AAAAAAAAAAAAAAAAAAA42:awsudo",
    "Account": "867530900042",
    "Arn": "arn:aws:sts::867530900042:assumed-role/SomethingAccess/awsudo"
}
```
