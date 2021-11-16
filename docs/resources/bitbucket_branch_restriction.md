# Resource: bitbucket_branch_restriction
Manage a branch restriction for a repository within Bitbucket.

## Example Usage
```hcl
resource "bitbucket_branch_restriction" "example" {
  workspace  = "workspace-slug"
  repository = "example-repo"
  pattern    = "master"
  kind       = "require_tasks_to_be_completed"
}
```
```hcl
resource "bitbucket_branch_restriction" "min-two-approvals-for-merge" {
  workspace  = "workspace-slug"
  repository = "example-repo"
  pattern    = "master"
  kind       = "require_approvals_to_merge"
  value      = 2
}
```
```hcl
resource "bitbucket_branch_restriction" "restrict-pushes-to-user" {
  workspace  = "workspace-slug"
  repository = "example-repo"
  pattern    = "master"
  kind       = "push"
  users      = ["some-user"]
}
```
```hcl
resource "bitbucket_branch_restriction" "restrict-pushes-to-group" {
  workspace  = "workspace-slug"
  repository = "example-repo"
  pattern    = "master"
  kind       = "push"
  groups      = ["some-user"]
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).
* `pattern` - The pattern to match against branches this restriction will apply to.
* `kind` - The type of restriction to apply - see list below.
* `value` - A configurable value used by the following restrictions: `require_passing_builds_to_merge` uses it to define the number of minimum number of passing builds, `require_approvals_to_merge` uses it to define the minimum number of approvals before the PR can be merged, `require_default_reviewer_approvals_to_merge` uses it to define the minimum number of approvals from default reviewers before the PR can be merged.
* `users` - A list of users (usernames or user's UUID) that are exempt from this branch restriction. Can only be set if restriction type (`kind`) is set to `push` or `restrict_merges`.
* `groups` - A list of groups (group names or group UUID) that are exempt from this branch restriction. Can only be set if restriction type (`kind`) is set to `push` or `restrict_merges`.

<details>
  <summary>Click to view list of supported `kind` values.</summary>

  * `require_tasks_to_be_completed`
  * `allow_auto_merge_when_builds_pass`
  * `require_passing_builds_to_merge`
  * `force`
  * `require_all_dependencies_merged`
  * `require_commits_behind`
  * `restrict_merges`
  * `enforce_merge_checks`
  * `reset_pullrequest_changes_requested_on_change`
  * `require_no_changes_requested`
  * `smart_reset_pullrequest_approvals`
  * `push`
  * `require_approvals_to_merge`
  * `require_default_reviewer_approvals_to_merge`
  * `reset_pullrequest_approvals_on_change`
  * `delete`
</details>

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the branch restriction.

## Import
Bitbucket branch restriction's can be imported with a combination of its workspace slug/UUID, repository name & branch restriction ID.

**_Note: users & groups will not be imported!_**

### Example using workspace slug, repository name & branch restriction ID
```sh
$ terraform import bitbucket_branch_restriction.example "workspace-slug/example-repo/1234"
```

### Example using workspace UUID, repository name & branch restriction ID
```sh
$ terraform import bitbucket_branch_restriction.example "{123ab4cd-5678-9e01-f234-5678g9h01i2j}/example-repo/1234"
```
