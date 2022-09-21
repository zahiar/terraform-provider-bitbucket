package bitbucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceBitbucketProjectRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The UUID of the project.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace": {
				Description: "The slug or UUID (including the enclosing `{}`) of the workspace the project belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of the project.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"key": {
				Description:      "The key of the project.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateProjectKey,
			},
			"description": {
				Description: "The description of the project.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_private": {
				Description: "A boolean to state if the project is private or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"has_wiki": {
				Description: "A boolean to state if the project includes a wiki or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}
