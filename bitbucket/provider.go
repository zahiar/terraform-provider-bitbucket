package bitbucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gobb "github.com/ktrysmt/go-bitbucket"
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
			"bitbucket_project":    dataSourceBitbucketProject(),
			"bitbucket_repository": dataSourceBitbucketRepository(),
			"bitbucket_user":       dataSourceBitbucketUser(),
			"bitbucket_webhook":    dataSourceBitbucketWebhook(),
			"bitbucket_workspace":  dataSourceBitbucketWorkspace(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"bitbucket_project":    resourceBitbucketProject(),
			"bitbucket_repository": resourceBitbucketRepository(),
			"bitbucket_webhook":    resourceBitbucketWebhook(),
		},

		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, resourceData *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client := gobb.NewBasicAuth(
		resourceData.Get("username").(string),
		resourceData.Get("password").(string),
	)

	return client, nil
}
