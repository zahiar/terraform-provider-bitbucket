# Resource: bitbucket_group_permission
Manage a group permission for a repository within Bitbucket.

## Example Usage
```hcl
resource "bitbucket_group_permission" "example" {
  workspace  = "{workspace-uuid}"
  repository = "example-repository"
  group      = "example-group"
  permission = "read"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The slug of the repository.
* `group` - (Required) The slug of the group.
* `permission` - (Required) The permission this group will have. Is one of 'read', 'write', or 'admin'.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the group permission.

## Import
Bitbucket group permission can be imported with a combination of its workspace UUID, repository slug & group slug.

### Example using workspace UUID, repository slug & group slug
```sh
$ terraform import bitbucket_group_permission.example "{123ab4cd-5678-9e01-f234-5678g9h01i2j}/example-repo/example-group"
```
