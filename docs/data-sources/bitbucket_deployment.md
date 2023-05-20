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
```hcl
data "bitbucket_deployment" "example" {
  name        = "deployment-name"
  workspace   = "workspace-slug"
  repository  = "example-repo"
}
```
```hcl
data "bitbucket_deployment" "example" {
  name        = "deployment-name"
  workspace   = "{workspace-uuid}"
  repository  = "example-repo"
}
```

## Argument Reference
The following arguments are supported:
* `id` - (Optional) The ID of the deployment.
* `name` - (Optional) The name of the deployment.
~> **NOTE:** Either `id` or `name` must be set. If both are set, `name` will be ignored.
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores, hyphens and periods).


## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the deployment environment.
* `name` - The name of the deployment environment.
* `environment` - The environment of the deployment (will be one of 'Test', 'Staging', or 'Production').
