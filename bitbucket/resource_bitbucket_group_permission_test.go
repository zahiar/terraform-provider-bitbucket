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

func TestAccBitbucketGroupPermissionResource_basic(t *testing.T) {
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
					  permission = "write"
					}`, workspaceSlug, projectName, projectKey, repoName, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_group_permission.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_group_permission.testacc", "group", groupName),
					resource.TestCheckResourceAttr("bitbucket_group_permission.testacc", "permission", "write"),

					resource.TestCheckResourceAttrSet("bitbucket_group_permission.testacc", "id"),
					resource.TestCheckResourceAttrSet("bitbucket_group_permission.testacc", "workspace"),
				),
			},
			{
				ResourceName:      "bitbucket_group_permission.testacc",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					resources := state.Modules[0].Resources
					workspaceResourceAttr := resources["data.bitbucket_workspace.testacc"].Primary.Attributes
					return fmt.Sprintf("%s/%s/%s", workspaceResourceAttr["uuid"], repoName, groupName), nil
				},
			},
		},
	})
}

func TestGenerateGroupPermissionResourceId(t *testing.T) {
	expected := "{my-workspace-uuid}-my-test-repo-my-test-group"
	result := generateGroupPermissionResourceId("{my-workspace-uuid}", "my-test-repo", "my-test-group")
	assert.Equal(t, expected, result)
}
