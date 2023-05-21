package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketGroupDataSource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
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
	
					resource "bitbucket_group" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  name       = "%s"
					  permission = "read"
					}
	
					data "bitbucket_group" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  slug       = bitbucket_group.testacc.slug
					}`, workspaceSlug, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_group.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("data.bitbucket_group.testacc", "name", groupName),
					resource.TestCheckResourceAttr("data.bitbucket_group.testacc", "permission", "read"),
					resource.TestCheckResourceAttr("data.bitbucket_group.testacc", "slug", groupName),
					resource.TestCheckResourceAttrSet("data.bitbucket_group.testacc", "id"),
				),
			},
		},
	})
}
