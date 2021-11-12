# Data Source: bitbucket_deploy_key
Use this data source to get the deploy key resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_deploy_key" "example" {
  id          = "{deploy-key-id}"
  workspace   = "workspace-slug"
  repository  = "example-repo"
}
```
```hcl
data "bitbucket_deploy_key" "example" {
  id          = "{deploy-key-id}"
  workspace   = "{workspace-uuid}"
  repository  = "example-repo"
}
```

## Argument Reference
The following arguments are supported:
* `id` - (Required) The ID of the deploy key.
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `label` - The label of the deploy key.
* `key` - The public SSH key attached to this deploy key.
