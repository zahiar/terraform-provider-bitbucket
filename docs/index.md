# Terraform Provider: Bitbucket Cloud
This is a Terraform provider for managing resources within a Bitbucket Cloud account.

In terms of authentication, you must use your Bitbucket username (not your email address) & either your password, or
if you have two-factor authentication enabled, then you must use an app password.
Visit here for more information on [app passwords](https://support.atlassian.com/bitbucket-cloud/docs/app-passwords/).

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
