package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketWorkspaceDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}`, os.Getenv("BITBUCKET_WORKSPACE")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_workspace.testacc", "id", os.Getenv("BITBUCKET_WORKSPACE")),
					resource.TestCheckResourceAttr("data.bitbucket_workspace.testacc", "type", "workspace"),
					resource.TestCheckResourceAttrSet("data.bitbucket_workspace.testacc", "uuid"),
					resource.TestCheckResourceAttrSet("data.bitbucket_workspace.testacc", "is_private"),
					resource.TestCheckResourceAttrSet("data.bitbucket_workspace.testacc", "name"),
				),
			},
		},
	})
}
