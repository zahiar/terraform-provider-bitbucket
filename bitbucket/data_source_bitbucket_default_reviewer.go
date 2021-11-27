package bitbucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketDefaultReviewer() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceBitbucketDefaultReviewerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the default reviewer.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace": {
				Description: "The slug or UUID (including the enclosing `{}`) of the workspace.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"repository": {
				Description: "The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"user": {
				Description: "The user's UUID.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}
