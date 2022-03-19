# Data Source: bitbucket_webhook
Use this data source to get the webhook resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_webhook" "example" {
  id          = "{webhook-uuid}"
  workspace   = "workspace-slug"
  repository  = "example-repo"
}
```
```hcl
data "bitbucket_webhook" "example" {
  id          = "{webhook-uuid}"
  workspace   = "{workspace-uuid}"
  repository  = "example-repo"
}
```

## Argument Reference
The following arguments are supported:
* `id` - (Required) The UUID (including the enclosing `{}`) of the webhook.
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace this webhook belongs to.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores, hyphens and periods).

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `name` - The name of the webhook.
* `url` - The url the webhook is configured with.
* `events` - A list of events that will trigger the webhook.
* `is_active` - A boolean to state if the webhook is active or not.
