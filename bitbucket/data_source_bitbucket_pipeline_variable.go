package bitbucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketPipelineVariable() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceBitbucketPipelineVariableRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the pipeline variable.",
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
			"key": {
				Description: "The name of the variable.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"value": {
				Description: "The value of the variable (note: if this variable is marked 'secured', this attribute will be blank).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"secured": {
				Description: "Whether this variable is considered secure/sensitive. If true, then it's value will not be exposed in any logs or API requests.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}
