package bitbucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketGroupPermission() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceBitbucketGroupPermissionRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the group permission.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace": {
				Description: "The UUID (including the enclosing `{}`) of the workspace.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"repository": {
				Description:      "The slug of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateRepositoryName,
			},
			"group": {
				Description: "The slug of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"permission": {
				Description: "The permission this group will have. Must be one of 'read', 'write', 'admin'.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
