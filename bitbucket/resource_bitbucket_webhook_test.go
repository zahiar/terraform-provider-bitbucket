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

func TestAccBitbucketWebhookResource_basic(t *testing.T) {
	workspaceSlug := os.Getenv("BITBUCKET_WORKSPACE")
	projectName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectKey := strings.ToUpper(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	repoName := "tf-acc-test-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	webhookName := "TF ACC Test Webhook"
	webhookUrl := "https://example.com"

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
	
					resource "bitbucket_webhook" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  name       = "%s"
					  url        = "%s"
					  events      = ["pullrequest:approved"]
					}`, workspaceSlug, projectName, projectKey, repoName, webhookName, webhookUrl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "name", webhookName),
					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "url", webhookUrl),
					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "is_active", "false"),

					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "events.#", "1"),
					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "events.0", "pullrequest:approved"),

					resource.TestCheckResourceAttrSet("bitbucket_webhook.testacc", "id"),
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
					}
	
					resource "bitbucket_webhook" "testacc" {
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					  name       = "%s"
					  url        = "%s"
					  events     = ["pullrequest:approved", "pullrequest:unapproved"]
					  is_active  = true
					}`, workspaceSlug, projectName, projectKey, repoName, webhookName, webhookUrl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "name", webhookName),
					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "url", webhookUrl),
					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "is_active", "true"),

					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "events.#", "2"),
					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "events.0", "pullrequest:approved"),
					resource.TestCheckResourceAttr("bitbucket_webhook.testacc", "events.1", "pullrequest:unapproved"),

					resource.TestCheckResourceAttrSet("bitbucket_webhook.testacc", "id"),
				),
			},
			{
				ResourceName:      "bitbucket_webhook.testacc",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					resources := state.Modules[0].Resources
					webhookResourceAttr := resources["bitbucket_webhook.testacc"].Primary.Attributes
					return fmt.Sprintf("%s/%s/%s", workspaceSlug, repoName, webhookResourceAttr["id"]), nil
				},
			},
		},
	})
}

func TestConvertEventsToStringArray(t *testing.T) {
	events := []interface{}{"a", "b", "c"}
	eventsStrArr := convertEventsToStringArray(events)

	expected := []string{"a", "b", "c"}
	assert.Equal(t, expected, eventsStrArr)
}
