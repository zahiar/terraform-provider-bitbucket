package bitbucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketDeploymentVariable() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceBitbucketDeploymentVariableRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the deployment variable.",
				Type:        schema.TypeString,
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
			"deployment": {
				Description: "The UUID (including the enclosing `{}`) of the deployment environment.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"key": {
				Description:      "The name of the variable (must consist of only ASCII letters, numbers, underscores & not begin with a number).",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateDeploymentVariableName,
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
