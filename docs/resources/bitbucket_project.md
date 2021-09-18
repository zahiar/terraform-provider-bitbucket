# Resource: bitbucket_project
Manage a project within Bitbucket.

## Example Usage
```hcl
resource "bitbucket_project" "example" {
  workspace   = "workspace-slug"
  name        = "Example Project"
  key         = "EXP"
  description = "An example project"
  is_private  = true
}
```
```hcl
resource "bitbucket_project" "example" {
  workspace   = "{workspace-uuid}"
  name        = "Example Project"
  key         = "EXP"
  description = "An example project"
  is_private  = true
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace the project belongs to.
* `name` - (Required) The name of the project.
* `key` - (Required) The key of the project.
* `description` - (Optional) The description of the project. Defaults to empty string.
* `is_private` - (Optional) A boolean to state if the project is private or not. Defaults to `true`.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The UUID of the project.

## Import
Bitbucket project can be imported with a combination of its workspace slug/UUID & key.

### Example using workspace slug & project key
```sh
$ terraform import bitbucket_project.example "workspace-slug/EXP"
```

### Example using workspace UUID & project key
```sh
$ terraform import bitbucket_project.example "{123ab4cd-5678-9e01-f234-5678g9h01i2j}/EXP"
```
