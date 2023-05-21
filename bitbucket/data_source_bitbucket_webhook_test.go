package bitbucket

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketWebhookDataSource_basic(t *testing.T) {
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
					  workspace  = data.bitbucket_workspace.testacc.id
					  name       = "%s"
					  key        = "%s"
					  is_private = true
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
					}
	
					data "bitbucket_webhook" "testacc" {
					  id         = bitbucket_webhook.testacc.id
					  workspace  = data.bitbucket_workspace.testacc.id
					  repository = bitbucket_repository.testacc.name
					}`, workspaceSlug, projectName, projectKey, repoName, webhookName, webhookUrl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucket_webhook.testacc", "workspace", workspaceSlug),
					resource.TestCheckResourceAttr("data.bitbucket_webhook.testacc", "repository", repoName),
					resource.TestCheckResourceAttr("data.bitbucket_webhook.testacc", "name", webhookName),
					resource.TestCheckResourceAttr("data.bitbucket_webhook.testacc", "url", webhookUrl),
					resource.TestCheckResourceAttr("data.bitbucket_webhook.testacc", "is_active", "true"),

					resource.TestCheckResourceAttr("data.bitbucket_webhook.testacc", "events.#", "2"),
					resource.TestCheckResourceAttr("data.bitbucket_webhook.testacc", "events.0", "pullrequest:approved"),
					resource.TestCheckResourceAttr("data.bitbucket_webhook.testacc", "events.1", "pullrequest:unapproved"),

					resource.TestCheckResourceAttrSet("data.bitbucket_webhook.testacc", "id"),
				),
			},
		},
	})
}
