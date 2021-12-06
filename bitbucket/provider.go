package bitbucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gobb "github.com/ktrysmt/go-bitbucket"

	v1 "github.com/zahiar/terraform-provider-bitbucket/bitbucket/api/v1"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_USERNAME", nil),
				Description: "Username to authenticate with Bitbucket.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_PASSWORD", nil),
				Description: "Password to authenticate with Bitbucket.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"bitbucket_branch_restriction": dataSourceBitbucketBranchRestriction(),
			"bitbucket_default_reviewer":   dataSourceBitbucketDefaultReviewer(),
			"bitbucket_deploy_key":         dataSourceBitbucketDeployKey(),
			"bitbucket_group":              dataSourceBitbucketGroup(),
			"bitbucket_group_permission":   dataSourceBitbucketGroupPermission(),
			"bitbucket_project":            dataSourceBitbucketProject(),
			"bitbucket_repository":         dataSourceBitbucketRepository(),
			"bitbucket_user":               dataSourceBitbucketUser(),
			"bitbucket_webhook":            dataSourceBitbucketWebhook(),
			"bitbucket_workspace":          dataSourceBitbucketWorkspace(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"bitbucket_branch_restriction": resourceBitbucketBranchRestriction(),
			"bitbucket_default_reviewer":   resourceBitbucketDefaultReviewer(),
			"bitbucket_deploy_key":         resourceBitbucketDeployKey(),
			"bitbucket_group":              resourceBitbucketGroup(),
			"bitbucket_group_member":       resourceBitbucketGroupMember(),
			"bitbucket_group_permission":   resourceBitbucketGroupPermission(),
			"bitbucket_project":            resourceBitbucketProject(),
			"bitbucket_repository":         resourceBitbucketRepository(),
			"bitbucket_webhook":            resourceBitbucketWebhook(),
		},

		ConfigureContextFunc: configureProvider,
	}
}

type Clients struct {
	V1 *v1.Client
	V2 *gobb.Client
}

func configureProvider(ctx context.Context, resourceData *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client := gobb.NewBasicAuth(
		resourceData.Get("username").(string),
		resourceData.Get("password").(string),
	)

	v1Client := v1.NewClient(
		&v1.Auth{
			Username: resourceData.Get("username").(string),
			Password: resourceData.Get("password").(string),
		},
	)

	clients := &Clients{
		V1: v1Client,
		V2: client,
	}

	return clients, nil
}
