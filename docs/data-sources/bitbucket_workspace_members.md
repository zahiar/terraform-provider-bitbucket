# Data Source: bitbucket_workspace_members
Use this data source to get a list of members belonging to a workspace, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_workspace" "example" {
  id = "example-slug"
}

data "bitbucket_workspace_members" "example" {
  workspace = data.bitbucket_workspace.id
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace the user belongs to.

## Attribute Reference
In addition to the arguments above, the following attributes are exported:
* `members` - A list of Member information, of which each entry in the list contains:
    * `id` - The Member's UUID.
    * `nickname` - The Member's nickname.
