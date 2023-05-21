package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	gobb "github.com/ktrysmt/go-bitbucket"
)

func TestAccBitbucketUserDataSource_basic(t *testing.T) {
	user, _ := getCurrentUser()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "bitbucket_user" "testacc" {
						id = "%s"
					}`, user.Uuid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_user.testacc", "id", user.Uuid),
					resource.TestCheckResourceAttr("data.bitbucket_user.testacc", "nickname", user.Nickname),
					resource.TestCheckResourceAttr("data.bitbucket_user.testacc", "display_name", user.DisplayName),
					resource.TestCheckResourceAttr("data.bitbucket_user.testacc", "account_status", user.AccountStatus),
					resource.TestCheckResourceAttr("data.bitbucket_user.testacc", "account_id", user.AccountId),
				),
			},
		},
	})
}

func getCurrentUser() (*gobb.User, error) {
	if _, isSet := os.LookupEnv("TF_ACC"); isSet {
		client := NewClient()

		return client.User.Profile()
	} else {
		return &gobb.User{
			Uuid:          "",
			Nickname:      "",
			DisplayName:   "",
			AccountStatus: "",
		}, nil
	}
}
