# Data Source: bitbucket_project
Use this data source to get the project resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_project" "example" {
  workspace = "workspace-slug"
  key       = "EXAMPLE"
}
```
```hcl
data "bitbucket_project" "example" {
  workspace = "{workspace-uuid}"
  key       = "EXAMPLE"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace the project belongs to.
* `key` - (Required) The key of the project.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The UUID of the project.
* `name` - A human-readable name of the project.
* `description` - A description of the project.
* `is_private` - A boolean to state if the project is private or not.
