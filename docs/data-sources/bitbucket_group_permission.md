# Data Source: bitbucket_group_permission
Use this data source to get the group permission resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_group_permission" "example" {
  workspace  = "{workspace-uuid}"
  repository = "example-repository"
  group      = "example-group"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The slug of the repository.
* `group` - (Required) The slug of the group.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the group permission.
* `permission` - The permission this group will have. Is one of 'read', 'write', or 'admin'.
