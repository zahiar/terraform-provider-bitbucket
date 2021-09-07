package bitbucket

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"bitbucket": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	err := testAccProvider.InternalValidate()
	assert.NoError(t, err)
}

func testAccPreCheck(t *testing.T) {
	username := os.Getenv("BITBUCKET_USERNAME")
	assert.NotEqual(t, "", username, "BITBUCKET_USERNAME must be set for acceptance tests")

	password := os.Getenv("BITBUCKET_PASSWORD")
	assert.NotEqual(t, "", password, "BITBUCKET_PASSWORD must be set for acceptance tests")
}
