package bitbucket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gobb "github.com/ktrysmt/go-bitbucket"
)

func resourceBitBucketWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBitBucketWebhookCreate,
		ReadContext:   resourceBitBucketWebhookRead,
		UpdateContext: resourceBitBucketWebhookUpdate,
		DeleteContext: resourceBitBucketWebhookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBitBucketWebhookImport,
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

func resourceBitBucketWebhookCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	webhookResponse, err := client.Repositories.Webhooks.Create(
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

	webhook, err := decodeWebhookResponse(webhookResponse)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to decode webhook response with error: %s", err))
	}

	resourceData.SetId(webhook.Uuid)

	return resourceBitBucketWebhookRead(ctx, resourceData, meta)
}

func resourceBitBucketWebhookRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	webhookResponse, err := client.Repositories.Webhooks.Get(
		&gobb.WebhooksOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
			Uuid:     resourceData.Get("id").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get webhook with error: %s", err))
	}

	webhook, err := decodeWebhookResponse(webhookResponse)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to decode webhook response with error: %s", err))
	}

	_ = resourceData.Set("name", webhook.Description)
	_ = resourceData.Set("url", webhook.Url)
	_ = resourceData.Set("is_active", webhook.Active)
	_ = resourceData.Set("events", webhook.Events)
	resourceData.SetId(webhook.Uuid)

	return nil
}

func resourceBitBucketWebhookUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return resourceBitBucketWebhookRead(ctx, resourceData, meta)
}

func resourceBitBucketWebhookDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceBitBucketWebhookImport(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{resourceData}

	splitID := strings.Split(resourceData.Id(), "/")
	if len(splitID) < 2 {
		return ret, fmt.Errorf("invalid import ID. It must to be in this format \"<workspace-slug|workspace-uuid>/<repository-name>/<webhook-uuid>\"")
	}

	_ = resourceData.Set("workspace", splitID[0])
	_ = resourceData.Set("repository", splitID[1])
	resourceData.SetId(splitID[2])

	_ = resourceBitBucketWebhookRead(ctx, resourceData, meta)

	return ret, nil
}

func decodeWebhookResponse(response interface{}) (*gobb.WebhooksOptions, error) {
	webhookMap := response.(map[string]interface{})

	if webhookMap["type"] == "error" {
		return nil, errors.New("unable able to decode webhook API response")
	}

	webHookOptions := &gobb.WebhooksOptions{}
	jsonString, err := json.Marshal(webhookMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonString, &webHookOptions)

	return webHookOptions, err
}

func convertEventsToStringArray(events []interface{}) []string {
	var eventArray []string

	for _, event := range events {
		eventArray = append(eventArray, event.(string))
	}
	return eventArray
}
