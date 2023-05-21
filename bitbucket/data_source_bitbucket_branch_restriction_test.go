package bitbucket

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketBranchRestrictionDataSource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	branchRestrictionPattern := "master"
	branchRestrictionKind := "require_tasks_to_be_completed"

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
					}
	
					resource "bitbucket_branch_restriction" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  pattern    = "%s"
					  kind       = "%s"
					}

					data "bitbucket_branch_restriction" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  id         = bitbucket_branch_restriction.testacc.id
					}`, workspaceSlug, projectName, projectKey, repoName, branchRestrictionPattern, branchRestrictionKind),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_branch_restriction.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("data.bitbucket_branch_restriction.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("data.bitbucket_branch_restriction.testacc", "pattern", branchRestrictionPattern),
					resource.TestCheckResourceAttr("data.bitbucket_branch_restriction.testacc", "kind", branchRestrictionKind),
					resource.TestCheckResourceAttr("data.bitbucket_branch_restriction.testacc", "value", "0"),

					resource.TestCheckNoResourceAttr("data.bitbucket_branch_restriction.testacc", "users"),
					resource.TestCheckNoResourceAttr("data.bitbucket_branch_restriction.testacc", "groups"),

					resource.TestCheckResourceAttrSet("data.bitbucket_branch_restriction.testacc", "id"),
				),
			},
		},
	})
}

func TestAccBitbucketBranchRestrictionDataSource_withKindAndValueCombination(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	branchRestrictionPattern := "master"
	branchRestrictionKind := "require_passing_builds_to_merge"
	branchRestrictionValue := 3

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
					}
	
					resource "bitbucket_branch_restriction" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  pattern    = "%s"
					  kind       = "%s"
					  value      = %d
					}

					data "bitbucket_branch_restriction" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  id         = bitbucket_branch_restriction.testacc.id
					}`, workspaceSlug, projectName, projectKey, repoName, branchRestrictionPattern, branchRestrictionKind, branchRestrictionValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_branch_restriction.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("data.bitbucket_branch_restriction.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("data.bitbucket_branch_restriction.testacc", "pattern", branchRestrictionPattern),
					resource.TestCheckResourceAttr("data.bitbucket_branch_restriction.testacc", "kind", branchRestrictionKind),
					resource.TestCheckResourceAttr("data.bitbucket_branch_restriction.testacc", "value", strconv.Itoa(branchRestrictionValue)),

					resource.TestCheckNoResourceAttr("data.bitbucket_branch_restriction.testacc", "users"),
					resource.TestCheckNoResourceAttr("data.bitbucket_branch_restriction.testacc", "groups"),

					resource.TestCheckResourceAttrSet("data.bitbucket_branch_restriction.testacc", "id"),
				),
			},
		},
	})
}
