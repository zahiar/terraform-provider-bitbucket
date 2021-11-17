package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	gobb "github.com/ktrysmt/go-bitbucket"
	"github.com/stretchr/testify/assert"
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
					resource.TestCheckResourceAttr("data.bitbucket_user.testacc", "nickname", user.Nickname),
					resource.TestCheckResourceAttr("data.bitbucket_user.testacc", "display_name", user.DisplayName),
					resource.TestCheckResourceAttr("data.bitbucket_user.testacc", "account_status", user.AccountStatus),
					resource.TestCheckResourceAttrSet("data.bitbucket_user.testacc", "id"),
					resource.TestCheckResourceAttrSet("data.bitbucket_user.testacc", "account_id"),
				),
			},
		},
	})
}

func getCurrentUser() (*gobb.User, error) {
	if _, isSet := os.LookupEnv("TF_ACC"); isSet {
		client := gobb.NewBasicAuth(
			os.Getenv("BITBUCKET_USERNAME"),
			os.Getenv("BITBUCKET_PASSWORD"),
		)

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

func TestDecodeUserResponseSuccess(t *testing.T) {
	resp := map[string]interface{}{
		"uuid":           "example-uuid",
		"display_name":   "example-display-name",
		"nickname":       "example-nickname",
		"account_id":     "example-account-id",
		"account_status": "example-account-status",
	}

	user, _ := decodeUserResponse(resp)

	expectedUser := &User{
		Uuid:          "example-uuid",
		DisplayName:   "example-display-name",
		Nickname:      "example-nickname",
		AccountId:     "example-account-id",
		AccountStatus: "example-account-status",
	}
	assert.Equal(t, expectedUser, user)
}

func TestDecodeUserResponseError(t *testing.T) {
	resp := map[string]interface{}{
		"type": "error",
	}

	_, err := decodeUserResponse(resp)
	assert.EqualError(t, err, "unable able to decode user API response")
}
