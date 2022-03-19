# Data Source: bitbucket_branch_restriction
Use this data source to get the branch restriction resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_branch_restriction" "example" {
  id          = "{branch-restriction-id}"
  workspace   = "workspace-slug"
  repository  = "example-repo"
}
```
```hcl
data "bitbucket_branch_restriction" "example" {
  id          = "{branch-restriction-id}"
  workspace   = "{workspace-uuid}"
  repository  = "example-repo"
}
```

## Argument Reference
The following arguments are supported:
* `id` - (Required) The ID of the branch restriction key.
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores, hyphens and periods).

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `pattern` - The pattern to match against branches this restriction will apply to.
* `kind` - The type of restriction to apply.
* `value` - A configurable value used by the following restrictions: `require_passing_builds_to_merge` uses it to define the number of minimum number of passing builds, `require_approvals_to_merge` uses it to define the minimum number of approvals before the PR can be merged, `require_default_reviewer_approvals_to_merge` uses it to define the minimum number of approvals from default reviewers before the PR can be merged.
