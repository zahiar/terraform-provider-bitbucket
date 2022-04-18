# Data Source: bitbucket_user_workspace

Note: this is an **experimental** data source - please raise tickets if you encounter issues.


Use this data source to get the user resource, you can then reference its attributes without having to hardcode them.

This data source is very similar to `bitbucket_user` - only difference being, this allows you to find users in a given
workspace by their Bitbucket nickname. Be mindful on very large teams, this resource may need to do multiple requests
to Bitbucket's API in order to find them.

## Example Usage
```hcl
data "bitbucket_user_workspace" "example" {
  workspace = "{workspace-uuid}"
  nickname  = "EXAMPLE"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace the user belongs to.
* `nickname` - (Required) The User's nickname.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The User's UUID.
* `display_name` - The User's display name.
* `account_id` - The User's account ID.
* `account_status` - The User's account status. Will be one of active, inactive or closed.
