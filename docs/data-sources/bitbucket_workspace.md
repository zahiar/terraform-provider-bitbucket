# Data Source: bitbucket_workspace
Use this data source to get the workspace resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_workspace" "example" {
  id = "example-slug"
}
```

## Argument Reference
The following arguments are supported:
* `id` - (Required) The slug of the workspace.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `type` - The type of the workspace.
* `name` - The name of the workspace.
* `uuid` - The UUID of the workspace.
* `is_private` - A boolean to state if the project is private or not.
