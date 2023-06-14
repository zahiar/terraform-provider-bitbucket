# Resource: bitbucket_group_member
Manage a user's membership in a group within Bitbucket.

## Example Usage
```hcl
resource "bitbucket_group_member" "example" {
  workspace = "{workspace-uuid}"
  group     = "example-group"
  user      = "{user-uuid}"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The UUID (including the enclosing `{}`) of the workspace.
* `group` - (Required) The slug of the group.
* `user` - (Required) The UUID (including the enclosing `{}`) of the user.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the group membership.

## Import
Bitbucket group members can be imported with a combination of its workspace UUID, group slug & user UUID.

### Example using workspace UUID, group slug & user UUID
```sh
$ terraform import bitbucket_group_member.example "{123ab4cd-5678-9e01-f234-5678g9h01i2j}/example-group/{123ab4cd-5678-9e01-f234-5678g9h01i2j}"
```
