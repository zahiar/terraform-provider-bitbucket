# Resource: bitbucket_deploy_key
Manage a deploy key for a repository within Bitbucket.

## Example Usage
```hcl
resource "bitbucket_deploy_key" "example" {
  workspace  = "workspace-slug"
  repository = "example-repo"
  label      = "Example deploy key"
  key        = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAK..."
}
```
```hcl
resource "bitbucket_deploy_key" "example" {
  workspace  = "{workspace-uuid}"
  repository = "example-repo"
  label      = "Example deploy key"
  key        = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAK..."
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).
* `label` - (Required) The label of the deploy key.
* `key` - (Required) The public SSH key attached to this deploy key.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the deploy key.

## Import
Bitbucket deploy key's can be imported with a combination of its workspace slug/UUID, repository name & deploy key ID.

### Example using workspace slug, repository name & deploy key ID
```sh
$ terraform import bitbucket_deploy_key.example "workspace-slug/example-repo/1234"
```

### Example using workspace UUID, repository name & deploy key ID
```sh
$ terraform import bitbucket_deploy_key.example "{123ab4cd-5678-9e01-f234-5678g9h01i2j}/example-repo/1234"
```
