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

func TestAccBitbucketDefaultReviewerResource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoDescription := "TF ACC Test Repository"
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

					resource "bitbucket_project" "testacc" {
					  workspace   = data.bitbucket_workspace.testacc.id
					  name        = "%s"
					  key         = "%s"
					  is_private  = true
					}

					resource "bitbucket_repository" "testacc" {
					  workspace   = data.bitbucket_workspace.testacc.id
					  project_key = bitbucket_project.testacc.key
					  name        = "%s"
					  description = "%s"
					  is_private  = true
					  fork_policy = "no_forks"
					}

					data "bitbucket_user" "testacc" {
						id = "%s"
					}

					resource "bitbucket_default_reviewer" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  user       = data.bitbucket_user.testacc.id
					}`, workspaceSlug, projectName, projectKey, repoName, repoDescription, user.Uuid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_default_reviewer.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_default_reviewer.testacc", "repository", repoName),
					resource.TestCheckResourceAttrSet("bitbucket_default_reviewer.testacc", "id"),
					resource.TestCheckResourceAttrSet("bitbucket_default_reviewer.testacc", "user"),
				),
			},
			{
				ResourceName:      "bitbucket_default_reviewer.testacc",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", workspaceSlug, repoName, user.Uuid), nil
				},
			},
		},
	})
}

func TestGenerateDefaultReviewerResourceId(t *testing.T) {
	workspace := "my-workspace"
	repository := "my-repository"
	user := "{123ab4cd-5678-9e01-f234-5678g9h01i2j}"

	result := generateDefaultReviewerResourceId(workspace, repository, user)
	assert.Equal(t, "my-workspace-my-repository-{123ab4cd-5678-9e01-f234-5678g9h01i2j}", result)
}
