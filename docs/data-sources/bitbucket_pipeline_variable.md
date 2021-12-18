# Data Source: bitbucket_pipeline_variable
Use this data source to get the pipeline variable resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_pipeline_variable" "example" {
  id          = "{pipeline-variable-id}"
  workspace   = "workspace-slug"
  repository  = "example-repo"
}
```
```hcl
data "bitbucket_pipeline_variable" "example" {
  id          = "{pipeline-variable-id}"
  workspace   = "{workspace-uuid}"
  repository  = "example-repo"
}
```

## Argument Reference
The following arguments are supported:
* `id` - (Required) The ID of the pipeline variable.
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `key` - The name of the variable.
* `value` - The value of the variable (note: if this variable is marked 'secured', this attribute will be blank).
* `secured` - Whether this variable is considered secure/sensitive. If true, then it's value will not be exposed in any logs or API requests.
