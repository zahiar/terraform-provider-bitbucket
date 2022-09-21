# Data Source: bitbucket_repository
Use this data source to get the repository resource, you can then reference its attributes without having to hardcode them.

## Example Usage
```hcl
data "bitbucket_repository" "example" {
  workspace   = "workspace-slug"
  name        = "example-repo"
}
```
```hcl
data "bitbucket_repository" "example" {
  workspace   = "{workspace-uuid}"
  name        = "example-repo"
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace this repository belongs to.
* `name` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores, hyphens and periods).

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The UUID of the repository.
* `project_key` - The key of the project this repository belongs to.
* `description` - The description of the repository.
* `is_private` - A boolean to state if the repository is private or not.
* `has_wiki` - A boolean to state if the repository includes a wiki or not.
* `fork_policy` - The name of the fork policy set on the repository.
* `enable_pipelines` - A boolean to state if pipelines have been enabled for this repository.
