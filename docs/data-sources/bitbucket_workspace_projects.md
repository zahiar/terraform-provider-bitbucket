# Data Source: bitbucket_workspace_projects
Use this data source to get a list of projects belonging to a workspace, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_workspace" "example" {
  id = "example-slug"
}

data "bitbucket_workspace_projects" "example" {
  workspace = data.bitbucket_workspace.id
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace the user belongs to.

## Attribute Reference
In addition to the arguments above, the following attributes are exported:
* `projects` - A list of Project information, of which each entry in the list contains:
    * `id` - The Project's UUID.
    * `name` - The Project's name.
    * `key` - The Project's key.
