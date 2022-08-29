# Data Source: bitbucket_user_permission
Use this data source to get the user permission resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_user_permission" "example" {
  workspace  = "{workspace-uuid}"
  repository = "example-repository"
  user       = "{user-uuid}"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The slug of the repository.
* `user` - (Required) The UUID (including the enclosing `{}`) of the user.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the user permission.
* `permission` - The permission this user will have. Is one of 'read', 'write', or 'admin'.
