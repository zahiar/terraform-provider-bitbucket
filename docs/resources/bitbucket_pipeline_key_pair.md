# Resource: bitbucket_pipeline_key_pair
Manage a pipeline key pair for a repository within Bitbucket.

**Note**: only a single pipeline key pair can exist per repository.

## Example Usage
```hcl
resource "bitbucket_pipeline_key_pair" "example" {
  workspace   = "workspace-slug"
  repository  = "example-repo"
  public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAK..."
  private_key = "-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEZZZZABG5vbmUZZZZEbm9uZQZZZZZZZZABZZABlwZZZZdzc2gtcn\nNhZZZZAwEZZQZZAY..."
}
```
```hcl
resource "bitbucket_pipeline_key_pair" "example" {
  workspace  = "{workspace-uuid}"
  repository = "example-repo"
  public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAK..."
  private_key = "-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEZZZZABG5vbmUZZZZEbm9uZQZZZZZZZZABZZABlwZZZZdzc2gtcn\nNhZZZZAwEZZQZZAY..."
}
```

## Argument Reference
The following arguments are supported:
* `workspace` - (Required) The slug or UUID (including the enclosing `{}`) of the workspace.
* `repository` - (Required) The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores, hyphens and periods).
* `public_key` - (Required) The public SSH key part of this pipeline key pair.
* `private_key` - (Required) The private SSH key part of this pipeline key pair.

## Attribute Reference
In addition to the arguments above, the following additional attributes are exported:
* `id` - The ID of the pipeline key pair.
