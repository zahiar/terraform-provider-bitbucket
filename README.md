# Terraform Provider: Bitbucket Cloud
This is a Terraform provider for managing resources within a Bitbucket Cloud account.

## Getting Started
As this provider is published to the public [Terraform Registry](https://registry.terraform.io/providers/zahiar/bitbucket),
you can install it like so (for Terraform 0.14+):
```hcl
provider "bitbucket" {
  username = "myUsername"
  password = "myPassword"
}

terraform {
  required_providers {
    bitbucket = {
      source  = "zahiar/bitbucket"
    }
  }
}
```

For more detailed instructions and documentation on the resources and data sources supported, please go to
[Terraform Registry](https://registry.terraform.io/providers/zahiar/bitbucket/latest/docs).

## Maintenance
This provider is maintained during free time, so if you are interested in helping to develop this further, you
are more than welcome to submit a pull request or raise a ticket if you'd prefer.

## Development

### Requirements
If you do wish to help develop this, you will need the following installed:
* [Go](http://www.golang.org) (see `go.mod` file for the correct version to install)
* [Go Linter](https://formulae.brew.sh/formula/golangci-lint)
* [GOPATH](http://golang.org/doc/code.html#GOPATH) (is correctly setup)
* [Terraform](https://www.terraform.io/downloads.html) (0.14+)

### Building
Simply run `make build`, and it will compile and create a binary, as well as print-out instructions
on how to configure Terraform to use this locally built provider.
```shell
$ make build
```

### Testing

#### Unit Tests 
```shell
$ make test
```

### Acceptance Tests
This will require you to specify the following environment variables, as these tests will provision actual resources in
your account, and it will tear them down afterwards to ensure it leaves your account clean.

You will also require a UUID of another account that is a member of your workspace in order for the `bitbucket_user_permission` 
tests to run, as Bitbucket's API will reject the account owner's UUID.

* `BITBUCKET_USERNAME` - Username of the account to run the tests against
* `BITBUCKET_PASSWORD` - Password of the account to run the tests against
* `BITBUCKET_MEMBER_ACCOUNT_UUID` - Account UUID of the member who is part of your account

**NOTE**: if a test fails, it may leave dangling resources in your account so please bear this in mind.

If you have two-factor authentication enabled, then be sure to set up an [app password](https://support.atlassian.com/bitbucket-cloud/docs/app-passwords/) and use that instead.
```shell
$ BITBUCKET_USERNAME=myUsername BITBUCKET_PASSWORD=myPassword BITBUCKET_MEMBER_ACCOUNT_UUID=myMemberUUID make testacc
```

### Documentation
Every data source or resource added must have an accompanying docs page (see `docs` directory for examples).

Docs are written using Markdown, and you can use [this page](https://registry.terraform.io/tools/doc-preview) to preview what your docs will look like when rendered.
