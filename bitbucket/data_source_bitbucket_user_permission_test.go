package bitbucket

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketUserPermissionDataSource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	user := os.Getenv("BITBUCKET_MEMBER_ACCOUNT_UUID")

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
					  workspace   = data.bitbucket_workspace.testacc.id
					  project_key = bitbucket_project.testacc.key
					  name        = "%s"
					}

					data "bitbucket_user" "testacc" {
						id = "%s"
					}

					resource "bitbucket_user_permission" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.uuid
					  repository = bitbucket_repository.testacc.name
					  user       = data.bitbucket_user.testacc.id
					  permission = "read"
					}

					data "bitbucket_user_permission" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.uuid
					  repository = bitbucket_repository.testacc.name
					  user       = data.bitbucket_user.testacc.id
					}`, workspaceSlug, projectName, projectKey, repoName, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_user_permission.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("data.bitbucket_user_permission.testacc", "user", user),

					resource.TestCheckResourceAttrSet("data.bitbucket_user_permission.testacc", "id"),
					resource.TestCheckResourceAttrSet("data.bitbucket_user_permission.testacc", "workspace"),
					resource.TestCheckResourceAttrSet("data.bitbucket_user_permission.testacc", "permission"),
				),
			},
		},
	})
}
