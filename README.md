# treasury

Treasury is a very simple and easy to use tool for managing secrets. It uses Amazon S3 service to store secrets. The secrets are encrypt before saving it on disks in its data centers and decrypt it when we read the secret. Treasury uses Server-Side Encryption.

## Architecture

![Architecture overwiev](doc/Treasure_diagram_v2.png)

## Command Line interface (CLI)

Treasury is controlled via a very easy to use command-line interface (CLI). Treasury is only a single command-line application: treasury. This application then takes a subcommand such as "read" or "write".

The Treasury CLI is a well-behaved command line application. In erroneous cases, a non-zero exit status will be returned. It also responds to -h and --help as you'd most likely expect.

To view a list of the available commands at any time, just run `treasury` with no arguments. To get help for any specific subcommand, run the subcommand with the -h argument.

** TO DO: add homebrew for cli install **

### Requirements

* Treasury S3 over environment variable

For example:
```
export TREASURY_S3=st-treasury-st-staging
```

* AWS Credentials

Before using the Treasury CLI, ensure that you've configured credentials. The best way to configure credentials on your machine is to use the ~/.aws/credentials file, which might look like:

[default]
aws_access_key_id = AKID1234567890
aws_secret_access_key = MY-SECRET-KEY

Alternatively, you can set the following environment variables:

AWS_ACCESS_KEY_ID=AKID1234567890
AWS_SECRET_ACCESS_KEY=MY-SECRET-KEY

### CLI Usage

Write secret
```
> treasury write integration/webapp/cockpit_api_pass superSecretPassword
Success! Data written to: webapp/integration/cockpit_api_pass
```

Read secret:
```
> treasury read integration/webapp/cockpit_api_pass
superSecretPassword
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
