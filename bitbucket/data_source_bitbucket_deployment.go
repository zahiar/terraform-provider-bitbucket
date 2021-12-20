package bitbucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketDeployment() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceBitbucketDeploymentRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the deployment.",
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
			"name": {
				Description: "The name of the deployment environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"environment": {
				Description: "The environment of the deployment (will be one of 'Test', 'Staging', or 'Production').",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
