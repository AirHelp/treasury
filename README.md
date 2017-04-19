# treasury

Treasury is a very simple tool for managing secrets. It uses Amazon S3 service to store secrets. The secrets are encrypted before saving them on disks in their data centers and decrypted when we read the secrets. Treasury uses Server-Side Encryption with AWS KMS-Managed Keys ([SSE-KMS](http://docs.aws.amazon.com/AmazonS3/latest/dev/UsingKMSEncryption.html)).

## Architecture

![Architecture overwiev](doc/Treasure_diagram_v2.png)

## Command Line interface (CLI)

Treasury is controlled via a very easy to use command-line interface (CLI). Treasury is only a single command-line application: treasury. This application takes a subcommand such as "read", "write", "import" or "export".

The Treasury CLI is a well-behaved command line application. In erroneous cases, a non-zero exit status will be returned. It also responds to -h and --help as you'd most likely expect.

To view a list of the available commands at any time, just run `treasury` with no arguments. To get help for any specific subcommand, run the subcommand with the -h argument.

** TO DO: add homebrew for cli install **

### Requirements

* Export environment variables: treasury S3 bucket and region (environment variable or --region parameter) set

For example:

```
export TREASURY_S3=st-treasury-st-staging
export AWS_REGION=eu-west-1
```

* AWS Credentials

Before using the Treasury CLI, ensure that you've configured AWS credentials. The best way to configure credentials on your machine is to use the ~/.aws/credentials file, which might look like:

```bash
[default]
aws_access_key_id = AKID1234567890
aws_secret_access_key = MY-SECRET-KEY
```

Alternatively, you can set the following environment variables:

```bash
AWS_ACCESS_KEY_ID=AKID1234567890
AWS_SECRET_ACCESS_KEY=MY-SECRET-KEY
```

You can also use non-default awscli profile:

```
AWS_PROFILE=st-staging treasury read integration/webapp/cockpit_api_pass`
```

And non-default awscli profile without default region:

```
AWS_PROFILE=st-staging ./treasury --region eu-west-1 read test/webapp/cockpit_pass`
```

* Example AWS IAM Policy

Read and Write policy to `test/test/*` and `test/cockpit/*` keys
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Stmt1491319766000",
            "Effect": "Allow",
            "Action": [
                "kms:Encrypt",
                "kms:Decrypt",
                "kms:ListAliases",
                "kms:ListKeys",
                "kms:GenerateDataKey*"
            ],
            "Resource": [
                "arn:aws:kms:eu-west-1:064764542321:key/14b4a163-6c9d-4edb-a4bf-5adc4cd50ad8"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::st-treasury-st-staging"
            ]
        },
        {
            "Sid": "Stmt1491319793000",
            "Effect": "Allow",
            "Action": [
                "s3:PutObject*",
                "s3:GetObject*"
            ],
            "Resource": [
                "arn:aws:s3:::st-treasury-st-staging/test/test/*",
                "arn:aws:s3:::st-treasury-st-staging/test/cockpit/*"
            ]
        }
    ]
}
```

Read only policy for `test/*` keys
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Stmt1491319766000",
            "Effect": "Allow",
            "Action": [
                "kms:Decrypt",
                "kms:ListAliases",
                "kms:ListKeys",
            ],
            "Resource": [
                "arn:aws:kms:eu-west-1:064764542321:key/14b4a163-6c9d-4edb-a4bf-5adc4cd50ad8"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::st-treasury-st-staging"
            ]
        },
        {
            "Sid": "Stmt1491319793000",
            "Effect": "Allow",
            "Action": [
                "s3:GetObject*"
            ],
            "Resource": [
                "arn:aws:s3:::st-treasury-st-staging/test/*"
            ]
        }
    ]
}
```

* Example AWS S3 Policy

The following bucket policy denies upload object (s3:PutObject) permission to everyone if the request does not include the `x-amz-server-side-encryption` header requesting server-side encryption with SSE-KMS.

```json
{
  "Version": "2012-10-17",
  "Id": "PutObjPolicy",
  "Statement": [
    {
      "Sid": "DenyIncorrectEncryptionHeader",
      "Effect": "Deny",
      "Principal": "*",
      "Action": "s3:PutObject",
      "Resource": "arn:aws:s3:::st-treasury-st-staging/*",
      "Condition": {
        "StringNotEquals": {
          "s3:x-amz-server-side-encryption": "aws:kms"
        }
      }
    },
    {
      "Sid": "DenyUnEncryptedObjectUploads",
      "Effect": "Deny",
      "Principal": "*",
      "Action": "s3:PutObject",
      "Resource": "arn:aws:s3:::st-treasury-st-staging/*",
      "Condition": {
        "Null": {
          "s3:x-amz-server-side-encryption": true
        }
      }
    }
  ]
}
```

### CLI Usage

#### Write secret
```
> treasury write integration/webapp/cockpit_api_pass superSecretPassword
Success! Data written to: webapp/integration/cockpit_api_pass
```

Note: if secret value is equal to existing one, write is skipped. `--force` flag can be used to overwrite.

#### Read secret
```
> treasury read integration/webapp/cockpit_api_pass
superSecretPassword
```

#### Import secrets
Assuming properties file `./secrets.env` with content:
```bash
ke1=secret1
key2=secret2
```
To import these values into s3:
```bash
> treasury import integration/application/ ./secrets.env
Import successful
```

Note: Using `=` in secret value is not allowed.
      If secret value is equal to existing one, import skips this value. `--force` flag can be used to overwrite.

#### Export secrets
Assuming stored secrets pairs on s3
```bash
integration/webapp/key1 => superSecretPassword1
integration/webapp/key2 => superSecretPassword2
```

To see exported values:
```bash
> treasury export integration/webapp/
export key1=superSecretPassword1
export key2=superSecretPassword2
```

To export them into shell environment variables:
```bash
eval $(treasury export integration/webapp/)
```

## Go Client

Example:
```go
import "github.com/AirHelp/treasury/client"

// use default client options
treasuryOptions := client.Options{}
treasury, err := client.NewClient(treasuryURL, treasuryOptions)
if err != nil {
  return err
}
secret, err := treasury.Read(key)
if err != nil {
  return err
}

fmt.Println(secret.Value)
```

** TO DO: add terraform resource **


## Development

## Build for development

```
make build
```

## Tests

Go tests

```
make test
```

Bats tests

```
bats test/bats/tests.bats
```

If `bats` missing, install it:
```bash
brew install bats
```
