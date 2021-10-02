# Data Source: bitbucket_default_reviewer
Use this data source to get the default reviewer resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_default_reviewer" "example" {
  workspace  = "workspace-slug"
  repository = "example-repo"
  user       = "{user-uuid}"
}
```
```hcl
data "bitbucket_default_reviewer" "example" {
  workspace  = "{workspace-uuid}"
  repository = "example-repo"
  user       = "{user-uuid}"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).
* `user` - (Required) The user's UUID.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the default reviewer.
