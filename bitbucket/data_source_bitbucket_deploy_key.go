package bitbucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketDeployKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBitbucketDeployKeyRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the deploy key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"workspace": {
				Description: "The slug or UUID (including the enclosing `{}`) of the workspace.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"repository": {
				Description:      "The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateRepositoryName,
			},
			"label": {
				Description: "The label of the deploy key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"key": {
				Description: "The public SSH key to attach to this deploy key.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func dataSourceBitbucketDeployKeyRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceBitbucketDeployKeyRead(ctx, resourceData, meta)
}
