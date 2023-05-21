package bitbucket

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	gobb "github.com/ktrysmt/go-bitbucket"
	"github.com/stretchr/testify/assert"
)

func TestAccBitbucketDeploymentResource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	deploymentName := "TF ACC Test Deployment"

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
					  name        = "%s"
					  environment = "Test"
					}`, workspaceSlug, projectName, projectKey, repoName, deploymentName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_deployment.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_deployment.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_deployment.testacc", "name", deploymentName),
					resource.TestCheckResourceAttr("bitbucket_deployment.testacc", "environment", "Test"),

					resource.TestCheckResourceAttrSet("bitbucket_deployment.testacc", "id"),
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
					  name        = "%s"
					  environment = "Staging"
					}`, workspaceSlug, projectName, projectKey, repoName, deploymentName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_deployment.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_deployment.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_deployment.testacc", "name", deploymentName),
					resource.TestCheckResourceAttr("bitbucket_deployment.testacc", "environment", "Staging"),

					resource.TestCheckResourceAttrSet("bitbucket_deployment.testacc", "id"),
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
					  name        = "%s"
					  environment = "Production"
					}`, workspaceSlug, projectName, projectKey, repoName, deploymentName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_deployment.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_deployment.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_deployment.testacc", "name", deploymentName),
					resource.TestCheckResourceAttr("bitbucket_deployment.testacc", "environment", "Production"),

					resource.TestCheckResourceAttrSet("bitbucket_deployment.testacc", "id"),
				),
			},
		},
	})
}

func TestGetDeploymentEnvironmentIntValue(t *testing.T) {
	validEnvironments := []string{gobb.Test.String(), gobb.Staging.String(), gobb.Production.String()}
	for _, name := range validEnvironments {
		_, err := getDeploymentEnvironmentIntValue(name)
		assert.Nil(t, err)
	}

	_, err := getDeploymentEnvironmentIntValue("not-valid")
	assert.Error(t, err)
}
