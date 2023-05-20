package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketWorkspaceProjectsDataSource_basic(t *testing.T) {
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

					data "bitbucket_workspace_projects" "testacc" {
						workspace  = data.bitbucket_workspace.testacc.id
						depends_on = [bitbucket_project.testacc]
					}`, os.Getenv("BITBUCKET_USERNAME"), projectName, projectKey, projectDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_workspace_projects.testacc", "id", os.Getenv("BITBUCKET_USERNAME")),
					resource.TestCheckResourceAttr("data.bitbucket_workspace_projects.testacc", "workspace", os.Getenv("BITBUCKET_USERNAME")),

					resource.TestCheckResourceAttrSet("data.bitbucket_workspace_projects.testacc", "projects.#"),
					resource.TestCheckResourceAttrSet("data.bitbucket_workspace_projects.testacc", "projects.0.id"),
					resource.TestCheckResourceAttr("data.bitbucket_workspace_projects.testacc", "projects.0.name", projectName),
					resource.TestCheckResourceAttr("data.bitbucket_workspace_projects.testacc", "projects.0.key", projectKey),
				),
			},
		},
	})
}
