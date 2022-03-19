# Resource: bitbucket_webhook
Manage a webhook for a repository within Bitbucket.

## Example Usage
```hcl
resource "bitbucket_webhook" "example" {
  workspace  = "workspace-slug"
  repository = "example-repo"
  name       = "Example webhook"
  url        = "https://example.webook"
  events     = ["pullrequest:approved", "pullrequest:comment_updated"]
  is_active  = true
}
```
```hcl
resource "bitbucket_webhook" "example" {
  workspace  = "{workspace-uuid}"
  repository = "example-repo"
  name       = "Example webhook"
  url        = "https://example.webook"
  events     = ["pullrequest:approved", "pullrequest:comment_updated"]
  is_active  = true
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace this webhook belongs to.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores, hyphens and periods).
* `name` - (Required) The name of the webhook.
* `url` - (Required) The url to configure the webhook with.
* `events` - (Required) A list of events that will trigger the webhook - see list below.
* `is_active` - (Optional) A boolean to state if the webhook is active or not. Defaults to `false`.


<details>
  <summary>Click to view list of supported events.</summary>

  * `issue:comment_created`
  * `issue:created`
  * `issue:updated`
  * `project:updated`
  * `pullrequest:approved`
  * `pullrequest:changes_request_created`
  * `pullrequest:changes_request_removed`
  * `pullrequest:comment_created`
  * `pullrequest:comment_deleted`
  * `pullrequest:comment_updated`
  * `pullrequest:created`
  * `pullrequest:fulfilled`
  * `pullrequest:rejected`
  * `pullrequest:unapproved`
  * `pullrequest:updated`
  * `repo:commit_comment_created`
  * `repo:commit_status_created`
  * `repo:commit_status_updated`
  * `repo:created`
  * `repo:deleted`
  * `repo:fork`
  * `repo:imported`
  * `repo:push`
  * `repo:transfer`
  * `repo:updated`
</details>

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The UUID of the webhook.

## Import
Bitbucket webhook's can be imported with a combination of its workspace slug/UUID, repository name & webhook UUID.

### Example using workspace slug, repository name & webhook UUID
```sh
$ terraform import bitbucket_webhook.example "workspace-slug/example-repo/{123ab4cd-5678-9e01-f234-5678g9h01i2j}"
```

### Example using workspace UUID, repository name & webhook UUID
```sh
$ terraform import bitbucket_webhook.example "{123ab4cd-5678-9e01-f234-5678g9h01i2j}/example-repo/{123ab4cd-5678-9e01-f234-5678g9h01i2j}"
```
