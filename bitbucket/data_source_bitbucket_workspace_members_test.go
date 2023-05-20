package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketWorkspaceMembersDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}

					data "bitbucket_workspace_members" "testacc" {
						workspace = data.bitbucket_workspace.testacc.id
					}`, os.Getenv("BITBUCKET_USERNAME")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_workspace_members.testacc", "id", os.Getenv("BITBUCKET_USERNAME")),
					resource.TestCheckResourceAttr("data.bitbucket_workspace_members.testacc", "workspace", os.Getenv("BITBUCKET_USERNAME")),

					resource.TestCheckResourceAttrSet("data.bitbucket_workspace_members.testacc", "members.#"),
					resource.TestCheckResourceAttrSet("data.bitbucket_workspace_members.testacc", "members.0.id"),
					resource.TestCheckResourceAttrSet("data.bitbucket_workspace_members.testacc", "members.0.nickname"),
				),
			},
		},
	})
}
