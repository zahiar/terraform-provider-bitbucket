package bitbucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitBucketWebhook() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBitBucketWebhookRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateRepositoryName,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"events": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceBitBucketWebhookRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceBitBucketWebhookRead(ctx, resourceData, meta)
}
