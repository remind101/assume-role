This tool will request and set temporary credentials in your shell environment variables for a given role.

## Installation

On OS X, the best way to get it is to use homebrew:

```bash
brew install remind101/formulae/assume-role
```

If you have a working Go 1.6/1.7 environment:

```bash
$ go get -u github.com/remind101/assume-role
```

On Windows with PowerShell, you can use [scoop.sh](http://scoop.sh/)

```cmd
$ scoop bucket add extras
$ scoop install assume-role
```

## Configuration

Setup a profile for each role you would like to assume in `~/.aws/config`.

For example:

`~/.aws/config`:

```ini
[profile usermgt]
region = us-east-1

[profile stage]
# Stage AWS Account.
region = us-east-1
role_arn = arn:aws:iam::1234:role/SuperUser
source_profile = usermgt

[profile prod]
# Production AWS Account.
region = us-east-1
role_arn = arn:aws:iam::9012:role/SuperUser
mfa_serial = arn:aws:iam::5678:mfa/eric-holmes
source_profile = usermgt
```

`~/.aws/credentials`:

```ini
[usermgt]
aws_access_key_id = AKIMYFAKEEXAMPLE
aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/MYxFAKEYEXAMPLEKEY
```

Reference: https://docs.aws.amazon.com/cli/latest/userguide/cli-roles.html

In this example, we have three AWS Account profiles:

 * usermgt
 * stage
 * prod

Each member of the org has their own IAM user and access/secret key for the `usermgt` AWS Account.
The keys are stored in the `~/.aws/credentials` file.

The `stage` and `prod` AWS Accounts have an IAM role named `SuperUser`.
The `assume-role` tool helps a user authenticate (using their keys) and then assume the privilege of the `SuperUser` role, even across AWS accounts!

## Usage

Perform an action as the given IAM role:

```bash
$ assume-role stage aws iam get-user
```

The `assume-role` tool sets `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_SESSION_TOKEN` environment variables and then executes the command provided.

If the role requires MFA, you will be asked for the token first:

```bash
$ assume-role prod aws iam get-user
MFA code: 123456
```

If no command is provided, `assume-role` will output the temporary security credentials:

```bash
$ assume-role prod
export AWS_ACCESS_KEY_ID="ASIAI....UOCA"
export AWS_SECRET_ACCESS_KEY="DuH...G1d"
export AWS_SESSION_TOKEN="AQ...1BQ=="
export AWS_SECURITY_TOKEN="AQ...1BQ=="
export ASSUMED_ROLE="prod"
# Run this to configure your shell:
# eval $(assume-role prod)
```

Or windows PowerShell:
```cmd
$env:AWS_ACCESS_KEY_ID="ASIAI....UOCA"
$env:AWS_SECRET_ACCESS_KEY="DuH...G1d"
$env:AWS_SESSION_TOKEN="AQ...1BQ=="
$env:AWS_SECURITY_TOKEN="AQ...1BQ=="
$env:ASSUMED_ROLE="prod"
# Run this to configure your shell:
# assume-role.exe prod | Invoke-Expression
```

If you use `eval $(assume-role)` frequently, you may want to create a alias for it:

* zsh
```shell
alias assume-role='function(){eval $(command assume-role $@);}'
```
* bash
```shell
function assume-role { eval $( $(which assume-role) $@); }
```

## TODO

* [ ] Cache credentials.
