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

func TestAccBitbucketRepositoryResource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoDescription := "TF ACC Test Repository"
	repoForkPolicy := "no_forks"

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
					}

					resource "bitbucket_repository" "testacc" {
					  workspace   = data.bitbucket_workspace.testacc.id
					  project_key = bitbucket_project.testacc.key
					  name        = "%s"
					}`, workspaceSlug, projectName, projectKey, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "name", repoName),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "project_key", projectKey),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "description", ""),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "is_private", "true"),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "has_wiki", "false"),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "fork_policy", "no_forks"),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "enable_pipelines", "false"),
					resource.TestCheckResourceAttrSet("bitbucket_repository.testacc", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}

					resource "bitbucket_project" "testacc" {
					  workspace = data.bitbucket_workspace.testacc.id
					  name      = "%s"
					  key       = "%s"
					}

					resource "bitbucket_repository" "testacc" {
					  workspace   = data.bitbucket_workspace.testacc.id
					  project_key = bitbucket_project.testacc.key
					  name        = "%s"
					  description = "%s"
					  fork_policy = "%s"
					}`, workspaceSlug, projectName, projectKey, repoName, repoDescription, repoForkPolicy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "name", repoName),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "project_key", projectKey),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "description", repoDescription),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "is_private", "true"),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "has_wiki", "false"),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "fork_policy", repoForkPolicy),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "enable_pipelines", "false"),
					resource.TestCheckResourceAttrSet("bitbucket_repository.testacc", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}

					resource "bitbucket_project" "testacc" {
					  workspace = data.bitbucket_workspace.testacc.id
					  name      = "%s"
					  key       = "%s"
					}

					resource "bitbucket_repository" "testacc" {
					  workspace        = data.bitbucket_workspace.testacc.id
					  project_key      = bitbucket_project.testacc.key
					  name             = "%s"
					  description      = "%s"
					  fork_policy      = "%s"
					  enable_pipelines = true
					  has_wiki         = true
					}`, workspaceSlug, projectName, projectKey, repoName, repoDescription, repoForkPolicy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "name", repoName),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "project_key", projectKey),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "description", repoDescription),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "is_private", "true"),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "has_wiki", "true"),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "fork_policy", repoForkPolicy),
					resource.TestCheckResourceAttr("bitbucket_repository.testacc", "enable_pipelines", "true"),
					resource.TestCheckResourceAttrSet("bitbucket_repository.testacc", "id"),
				),
			},
			{
				ResourceName:      "bitbucket_repository.testacc",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s", workspaceSlug, repoName), nil
				},
			},
		},
	})
}

func TestValidateRepositoryName(t *testing.T) {
	invalidName := "ABC!@£"
	validator := validateRepositoryName(invalidName, nil)
	assert.True(t, validator.HasError())

	validNames := []string{"abc-def", "foo.bar"}

	for _, validName := range validNames {
		validator = validateRepositoryName(validName, nil)
		assert.False(t, validator.HasError())
	}
}
