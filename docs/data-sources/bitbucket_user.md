# Data Source: bitbucket_user
Use this data source to get the user resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_user" "example" {
  id = "{user-uuid}"
}
```

## Argument Reference
The following arguments are supported:
* `id` - (Required) The User's UUID.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `nickname` - The User's nickname.
* `display_name` - The User's display name.
* `account_id` - The User's account ID.
* `account_status` - The User's account status. Will be one of active, inactive or closed.
