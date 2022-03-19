# Resource: bitbucket_default_reviewer
Manage a default reviewer for a repository within Bitbucket.

## Example Usage
```hcl
resource "bitbucket_default_reviewer" "example" {
  workspace  = "workspace-slug"
  repository = "example-repo"
  user       = "{user-uuid}"
}
```
```hcl
resource "bitbucket_default_reviewer" "example" {
  workspace  = "{workspace-uuid}"
  repository = "example-repo"
  user       = "{user-uuid}"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores, hyphens and periods).
* `user` - (Required) The user's UUID.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the default reviewer.

## Import
Bitbucket default reviewers can be imported with a combination of its workspace slug/UUID, repository name & user's UUID.

### Example using workspace slug, repository name & user's UUID
```sh
$ terraform import bitbucket_default_reviewer.example "workspace-slug/example-repo/{123ab4cd-5678-9e01-f234-5678g9h01i2j}"
```

### Example using workspace UUID, repository name & user's UUID
```sh
$ terraform import bitbucket_default_reviewer.example "{123ab4cd-5678-9e01-f234-5678g9h01i2j}/example-repo/{123ab4cd-5678-9e01-f234-5678g9h01i2j}"
```
