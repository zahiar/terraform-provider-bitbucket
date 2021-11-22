package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBitbucketGroupResource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_USERNAME")
	groupName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}
	
					resource "bitbucket_group" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  name       = "%s"
                      auto_add   = true
                      permission = "read"
					}`, workspaceSlug, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "name", groupName),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "auto_add", "true"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "permission", "read"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "slug", groupName),
					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "id"),
				),
			},
		},
	})
}

func TestAccBitbucketGroupResource_changeName(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_USERNAME")
	groupName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	newGroupName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}

					resource "bitbucket_group" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  name       = "%s"
                      auto_add   = true
                      permission = "read"
					}`, workspaceSlug, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "name", groupName),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "auto_add", "true"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "permission", "read"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "slug", groupName),
					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}

					resource "bitbucket_group" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  name       = "%s"
                      auto_add   = true
                      permission = "read"
					}`, workspaceSlug, newGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "name", newGroupName),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "auto_add", "true"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "permission", "read"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "slug", newGroupName),
					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "id"),
				),
			},
		},
	})
}

func TestAccBitbucketGroupResource_changeProperties(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_USERNAME")
	groupName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}

					resource "bitbucket_group" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  name       = "%s"
                      auto_add   = true
                      permission = "read"
					}`, workspaceSlug, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "name", groupName),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "auto_add", "true"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "permission", "read"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "slug", groupName),
					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}

					resource "bitbucket_group" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  name       = "%s"
                      auto_add   = false
                      permission = "write"
					}`, workspaceSlug, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "name", groupName),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "auto_add", "false"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "permission", "write"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "slug", groupName),
					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "id"),
				),
			},
		},
	})
}

func TestGenerateGroupResourceId(t *testing.T) {
	expected := "{my-workspace-uuid}-my-test-group"
	result := generateGroupResourceId("{my-workspace-uuid}", "my-test-group")
	assert.Equal(t, expected, result)
}
