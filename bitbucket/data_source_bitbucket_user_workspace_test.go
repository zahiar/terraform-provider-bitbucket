package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketUserWorkspaceDataSource_basic(t *testing.T) {
	workspace := os.Getenv("BITBUCKET_WORKSPACE")
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

					data "bitbucket_user_workspace" "testacc" {
						nickname  = "%s"
						workspace = data.bitbucket_workspace.testacc.id
					}`, workspace, user.Nickname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_user_workspace.testacc", "id", user.Uuid),
					resource.TestCheckResourceAttr("data.bitbucket_user_workspace.testacc", "nickname", user.Nickname),
					resource.TestCheckResourceAttr("data.bitbucket_user_workspace.testacc", "display_name", user.DisplayName),
					resource.TestCheckResourceAttr("data.bitbucket_user_workspace.testacc", "account_status", user.AccountStatus),
					resource.TestCheckResourceAttr("data.bitbucket_user_workspace.testacc", "account_id", user.AccountId),
					resource.TestCheckResourceAttr("data.bitbucket_user_workspace.testacc", "workspace", workspace),
				),
			},
		},
	})
}
