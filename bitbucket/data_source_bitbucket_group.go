package bitbucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceBitbucketGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace": {
				Description: "The UUID (including the enclosing `{}`) of the workspace this group belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of the group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"slug": {
				Description: "The slug of the group (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"auto_add": {
				Description: "Whether this group is auto-added to all future repositories.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"permission": {
				Description: "The permission this group will have over repositories. Must be one of 'read', 'write', 'admin'.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
