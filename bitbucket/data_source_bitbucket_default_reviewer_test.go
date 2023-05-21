package bitbucket

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketDefaultReviewerDataSource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoDescription := "TF ACC Test Repository"
	user, _ := getCurrentUser()

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
					  fork_policy = "no_forks"
					}

					data "bitbucket_user" "testacc" {
						id = "%s"
					}

					resource "bitbucket_default_reviewer" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  user       = data.bitbucket_user.testacc.id
					}

					data "bitbucket_default_reviewer" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  user       = data.bitbucket_user.testacc.id
					  depends_on = [bitbucket_default_reviewer.testacc]
					}`, workspaceSlug, projectName, projectKey, repoName, repoDescription, user.Uuid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_default_reviewer.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("data.bitbucket_default_reviewer.testacc", "repository", repoName),
					resource.TestCheckResourceAttrSet("data.bitbucket_default_reviewer.testacc", "id"),
					resource.TestCheckResourceAttrSet("data.bitbucket_default_reviewer.testacc", "user"),
				),
			},
		},
	})
}
