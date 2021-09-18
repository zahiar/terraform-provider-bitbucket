# Bitbucket Terraform Provider
The Bitbucket Terraform provider is used to interact with the Bitbucket API.

In terms of authentication, you must use your Bitbucket username (not your email address) & either your password, or
if you have two-factor authentication enabled, then you must use an app password.
Go here for more info on app passwords: https://support.atlassian.com/bitbucket-cloud/docs/app-passwords/

## Example Usage
### With Embedded Credentials
```hcl
provider "bitbucket" {
  username = "my-username" 
  password = "my-password"
}

resource "bitbucket_xxx" "example" {
  ...
}
```

### With Environment Variables
Please set the following environment variables:
```shell
BITBUCKET_USERNAME=my-username
BITBUCKET_PASSWORD=my-password
```

```hcl
provider "bitbucket" {}

resource "bitbucket_xxx" "example" {
  ...
}
```
