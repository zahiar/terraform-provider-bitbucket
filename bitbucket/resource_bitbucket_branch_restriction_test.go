package bitbucket

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBitbucketBranchRestrictionResource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_USERNAME")
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
					}`, workspaceSlug, projectName, projectKey, repoName, branchRestrictionPattern, branchRestrictionKind),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "pattern", branchRestrictionPattern),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "kind", branchRestrictionKind),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "value", "0"),

					resource.TestCheckNoResourceAttr("bitbucket_branch_restriction.testacc", "users"),
					resource.TestCheckNoResourceAttr("bitbucket_branch_restriction.testacc", "groups"),

					resource.TestCheckResourceAttrSet("bitbucket_branch_restriction.testacc", "id"),
				),
			},
		},
	})
}

func TestAccBitbucketBranchRestrictionResource_withKindAndValueCombination(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_USERNAME")
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
					}`, workspaceSlug, projectName, projectKey, repoName, branchRestrictionPattern, branchRestrictionKind, branchRestrictionValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "pattern", branchRestrictionPattern),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "kind", branchRestrictionKind),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "value", strconv.Itoa(branchRestrictionValue)),

					resource.TestCheckNoResourceAttr("bitbucket_branch_restriction.testacc", "users"),
					resource.TestCheckNoResourceAttr("bitbucket_branch_restriction.testacc", "groups"),

					resource.TestCheckResourceAttrSet("bitbucket_branch_restriction.testacc", "id"),
				),
			},
		},
	})
}

func TestAccBitbucketBranchRestrictionResource_withUsers(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_USERNAME")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	branchRestrictionPattern := "master"
	branchRestrictionKind := "push"
	branchRestrictionUser := os.Getenv("BITBUCKET_USERNAME")

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
					  users      = ["%s"]
					}`, workspaceSlug, projectName, projectKey, repoName, branchRestrictionPattern, branchRestrictionKind, branchRestrictionUser),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "pattern", branchRestrictionPattern),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "kind", branchRestrictionKind),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "value", "0"),

					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "users.#", "1"),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "users.0", branchRestrictionUser),

					resource.TestCheckNoResourceAttr("bitbucket_branch_restriction.testacc", "groups"),

					resource.TestCheckResourceAttrSet("bitbucket_branch_restriction.testacc", "id"),
				),
			},
		},
	})
}

func TestAccBitbucketBranchRestrictionResource_withEmptyUsersAndEmptyGroups(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_USERNAME")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	branchRestrictionPattern := "master"
	branchRestrictionKind := "push"

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
					  users      = []
					  groups     = []
					}`, workspaceSlug, projectName, projectKey, repoName, branchRestrictionPattern, branchRestrictionKind),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "pattern", branchRestrictionPattern),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "kind", branchRestrictionKind),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "value", "0"),

					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "users.#", "0"),
					resource.TestCheckResourceAttr("bitbucket_branch_restriction.testacc", "groups.#", "0"),

					resource.TestCheckResourceAttrSet("bitbucket_branch_restriction.testacc", "id"),
				),
			},
		},
	})
}

func TestParseBranchRestrictionUserFields(t *testing.T) {
	users := []interface{}{"user-a", "user-b", "user-c"}
	usersStrArr := parseBranchRestrictionUserFields(users)

	expected := []string{"user-a", "user-b", "user-c"}
	assert.Equal(t, expected, usersStrArr)
}

func TestParseBranchRestrictionGroupFields(t *testing.T) {
	groups := []interface{}{"group-a", "group-b", "group-c"}
	usersStrArr := parseBranchRestrictionUserGroupFields(groups)

	expected := map[string]string{
		"group-a": "group-a",
		"group-b": "group-b",
		"group-c": "group-c",
	}
	assert.Equal(t, expected, usersStrArr)
}
