# Resource: bitbucket_deployment
Manage a deployment environment for a repository within Bitbucket.

## Example Usage
```hcl
resource "bitbucket_deployment" "example" {
  workspace   = "workspace-slug"
  repository  = "example-repo"
  name        = "Example environment"
  environment = "Staging"
}
```
```hcl
resource "bitbucket_deployment" "example" {
  workspace   = "{workspace-uuid}"
  repository  = "example-repo"
  name        = "Example environment"
  environment = "Staging"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).
* `name` - (Required) The name of the deployment environment.
* `environment` - (Required)The environment of the deployment (must be either 'Test', 'Staging' or 'Production').

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the deployment.

## Import
Bitbucket deployment environment's can be imported with a combination of its workspace slug/UUID, repository name & deployment environment ID.

### Example using workspace slug, repository name & deployment environment ID
```sh
$ terraform import bitbucket_deployment.example "workspace-slug/example-repo/1234"
```

### Example using workspace UUID, repository name & deployment environment ID
```sh
$ terraform import bitbucket_deployment.example "{123ab4cd-5678-9e01-f234-5678g9h01i2j}/example-repo/1234"
```
