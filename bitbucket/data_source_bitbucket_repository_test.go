package bitbucket

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketRepositoryDataSource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoDescription := "TF ACC Test Repository"

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
					  workspace   = data.bitbucket_workspace.testacc.id
					  name        = "%s"
					  key         = "%s"
					  is_private  = true
					}
	
					resource "bitbucket_repository" "testacc" {
					  workspace   = data.bitbucket_workspace.testacc.id
					  project_key = bitbucket_project.testacc.key
					  name        = "%s"
					  description = "%s"
					  is_private  = true
					  has_wiki    = false
					  fork_policy = "no_forks"
					}
	
					data "bitbucket_repository" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  name       = "%s"
					  depends_on = [bitbucket_repository.testacc]
					}`, workspaceSlug, projectName, projectKey, repoName, repoDescription, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "name", repoName),
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "project_key", projectKey),
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "description", repoDescription),
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "is_private", "true"),
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "has_wiki", "false"),
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "fork_policy", "no_forks"),
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "enable_pipelines", "false"),
					resource.TestCheckResourceAttrSet("data.bitbucket_repository.testacc", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}
	
					resource "bitbucket_project" "testacc" {
					  workspace   = data.bitbucket_workspace.testacc.id
					  name        = "%s"
					  key         = "%s"
					  is_private  = true
					}
	
					resource "bitbucket_repository" "testacc" {
					  workspace        = data.bitbucket_workspace.testacc.id
					  project_key      = bitbucket_project.testacc.key
					  name             = "%s"
					  description      = "%s"
					  is_private       = true
					  has_wiki         = true
					  fork_policy      = "no_forks"
					  enable_pipelines = true
					}
	
					data "bitbucket_repository" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  name       = "%s"
					  depends_on = [bitbucket_repository.testacc]
					}`, workspaceSlug, projectName, projectKey, repoName, repoDescription, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "name", repoName),
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "project_key", projectKey),
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "description", repoDescription),
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "is_private", "true"),
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "has_wiki", "true"),
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "fork_policy", "no_forks"),
					resource.TestCheckResourceAttr("data.bitbucket_repository.testacc", "enable_pipelines", "true"),
					resource.TestCheckResourceAttrSet("data.bitbucket_repository.testacc", "id"),
				),
			},
		},
	})
}
