# Resource: bitbucket_repository
Manage a repository within Bitbucket.

## Example Usage
```hcl
resource "bitbucket_repository" "example" {
  workspace        = "workspace-slug"
  name             = "example-repo"
  project_key      = "EXP"
  description      = "An example repository"
  is_private       = true
  has_wiki         = true
  fork_policy      = "no_forks"
  enable_pipelines = false
  language         = "go"
}
```
```hcl
resource "bitbucket_repository" "example" {
  workspace        = "{workspace-uuid}"
  name             = "example-repo"
  project_key      = "EXP"
  description      = "An example repository"
  is_private       = true
  has_wiki         = true
  fork_policy      = "no_forks"
  enable_pipelines = false
  language         = "go"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace this repository belongs to.
* `name` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores, hyphens and periods).
* `project_key` - (Required) The key of the project this repository belongs to.
* `description` - (Optional) The description of the repository. Defaults to empty string.
* `is_private` - (Optional) A boolean to state if the repository is private or not. Defaults to `true`.
* `has_wiki` - (Optional) A boolean to state if the repository includes a wiki or not. Defaults to `false`.
* `fork_policy` - (Optional) The name of the fork policy to apply to this repository. Defaults to `no_forks`. Only applies if `is_private` is set to `true`.
* `enable_pipelines` - (Optional) A boolean to state if pipelines have been enabled for this repository. Defaults to `false`.
* `language` - (Optional) The programming language of the repository. Defaults to empty string.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The UUID of the repository.

## Import
Bitbucket repository can be imported with a combination of its workspace slug/UUID & name.

### Example using workspace slug & repository name
```sh
$ terraform import bitbucket_repository.example "workspace-slug/example-repo"
```

### Example using workspace UUID & repository name
```sh
$ terraform import bitbucket_repository.example "{123ab4cd-5678-9e01-f234-5678g9h01i2j}/example-repo"
```
