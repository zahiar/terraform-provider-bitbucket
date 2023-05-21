package bitbucket

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketGroupPermissionDataSource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	groupName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

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

					resource "bitbucket_group" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.uuid
					  name       = "%s"
					  permission = "read"
					}

					resource "bitbucket_group_permission" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.uuid
					  repository = bitbucket_repository.testacc.name
					  group      = bitbucket_group.testacc.slug
					  permission = "read"
					}

					data "bitbucket_group_permission" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.uuid
					  repository = bitbucket_repository.testacc.name
					  group      = bitbucket_group.testacc.slug
					}`, workspaceSlug, projectName, projectKey, repoName, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_group_permission.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("data.bitbucket_group_permission.testacc", "group", groupName),
					resource.TestCheckResourceAttr("data.bitbucket_group_permission.testacc", "permission", "read"),

					resource.TestCheckResourceAttrSet("data.bitbucket_group_permission.testacc", "id"),
					resource.TestCheckResourceAttrSet("data.bitbucket_group_permission.testacc", "workspace"),
				),
			},
		},
	})
}
