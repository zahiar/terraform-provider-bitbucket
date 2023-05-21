package bitbucket

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketPipelineVariableDataSource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	pipelineVariableName := "tf_acc_test" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	pipelineVariableValue := "tf-acc-test" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

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

					resource "bitbucket_pipeline_variable" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  key        = "%s"
					  value      = "%s"
					  secured    = true
					}

					data "bitbucket_pipeline_variable" "testacc" {
					  id         = bitbucket_pipeline_variable.testacc.id
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					}`, workspaceSlug, projectName, projectKey, repoName, pipelineVariableName, pipelineVariableValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_pipeline_variable.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("data.bitbucket_pipeline_variable.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("data.bitbucket_pipeline_variable.testacc", "key", pipelineVariableName),
					resource.TestCheckResourceAttr("data.bitbucket_pipeline_variable.testacc", "value", ""),
					resource.TestCheckResourceAttr("data.bitbucket_pipeline_variable.testacc", "secured", "true"),
					resource.TestCheckResourceAttrSet("data.bitbucket_pipeline_variable.testacc", "id"),
				),
			},
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

					resource "bitbucket_pipeline_variable" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  key        = "%s"
					  value      = "%s"
					  secured    = false
					}

					data "bitbucket_pipeline_variable" "testacc" {
					  id         = bitbucket_pipeline_variable.testacc.id
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					}`, workspaceSlug, projectName, projectKey, repoName, pipelineVariableName, pipelineVariableValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_pipeline_variable.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("data.bitbucket_pipeline_variable.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("data.bitbucket_pipeline_variable.testacc", "key", pipelineVariableName),
					resource.TestCheckResourceAttr("data.bitbucket_pipeline_variable.testacc", "value", pipelineVariableValue),
					resource.TestCheckResourceAttr("data.bitbucket_pipeline_variable.testacc", "secured", "false"),
					resource.TestCheckResourceAttrSet("data.bitbucket_pipeline_variable.testacc", "id"),
				),
			},
		},
	})
}
