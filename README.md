This tool will request and set temporary credentials in your shell environment variables for a given role.

## Installation

To install version 1.0.0:
```
cd /tmp/
curl -L -o assume-role https://github.com/roadtrippers/assume-role/releases/download/1.0.0/assume-role-Darwin
echo 'dff2c8219d8f1ccf4574a5537bd04bf1c9f70f032d243e411b9f3ba724deead4  ./assume-role' > ./assume-role.sha256
shasum -c assume-role.sha256 || echo "DO NOT PROCEED. SHASUM DID NOT MATCH"
mv ./assume-role /usr/local/bin
chmod +x /usr/local/bin/assume-role
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

Useful bash profile setup:
```
function assume-role(){
    unset AWS_ACCESS_KEY_ID
    unset AWS_SECRET_ACCESS_KEY
    unset AWS_SESSION_TOKEN
    unset AWS_SECURITY_TOKEN
    unset ASSUMED_ROLE
    eval $(assume-role $@)
}

print_assumed_role(){
    if test ! -z "${ASSUMED_ROLE}"; then
	echo -n '['
	echo -n $ASSUMED_ROLE


	if test ! -z "${AWS_SESSION_EXPIRATION}"; then
	    echo -n ", ${AWS_SESSION_EXPIRATION}"
	fi

	echo -n '] '
    fi
}

if [ ! -z "${PROMPT_COMMAND}" ]; then
    PROMPT_COMMAND="${PROMPT_COMMAND} ;"
fi

export PROMPT_COMMAND="${PROMPT_COMMAND} print_assumed_role"
```

This will make bash prompt to show assumed role configuration
and its expiration time, as shown below:

```
[vpc-thlprod, Wed 17:41] YOUR_REGULAR_BASH_PROMPT
```
