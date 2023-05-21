package bitbucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketDeployment() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceBitbucketDeploymentReadByNameOrId,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the deployment.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Description: "The name of the deployment environment.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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
			"environment": {
				Description: "The environment of the deployment (will be one of 'Test', 'Staging', or 'Production').",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
