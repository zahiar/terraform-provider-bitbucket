# Resource: bitbucket_group
Manage a group within Bitbucket.

## Example Usage
```hcl
resource "bitbucket_group" "example" {
  workspace  = "{workspace-uuid}"
  name       = "Example Group"
  auto_add   = false
  permission = "read"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The UUID (including the enclosing `{}`) of the workspace the group belongs to.
* `name` - (Required) A human-readable name of the group.
* `auto_add` - (Optional) A boolean to state whether this group is auto-added to all future repositories.
* `permission` - (Optional) The permission this group will have over repositories. Is one of 'read', 'write', or 'admin'.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the group.
* `slug` - The slug of the group.

## Import
Bitbucket group can be imported with a combination of its workspace UUID & group slug.

### Example using workspace UUID & group slug
```sh
$ terraform import bitbucket_group.example "{123ab4cd-5678-9e01-f234-5678g9h01i2j}/example-group"
```
