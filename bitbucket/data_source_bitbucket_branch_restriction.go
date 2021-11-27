package bitbucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketBranchRestriction() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceBitbucketBranchRestrictionRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the branch restriction.",
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
			"pattern": {
				Description: "The pattern to match against branches this restriction will apply to.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"kind": {
				Description: "The type of restriction to apply.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"value": {
				Description: "A configurable value used by the following restrictions: `require_passing_builds_to_merge` uses it to define the number of minimum number of passing builds, `require_approvals_to_merge` uses it to define the minimum number of approvals before the PR can be merged, `require_default_reviewer_approvals_to_merge` uses it to define the minimum number of approvals from default reviewers before the PR can be merged.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}
