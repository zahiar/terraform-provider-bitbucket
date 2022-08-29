package bitbucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketUserPermission() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceBitbucketUserPermissionRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the user permission.",
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
			"user": {
				Description: "The UUID (including the enclosing `{}`) of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"permission": {
				Description: "The permission this user will have. Must be one of 'read', 'write', 'admin'.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
