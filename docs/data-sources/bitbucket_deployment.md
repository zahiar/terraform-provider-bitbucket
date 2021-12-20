# Data Source: bitbucket_deployment
Use this data source to get the deployment resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_deployment" "example" {
  id          = "{deployment-id}"
  workspace   = "workspace-slug"
  repository  = "example-repo"
}
```
```hcl
data "bitbucket_deployment" "example" {
  id          = "{deployment-id}"
  workspace   = "{workspace-uuid}"
  repository  = "example-repo"
}
```

## Argument Reference
The following arguments are supported:
* `id` - (Required) The ID of the deployment.
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `name` - The name of the deployment environment.
* `environment` - The environment of the deployment (will be one of 'Test', 'Staging', or 'Production').
