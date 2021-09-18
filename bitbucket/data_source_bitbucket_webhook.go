package bitbucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketWebhook() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBitbucketWebhookRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The UUID of the webhook.",
				Type:        schema.TypeString,
				Required:    true,
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
				Computed:    true,
			},
			"url": {
				Description: "The url to configure the webhook with.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"events": {
				Description: "A list of events that will trigger the webhook - see docs for a complete list.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"is_active": {
				Description: "A boolean to state if the webhook is active or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func dataSourceBitbucketWebhookRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceBitbucketWebhookRead(ctx, resourceData, meta)
}
