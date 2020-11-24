# Device42 Terraform provider for passwords
---
## Supports:
* Retrieval of passwords from Device42 by reference to their numeric `id`; this exposes the following attributes:
  * `id` (which you already know, since it's used as the selector)
  * `username`
  * `password`
  * `label`
* That's basically it, right now

## Requirements:
* A Device42 host (configured in the `provider {}` block, or from the `D42_HOSTNAME` environment variable; defaults to `swaggerdemo.device42.com`)
* A username with access (can be configured in the `provider {}` block, or from the `D42_USERNAME` environment variable)
* A password with access (can be configured in the `provider {}` block, or from the `D42_PASSWORD` environment variable)
* Terraform (tested with `0.12.29`, `0.13.5`, and `0.14.0`)

## Proof of concept (macOS-specific):

```sh
# build the binary in the test directory for 0.12.x
go build -o tf_test/0.12/terraform-provider-device42

# Terraform 0.13.x changes the filesystem layout and config schema for custom providers
mkdir -p ~/.terraform.d/plugins/terraform.example.com/example/device42/0.0.1/darwin_amd64
go build -o ~/.terraform.d/plugins/terraform.example.com/example/device42/0.0.1/darwin_amd64/terraform-provider-device42_v0.0.1

# set the hostname for your Device42 instance
export D42_HOSTNAME=device42.company.com

# set the username in an environment variable and export it for child processes
export D42_USERNAME=you

# read your (non-echoed) password from the shell prompt and export it for child processes
read -s D42_PASSWORD
export D42_PASSWORD

# use main.tf to demo retrieving a value
cd tf_test/${version}
< edit test.tf to have a valid ID >
terraform destroy -auto-approve && terraform apply -auto-approve
```

## Testing (or lack thereof...)

This is currently very difficult/borderline impossible because I no longer have access to a Device42 instance having now left the company where we used Device42 for secret storage. I have a trial version of `16.18.00.1601912441` which this has been tested against (until the trial license runs out), and it seems to work OK, with some quirks:

* Device42 doesn't always return an actual error message, sometimes it's just `msg="", code=0`, which you would _think_ means success. It doesn't.
* Device42 doesn't always return JSON (even if you specifically request it in headers), sometimes it's just a quoted string to say you don't have permissions. It's not clean, but it should be handled and spit out a sensible error message.

### Testing against the public demo

There is a public demo instance ([`https://swaggerdemo.device42.com`](https://swaggerdemo.device42.com) which ties in with the API reference at [`https://api.device42.com`](https://api.device42.com)) however as expected from a demo site, the data is reset periodically so there are no persistent passwords stored which could be referred to.

This would have been my motivation to also write `resource_password.go` in order to create a password via a POST request, but the instance is not currently configured to accept password storage:

```sh
~ » curl -X POST -H "Authorization: Basic AA==" -H 'Content-Type: application/x-www-form-urlencoded' -H 'Accept: application/json' --data 'username=test&password=testPW&view_edit_users=guest,api_user' https://swaggerdemo.device42.com/api/1.0/passwords/

{"msg": "Please enter the passphrase first. Go to Tools > Settings > Password Security", "code": 2}
```

So there's no persistent data that could be used for testing, and no way to add a new password as part of the tests. I've sent an email to the Device42 `support@` address to see if they could add an encryption passphrase to allow this demo instance to be testable, but won't hold my breath

## To do
* watch https://support.device42.com/hc/en-us/community/posts/360009940834-API-tokens-instead-of-user-pass in case D42 ever implements token-based API access for systems like Jenkins instead of just username/password auth
* support a "password command" env.var (i.e. `pass test/device42` to use passwordstore.org) rather than a password env.var
* maybe support inserting secrets
  * an example workflow might be to generate a random value in Terraform using the `random` provider, then store it in Device42 in the same session?
* work out how this could be used as a custom provider since I don't want to put it in the actual Terraform registry...
  * This may partly be addressed by pushing to an internal artifact repository and using the terraform 0.13 `required_providers` blocks
