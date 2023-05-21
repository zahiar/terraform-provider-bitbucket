package bitbucket

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketDeploymentVariableDataSource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	deploymentVariableNameNotSecure := "tf_acc_test" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	deploymentVariableNameSecure := "tf_acc_test" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	deploymentVariableValue := "tf-acc-test" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

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
					  workspace  = data.bitbucket_workspace.testacc.id
					  name       = "%s"
					  key        = "%s"
					  is_private = true
					}

					resource "bitbucket_repository" "testacc" {
					  workspace        = data.bitbucket_workspace.testacc.id
					  project_key      = bitbucket_project.testacc.key
					  name             = "%s"
					  enable_pipelines = true
					}

					resource "bitbucket_deployment" "testacc" {
					  workspace   = data.bitbucket_workspace.testacc.id
					  repository  = bitbucket_repository.testacc.name
					  name        = "TF ACC Test Deployment"
					  environment = "Test"
					}

					resource "bitbucket_deployment_variable" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  deployment = bitbucket_deployment.testacc.id
					  key        = "%s"
					  value      = "%s"
					  secured    = true
					}

					data "bitbucket_deployment_variable" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  deployment = bitbucket_deployment.testacc.id
					  key        = bitbucket_deployment_variable.testacc.key
					}`, workspaceSlug, projectName, projectKey, repoName, deploymentVariableNameSecure, deploymentVariableValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_deployment_variable.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("data.bitbucket_deployment_variable.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("data.bitbucket_deployment_variable.testacc", "key", deploymentVariableNameSecure),
					resource.TestCheckResourceAttr("data.bitbucket_deployment_variable.testacc", "value", ""),
					resource.TestCheckResourceAttr("data.bitbucket_deployment_variable.testacc", "secured", "true"),

					resource.TestCheckResourceAttrSet("data.bitbucket_deployment_variable.testacc", "id"),
					resource.TestCheckResourceAttrSet("data.bitbucket_deployment_variable.testacc", "deployment"),
				),
			},
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}
	
					resource "bitbucket_project" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  name       = "%s"
					  key        = "%s"
					  is_private = true
					}

					resource "bitbucket_repository" "testacc" {
					  workspace        = data.bitbucket_workspace.testacc.id
					  project_key      = bitbucket_project.testacc.key
					  name             = "%s"
					  enable_pipelines = true
					}

					resource "bitbucket_deployment" "testacc" {
					  workspace   = data.bitbucket_workspace.testacc.id
					  repository  = bitbucket_repository.testacc.name
					  name        = "TF ACC Test Deployment"
					  environment = "Test"
					}

					resource "bitbucket_deployment_variable" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  deployment = bitbucket_deployment.testacc.id
					  key        = "%s"
					  value      = "%s"
					  secured    = false
					}

					data "bitbucket_deployment_variable" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  deployment = bitbucket_deployment.testacc.id
					  key        = bitbucket_deployment_variable.testacc.key
					}`, workspaceSlug, projectName, projectKey, repoName, deploymentVariableNameNotSecure, deploymentVariableValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_deployment_variable.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("data.bitbucket_deployment_variable.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("data.bitbucket_deployment_variable.testacc", "key", deploymentVariableNameNotSecure),
					resource.TestCheckResourceAttr("data.bitbucket_deployment_variable.testacc", "value", deploymentVariableValue),
					resource.TestCheckResourceAttr("data.bitbucket_deployment_variable.testacc", "secured", "false"),

					resource.TestCheckResourceAttrSet("data.bitbucket_deployment_variable.testacc", "id"),
					resource.TestCheckResourceAttrSet("data.bitbucket_deployment_variable.testacc", "deployment"),
				),
			},
		},
	})
}
