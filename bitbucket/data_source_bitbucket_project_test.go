package bitbucket

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketProjectDataSource_basic(t *testing.T) {
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
					  workspace   = data.bitbucket_workspace.testacc.id
					  name        = "%s"
					  key         = "%s"
					  description = "%s"
					  is_private  = true
					}
	
					data "bitbucket_project" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  key        = "%s"
					  depends_on = [bitbucket_project.testacc]
					}`, workspaceSlug, projectName, projectKey, projectDescription, projectKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_project.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("data.bitbucket_project.testacc", "key", projectKey),
					resource.TestCheckResourceAttr("data.bitbucket_project.testacc", "name", projectName),
					resource.TestCheckResourceAttr("data.bitbucket_project.testacc", "description", projectDescription),
					resource.TestCheckResourceAttr("data.bitbucket_project.testacc", "is_private", "true"),
					resource.TestCheckResourceAttrSet("data.bitbucket_project.testacc", "id"),
				),
			},
		},
	})
}
