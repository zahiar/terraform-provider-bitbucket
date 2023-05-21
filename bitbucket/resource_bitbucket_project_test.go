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

func TestAccBitbucketProjectResource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	projectDescription := "TF ACC Test Project"

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
					}`, workspaceSlug, projectName, projectKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_project.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_project.testacc", "key", projectKey),
					resource.TestCheckResourceAttr("bitbucket_project.testacc", "name", projectName),
					resource.TestCheckResourceAttr("bitbucket_project.testacc", "description", ""),
					resource.TestCheckResourceAttr("bitbucket_project.testacc", "is_private", "true"),
					resource.TestCheckResourceAttrSet("bitbucket_project.testacc", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}

					resource "bitbucket_project" "testacc" {
					  workspace   = data.bitbucket_workspace.testacc.id
					  name        = "%s"
					  key         = "%s"
					  description = "%s"
					  is_private  = false
					}`, workspaceSlug, projectName, projectKey, projectDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_project.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_project.testacc", "key", projectKey),
					resource.TestCheckResourceAttr("bitbucket_project.testacc", "name", projectName),
					resource.TestCheckResourceAttr("bitbucket_project.testacc", "description", projectDescription),
					resource.TestCheckResourceAttr("bitbucket_project.testacc", "is_private", "false"),
					resource.TestCheckResourceAttrSet("bitbucket_project.testacc", "id"),
				),
			},
			{
				ResourceName:      "bitbucket_project.testacc",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s", workspaceSlug, projectKey), nil
				},
			},
		},
	})
}

func TestValidateProjectKey(t *testing.T) {
	invalidKey := "123-invalid-!@£"
	validator := validateProjectKey(invalidKey, nil)
	assert.True(t, validator.HasError())

	validKey := "SOME_123_PROJECT"
	validator = validateProjectKey(validKey, nil)
	assert.False(t, validator.HasError())
}
