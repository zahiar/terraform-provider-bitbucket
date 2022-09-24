# Data Source: bitbucket_group
Use this data source to get the group resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_group" "example" {
  workspace = "{workspace-uuid}"
  slug      = "example-group"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The UUID (including the enclosing `{}`) of the workspace the group belongs to.
* `slug` - (Required) The slug of the group.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the group.
* `name` - A human-readable name of the group.
* `permission` - The global permission this group will have over all repositories. Is one of 'none', 'read', 'write', 'admin'.
