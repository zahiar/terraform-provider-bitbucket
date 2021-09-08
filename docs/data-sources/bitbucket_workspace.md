# Data Source: bitbucket_workspace
Use this data source to get the workspace resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_workspace" "example" {
  name = "example"
}
```

## Argument Reference
The following arguments are supported:
* `name` - (Required) The name of the workspace.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The UUID of the workspace.
* `type` - The type of the workspace.
* `slug` - The slug of the workspace.
* `is_private` - A boolean to state if the project is private or not.
