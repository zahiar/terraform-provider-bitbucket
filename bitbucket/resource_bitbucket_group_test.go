package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccBitbucketGroupResource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
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
					  workspace  = data.bitbucket_workspace.testacc.uuid
					  name       = "%s"
					  permission = "read"
					}`, workspaceSlug, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "name", groupName),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "permission", "read"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "slug", groupName),

					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "workspace"),
					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "id"),
				),
			},
			{
				ResourceName:      "bitbucket_group.testacc",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					resources := state.Modules[0].Resources
					workspaceResourceAttr := resources["data.bitbucket_workspace.testacc"].Primary.Attributes
					return fmt.Sprintf("%s/%s", workspaceResourceAttr["uuid"], groupName), nil
				},
			},
		},
	})
}

func TestAccBitbucketGroupResource_changeName(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
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
					  workspace  = data.bitbucket_workspace.testacc.uuid
					  name       = "%s"
					  permission = "read"
					}`, workspaceSlug, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "name", groupName),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "permission", "read"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "slug", groupName),

					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "workspace"),
					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}

					resource "bitbucket_group" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.uuid
					  name       = "%s"
					  permission = "read"
					}`, workspaceSlug, newGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "name", newGroupName),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "permission", "read"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "slug", newGroupName),

					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "workspace"),
					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "id"),
				),
			},
			{
				ResourceName:      "bitbucket_group.testacc",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					resources := state.Modules[0].Resources
					workspaceResourceAttr := resources["data.bitbucket_workspace.testacc"].Primary.Attributes
					return fmt.Sprintf("%s/%s", workspaceResourceAttr["uuid"], newGroupName), nil
				},
			},
		},
	})
}

func TestAccBitbucketGroupResource_changeProperties(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
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
					  workspace  = data.bitbucket_workspace.testacc.uuid
					  name       = "%s"
					  permission = "read"
					}`, workspaceSlug, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "name", groupName),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "permission", "read"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "slug", groupName),

					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "workspace"),
					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}

					resource "bitbucket_group" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.uuid
					  name       = "%s"
					  permission = "write"
					}`, workspaceSlug, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "name", groupName),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "permission", "write"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "slug", groupName),

					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "workspace"),
					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "id"),
				),
			},
			{
				ResourceName:      "bitbucket_group.testacc",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					resources := state.Modules[0].Resources
					workspaceResourceAttr := resources["data.bitbucket_workspace.testacc"].Primary.Attributes
					return fmt.Sprintf("%s/%s", workspaceResourceAttr["uuid"], groupName), nil
				},
			},
		},
	})
}

func TestAccBitbucketGroupResource_withoutProperties(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
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
					  workspace = data.bitbucket_workspace.testacc.uuid
					  name      = "%s"
					}`, workspaceSlug, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "name", groupName),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "permission", "none"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "slug", groupName),

					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "workspace"),
					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					data "bitbucket_workspace" "testacc" {
						id = "%s"
					}

					resource "bitbucket_group" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.uuid
					  name       = "%s"
					  permission = "write"
					}`, workspaceSlug, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "name", groupName),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "permission", "write"),
					resource.TestCheckResourceAttr("bitbucket_group.testacc", "slug", groupName),

					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "workspace"),
					resource.TestCheckResourceAttrSet("bitbucket_group.testacc", "id"),
				),
			},
			{
				ResourceName:      "bitbucket_group.testacc",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					resources := state.Modules[0].Resources
					workspaceResourceAttr := resources["data.bitbucket_workspace.testacc"].Primary.Attributes
					return fmt.Sprintf("%s/%s", workspaceResourceAttr["uuid"], groupName), nil
				},
			},
		},
	})
}

func TestGenerateGroupResourceId(t *testing.T) {
	expected := "{my-workspace-uuid}-my-test-group"
	result := generateGroupResourceId("{my-workspace-uuid}", "my-test-group")
	assert.Equal(t, expected, result)
}
