# Resource: bitbucket_user_permission
Manage a user permission for a repository within Bitbucket.

## Example Usage
```hcl
resource "bitbucket_user_permission" "example" {
  workspace  = "{workspace-uuid}"
  repository = "example-repository"
  user       = "{user-uuid}"
  permission = "read"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The slug of the repository.
* `user` - (Required) The UUID (including the enclosing `{}`) of the user.
* `permission` - (Required) The permission this user will have. Is one of 'read', 'write', or 'admin'.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the user permission.

## Import
Bitbucket user permission can be imported with a combination of its workspace UUID, repository slug & user UUID.

### Example using workspace UUID, repository slug & user UUID
```sh
$ terraform import bitbucket_user_permission.example "{123ab4cd-5678-9e01-f234-5678g9h01i2j}/example-repo/{123ab4cd-5678-9e01-f234-5678g9h01i2j}"
```
