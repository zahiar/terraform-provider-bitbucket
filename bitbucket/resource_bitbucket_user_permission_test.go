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

func TestAccBitbucketUserPermissionResource_basic(t *testing.T) {
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
					}`, workspaceSlug, projectName, projectKey, repoName, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_user_permission.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_user_permission.testacc", "user", user),
					resource.TestCheckResourceAttr("bitbucket_user_permission.testacc", "permission", "read"),

					resource.TestCheckResourceAttrSet("bitbucket_user_permission.testacc", "id"),
					resource.TestCheckResourceAttrSet("bitbucket_user_permission.testacc", "workspace"),
				),
			},
			{
				ResourceName:      "bitbucket_user_permission.testacc",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					resources := state.Modules[0].Resources
					workspaceResourceAttr := resources["data.bitbucket_workspace.testacc"].Primary.Attributes
					return fmt.Sprintf("%s/%s/%s", workspaceResourceAttr["uuid"], repoName, user), nil
				},
			},
		},
	})
}

func TestGenerateUserPermissionResourceId(t *testing.T) {
	expected := "{my-workspace-uuid}-my-test-repo-{my-user-uuid}"
	result := generateUserPermissionResourceId("{my-workspace-uuid}", "my-test-repo", "{my-user-uuid}")
	assert.Equal(t, expected, result)
}
