package bitbucket

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBitbucketPipelineKeyPairResource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	pipelineKeyPairPublicSSHKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAK/b1cHHDr/TEV1JGQl+WjCwStKG6Bhrv0rFpEsYlyTBm1fzN0VOJJYn4ZOPCPJwqse6fGbXntEs+BbXiptR+++HycVgl65TMR0b5ul5AgwrVdZdT7qjCOCgaSV74/9xlHDK8oqgGnfA7ZoBBU+qpVyaloSjBdJfLtPY/xqj4yHnXKYzrtn/uFc4Kp9Tb7PUg9Io3qohSTGJGVHnsVblq/rToJG7L5xIo0OxK0SJSQ5vuId93ZuFZrCNMXj8JDHZeSEtjJzpRCBEXHxpOPhAcbm4MzULgkFHhAVgp4JbkrT99/wpvZ7r9AdkTg7HGqL3rlaDrEcWfL7Lu6TnhBdq5"
	pipelineKeyPairPrivateSSHKeyEscaped := "-----BEGIN OPENSSH PRIVATE KEY-----\\n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\\nAgMEBQ==\\n-----END OPENSSH PRIVATE KEY-----"
	pipelineKeyPairPrivateSSHKey := "-----BEGIN OPENSSH PRIVATE KEY-----\n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\nAgMEBQ==\n-----END OPENSSH PRIVATE KEY-----"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}
	
					resource "bitbucket_project" "testacc" {
					  workspace = data.bitbucket_workspace.testacc.id
					  name      = "%s"
					  key       = "%s"
					}
	
					resource "bitbucket_repository" "testacc" {
					  workspace        = data.bitbucket_workspace.testacc.id
					  project_key      = bitbucket_project.testacc.key
					  name             = "%s"
					  enable_pipelines = true
					}
	
					resource "bitbucket_pipeline_key_pair" "testacc" {
					  workspace   = data.bitbucket_workspace.testacc.id
					  repository  = bitbucket_repository.testacc.name
					  public_key  = "%s"
					  private_key = "%s"
					}`, workspaceSlug, projectName, projectKey, repoName, pipelineKeyPairPublicSSHKey, pipelineKeyPairPrivateSSHKeyEscaped),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_pipeline_key_pair.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_pipeline_key_pair.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_pipeline_key_pair.testacc", "public_key", pipelineKeyPairPublicSSHKey),
					resource.TestCheckResourceAttr("bitbucket_pipeline_key_pair.testacc", "private_key", pipelineKeyPairPrivateSSHKey),

					resource.TestCheckResourceAttrSet("bitbucket_pipeline_key_pair.testacc", "id"),
				),
			},
		},
	})
}

func TestGeneratePipelineKeyPairId(t *testing.T) {
	expected := "{my-workspace-uuid}-my-repo"
	result := generatePipelineKeyPairId("{my-workspace-uuid}", "my-repo")
	assert.Equal(t, expected, result)
}
