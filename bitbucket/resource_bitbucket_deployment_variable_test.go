package bitbucket

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketDeploymentVariableResource_basic(t *testing.T) {
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
					}`, workspaceSlug, projectName, projectKey, repoName, deploymentVariableNameSecure, deploymentVariableValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.testacc", "key", deploymentVariableNameSecure),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.testacc", "value", deploymentVariableValue),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.testacc", "secured", "true"),

					resource.TestCheckResourceAttrSet("bitbucket_deployment_variable.testacc", "id"),
					resource.TestCheckResourceAttrSet("bitbucket_deployment_variable.testacc", "deployment"),
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
					}`, workspaceSlug, projectName, projectKey, repoName, deploymentVariableNameNotSecure, deploymentVariableValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.testacc", "key", deploymentVariableNameNotSecure),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.testacc", "value", deploymentVariableValue),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.testacc", "secured", "false"),

					resource.TestCheckResourceAttrSet("bitbucket_deployment_variable.testacc", "id"),
					resource.TestCheckResourceAttrSet("bitbucket_deployment_variable.testacc", "deployment"),
				),
			},
		},
	})
}

func TestAccBitbucketDeploymentVariableResource_multipleVars(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	deploymentVariableNameOne := "tf_acc_test" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	deploymentVariableNameTwo := "tf_acc_test" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	deploymentVariableNameThree := "tf_acc_test" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
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

					resource "bitbucket_deployment_variable" "one" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  deployment = bitbucket_deployment.testacc.id
					  key        = "%s"
					  value      = "%s"
					  secured    = false
					}

					resource "bitbucket_deployment_variable" "two" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  deployment = bitbucket_deployment.testacc.id
					  key        = "%s"
					  value      = "%s"
					  secured    = false
					}

					resource "bitbucket_deployment_variable" "three" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  deployment = bitbucket_deployment.testacc.id
					  key        = "%s"
					  value      = "%s"
					  secured    = false
					}`, workspaceSlug, projectName, projectKey, repoName, deploymentVariableNameOne, deploymentVariableValue, deploymentVariableNameTwo, deploymentVariableValue, deploymentVariableNameThree, deploymentVariableValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.one", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.one", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.one", "key", deploymentVariableNameOne),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.one", "value", deploymentVariableValue),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.one", "secured", "false"),
					resource.TestCheckResourceAttrSet("bitbucket_deployment_variable.one", "id"),
					resource.TestCheckResourceAttrSet("bitbucket_deployment_variable.one", "deployment"),

					resource.TestCheckResourceAttr("bitbucket_deployment_variable.two", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.two", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.two", "key", deploymentVariableNameTwo),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.two", "value", deploymentVariableValue),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.two", "secured", "false"),
					resource.TestCheckResourceAttrSet("bitbucket_deployment_variable.two", "id"),
					resource.TestCheckResourceAttrSet("bitbucket_deployment_variable.two", "deployment"),

					resource.TestCheckResourceAttr("bitbucket_deployment_variable.three", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.three", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.three", "key", deploymentVariableNameThree),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.three", "value", deploymentVariableValue),
					resource.TestCheckResourceAttr("bitbucket_deployment_variable.three", "secured", "false"),
					resource.TestCheckResourceAttrSet("bitbucket_deployment_variable.three", "id"),
					resource.TestCheckResourceAttrSet("bitbucket_deployment_variable.three", "deployment"),
				),
			},
		},
	})
}
