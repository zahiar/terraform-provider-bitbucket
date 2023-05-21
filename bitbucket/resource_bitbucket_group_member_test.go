package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccBitbucketGroupMemberResource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	user, _ := getCurrentUser()
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

					data "bitbucket_user" "testacc" {
						id = "%s"
					}

					resource "bitbucket_group" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.uuid
					  name       = "%s"
					  permission = "read"
					}

					resource "bitbucket_group_member" "testacc" {
					  workspace = data.bitbucket_workspace.testacc.uuid
					  group     = bitbucket_group.testacc.slug
					  user      = data.bitbucket_user.testacc.id
					}`, workspaceSlug, user.Uuid, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_group_member.testacc", "workspace", user.Uuid),
					resource.TestCheckResourceAttr("bitbucket_group_member.testacc", "group", groupName),
					resource.TestCheckResourceAttr("bitbucket_group_member.testacc", "user", user.Uuid),

					resource.TestCheckResourceAttrSet("bitbucket_group_member.testacc", "id"),
				),
			},
			{
				ResourceName:      "bitbucket_group_member.testacc",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					resources := state.Modules[0].Resources
					workspaceResourceAttr := resources["data.bitbucket_workspace.testacc"].Primary.Attributes
					userResourceAttr := resources["data.bitbucket_user.testacc"].Primary.Attributes
					return fmt.Sprintf("%s/%s/%s", workspaceResourceAttr["uuid"], groupName, userResourceAttr["id"]), nil
				},
			},
		},
	})
}

func TestGenerateGroupMemberResourceId(t *testing.T) {
	expected := "{my-workspace-uuid}-my-test-group-{my-user-uuid}"
	result := generateGroupPermissionResourceId("{my-workspace-uuid}", "my-test-group", "{my-user-uuid}")
	assert.Equal(t, expected, result)
}
