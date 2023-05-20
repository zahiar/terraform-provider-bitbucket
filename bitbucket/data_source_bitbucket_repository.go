package bitbucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketRepository() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceBitbucketRepositoryRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The UUID of the repository.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace": {
				Description: "The slug or UUID (including the enclosing `{}`) of the workspace this repository belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description:      "The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateRepositoryName,
			},
			"project_key": {
				Description: "The key of the project this repository belongs to.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "The description of the repository.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_private": {
				Description: "A boolean to state if the repository is private or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"has_wiki": {
				Description: "A boolean to state if the repository includes a wiki or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"fork_policy": {
				Description: "The name of the fork policy to apply to this repository.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"enable_pipelines": {
				Description: "A boolean to state if pipelines have been enabled for this repository.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"language": {
				Description: "The language of the repository.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
