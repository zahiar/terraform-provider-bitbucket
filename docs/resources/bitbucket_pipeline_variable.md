# Resource: bitbucket_pipeline_variable
Manage a pipeline variable for a repository within Bitbucket.

## Example Usage
```hcl
resource "bitbucket_pipeline_variable" "example" {
  workspace  = "workspace-slug"
  repository = "example-repo"
  key        = "some-variable-name"
  value      = "some-variable-value"
  secured    = false
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).
* `key` - (Required) The name of the variable (must consist of only ASCII letters, numbers, underscores & not begin with a number).
* `value` - (Required) The value of the variable.
* `secured` - (Optional) Whether this variable is considered secure/sensitive. If true, then it's value will not be exposed in any logs or API requests. Defaults to `false`.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the pipeline variable.

## Import
Bitbucket pipeline variable's can be imported with a combination of its workspace slug/UUID, repository name & pipeline variable ID.

**_Note: secured values will not be imported!_**

### Example using workspace slug, repository name & pipeline variable ID
```sh
$ terraform import bitbucket_pipeline_variable.example "workspace-slug/example-repo/1234"
```

### Example using workspace UUID, repository name & pipeline variable ID
```sh
$ terraform import bitbucket_pipeline_variable.example "{123ab4cd-5678-9e01-f234-5678g9h01i2j}/example-repo/1234"
```
