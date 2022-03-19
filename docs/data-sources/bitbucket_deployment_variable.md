# Data Source: bitbucket_deployment_variable
Use this data source to get the deployment variable resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_deployment_variable" "example" {
  workspace  = "workspace-slug"
  repository = "example-repo"
  deployment = "{deployment-environment-id}"
  key        = "some_variable_name"
}
```
```hcl
data "bitbucket_deployment_variable" "example" {
  workspace  = "{workspace-uuid}"
  repository = "example-repo"
  deployment = "{deployment-environment-id}"
  key        = "some_variable_name"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores, hyphens and periods).
* `deployment` - (Required) The UUID (including the enclosing `{}`) of the deployment environment.
* `key` - (Required) The name of the variable (must consist of only ASCII letters, numbers, underscores & not begin with a number).

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the deployment variable.
* `value` - The value of the variable (note: if this variable is marked 'secured', this attribute will be blank).
* `secured` - Whether this variable is considered secure/sensitive. If true, then it's value will not be exposed in any logs or API requests.
