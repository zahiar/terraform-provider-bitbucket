package bitbucket

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccBitbucketPipelineVariableResource_basic(t *testing.T) {
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
					}`, workspaceSlug, projectName, projectKey, repoName, pipelineVariableName, pipelineVariableValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_pipeline_variable.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_pipeline_variable.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_pipeline_variable.testacc", "key", pipelineVariableName),
					resource.TestCheckResourceAttr("bitbucket_pipeline_variable.testacc", "value", pipelineVariableValue),
					resource.TestCheckResourceAttr("bitbucket_pipeline_variable.testacc", "secured", "true"),

					resource.TestCheckResourceAttrSet("bitbucket_pipeline_variable.testacc", "id"),
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
					}`, workspaceSlug, projectName, projectKey, repoName, pipelineVariableName, pipelineVariableValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_pipeline_variable.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_pipeline_variable.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_pipeline_variable.testacc", "key", pipelineVariableName),
					resource.TestCheckResourceAttr("bitbucket_pipeline_variable.testacc", "value", pipelineVariableValue),
					resource.TestCheckResourceAttr("bitbucket_pipeline_variable.testacc", "secured", "false"),

					resource.TestCheckResourceAttrSet("bitbucket_pipeline_variable.testacc", "id"),
				),
			},
			{
				ResourceName:      "bitbucket_pipeline_variable.testacc",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					resources := state.Modules[0].Resources
					pipelineVariableResourceAttr := resources["bitbucket_pipeline_variable.testacc"].Primary.Attributes
					return fmt.Sprintf("%s/%s/%s", workspaceSlug, repoName, pipelineVariableResourceAttr["id"]), nil
				},
			},
		},
	})
}

func TestValidateRepositoryVariableName(t *testing.T) {
	invalidNames := []string{"4asdasd", "$£$%$", "AS@£$@£$@", "as-adasd"}
	for _, name := range invalidNames {
		validator := validateRepositoryVariableName(name, nil)
		assert.True(t, validator.HasError())
	}

	validNames := []string{"adasd", "asds2342432", "asdadsa_asdasd23424242", "Asdfdfsdf", "AsfdsdfSDFDFSf"}
	for _, name := range validNames {
		validator := validateRepositoryVariableName(name, nil)
		assert.False(t, validator.HasError())
	}
}
