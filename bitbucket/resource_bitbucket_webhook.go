package bitbucket

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gobb "github.com/ktrysmt/go-bitbucket"
)

func resourceBitbucketWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBitbucketWebhookCreate,
		ReadContext:   resourceBitbucketWebhookRead,
		UpdateContext: resourceBitbucketWebhookUpdate,
		DeleteContext: resourceBitbucketWebhookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBitbucketWebhookImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The UUID of the webhook.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace": {
				Description: "The slug or UUID (including the enclosing `{}`) of the workspace this webhook belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"repository": {
				Description:      "The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateRepositoryName,
			},
			"name": {
				Description: "The name of the webhook.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"url": {
				Description:  "The url to configure the webhook with.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"events": {
				Description: "A list of events that will trigger the webhook - see docs for a complete list.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"pullrequest:unapproved", "issue:comment_created", "repo:imported", "repo:created", "repo:commit_comment_created", "pullrequest:approved", "pullrequest:comment_updated", "issue:updated", "project:updated", "repo:deleted", "pullrequest:changes_request_created", "pullrequest:comment_created", "repo:commit_status_updated", "pullrequest:updated", "issue:created", "repo:fork", "pullrequest:comment_deleted", "repo:commit_status_created", "repo:updated", "pullrequest:rejected", "pullrequest:fulfilled", "pullrequest:created", "pullrequest:changes_request_removed", "repo:transfer", "repo:push"}, false),
				},
				Required: true,
			},
			"is_active": {
				Description: "A boolean to state if the webhook is active or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceBitbucketWebhookCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	webhook, err := client.Repositories.Webhooks.Create(
		&gobb.WebhooksOptions{
			Owner:       resourceData.Get("workspace").(string),
			RepoSlug:    resourceData.Get("repository").(string),
			Description: resourceData.Get("name").(string),
			Url:         resourceData.Get("url").(string),
			Active:      resourceData.Get("is_active").(bool),
			Events:      convertEventsToStringArray(resourceData.Get("events").([]interface{})),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to create webhook with error: %s", err))
	}

	resourceData.SetId(webhook.Uuid)

	return resourceBitbucketWebhookRead(ctx, resourceData, meta)
}

func resourceBitbucketWebhookRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	webhook, err := client.Repositories.Webhooks.Get(
		&gobb.WebhooksOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
			Uuid:     resourceData.Get("id").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get webhook with error: %s", err))
	}

	_ = resourceData.Set("name", webhook.Description)
	_ = resourceData.Set("url", webhook.Url)
	_ = resourceData.Set("is_active", webhook.Active)
	_ = resourceData.Set("events", webhook.Events)
	resourceData.SetId(webhook.Uuid)

	return nil
}

func resourceBitbucketWebhookUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	_, err := client.Repositories.Webhooks.Update(
		&gobb.WebhooksOptions{
			Uuid:        resourceData.Id(),
			Owner:       resourceData.Get("workspace").(string),
			RepoSlug:    resourceData.Get("repository").(string),
			Description: resourceData.Get("name").(string),
			Url:         resourceData.Get("url").(string),
			Active:      resourceData.Get("is_active").(bool),
			Events:      convertEventsToStringArray(resourceData.Get("events").([]interface{})),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to update webhook with error: %s", err))
	}

	return resourceBitbucketWebhookRead(ctx, resourceData, meta)
}

func resourceBitbucketWebhookDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	_, err := client.Repositories.Webhooks.Delete(
		&gobb.WebhooksOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
			Uuid:     resourceData.Id(),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete webhook with error: %s", err))
	}

	resourceData.SetId("")

	return nil
}

func resourceBitbucketWebhookImport(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{resourceData}

	splitID := strings.Split(resourceData.Id(), "/")
	if len(splitID) < 2 {
		return ret, fmt.Errorf("invalid import ID. It must to be in this format \"<workspace-slug|workspace-uuid>/<repository-name>/<webhook-uuid>\"")
	}

	_ = resourceData.Set("workspace", splitID[0])
	_ = resourceData.Set("repository", splitID[1])
	resourceData.SetId(splitID[2])

	_ = resourceBitbucketWebhookRead(ctx, resourceData, meta)

	return ret, nil
}

func convertEventsToStringArray(events []interface{}) []string {
	var eventArray []string

	for _, event := range events {
		eventArray = append(eventArray, event.(string))
	}
	return eventArray
}
