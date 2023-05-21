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
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_USERNAME", nil),
				Description: "Username to authenticate with Bitbucket.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_PASSWORD", nil),
				Description: "Password to authenticate with Bitbucket.",
			},
			"oauth_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_OAUTH_CLIENT_ID", nil),
				Description: "Client ID for OAuth authentication with Bitbucket.",
			},
			"oauth_client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_OAUTH_CLIENT_SECRET", nil),
				Description: "Client secret for OAuth authentication with Bitbucket.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"bitbucket_branch_restriction":  dataSourceBitbucketBranchRestriction(),
			"bitbucket_default_reviewer":    dataSourceBitbucketDefaultReviewer(),
			"bitbucket_deploy_key":          dataSourceBitbucketDeployKey(),
			"bitbucket_deployment":          dataSourceBitbucketDeployment(),
			"bitbucket_deployment_variable": dataSourceBitbucketDeploymentVariable(),
			"bitbucket_group":               dataSourceBitbucketGroup(),
			"bitbucket_group_permission":    dataSourceBitbucketGroupPermission(),
			"bitbucket_ip_ranges":           dataSourceBitbucketIpRanges(),
			"bitbucket_pipeline_variable":   dataSourceBitbucketPipelineVariable(),
			"bitbucket_project":             dataSourceBitbucketProject(),
			"bitbucket_repository":          dataSourceBitbucketRepository(),
			"bitbucket_user":                dataSourceBitbucketUser(),
			"bitbucket_user_permission":     dataSourceBitbucketUserPermission(),
			"bitbucket_user_workspace":      dataSourceBitbucketUserWorkspace(),
			"bitbucket_webhook":             dataSourceBitbucketWebhook(),
			"bitbucket_workspace":           dataSourceBitbucketWorkspace(),
			"bitbucket_workspace_members":   dataSourceBitbucketWorkspaceMembers(),
			"bitbucket_workspace_projects":  dataSourceBitbucketWorkspaceProjects(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"bitbucket_branch_restriction":  resourceBitbucketBranchRestriction(),
			"bitbucket_default_reviewer":    resourceBitbucketDefaultReviewer(),
			"bitbucket_deploy_key":          resourceBitbucketDeployKey(),
			"bitbucket_deployment":          resourceBitbucketDeployment(),
			"bitbucket_deployment_variable": resourceBitbucketDeploymentVariable(),
			"bitbucket_group":               resourceBitbucketGroup(),
			"bitbucket_group_member":        resourceBitbucketGroupMember(),
			"bitbucket_group_permission":    resourceBitbucketGroupPermission(),
			"bitbucket_pipeline_key_pair":   resourceBitbucketPipelineKeyPair(),
			"bitbucket_pipeline_variable":   resourceBitbucketPipelineVariable(),
			"bitbucket_project":             resourceBitbucketProject(),
			"bitbucket_repository":          resourceBitbucketRepository(),
			"bitbucket_user_permission":     resourceBitbucketUserPermission(),
			"bitbucket_webhook":             resourceBitbucketWebhook(),
		},

		ConfigureContextFunc: configureProvider,
	}
}

type Clients struct {
	V1 *v1.Client
	V2 *gobb.Client
}

func configureProvider(ctx context.Context, resourceData *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := resourceData.Get("username").(string)
	password := resourceData.Get("password").(string)
	oauthClientId := resourceData.Get("oauth_client_id").(string)
	oauthClientSecret := resourceData.Get("oauth_client_secret").(string)

	var client *gobb.Client
	var v1Client *v1.Client

	if username != "" && password != "" {
		client = gobb.NewBasicAuth(username, password)
		v1Client = v1.NewBasicAuthClient(username, password)
	} else if oauthClientId != "" && oauthClientSecret != "" {
		client = gobb.NewOAuthClientCredentials(oauthClientId, oauthClientSecret)
		v1Client = v1.NewOAuthClient(oauthClientId, oauthClientSecret)
	} else if username != "" && password == "" {
		diag.Errorf("`username` is set but `password` is not.")
	} else if username == "" && password != "" {
		diag.Errorf("`password` is set but `username` is not.")
	} else if oauthClientId != "" && oauthClientSecret == "" {
		diag.Errorf("`oauth_client_id` is set but `oauth_client_secret` is not.")
	} else if oauthClientId == "" && oauthClientSecret != "" {
		diag.Errorf("`oauth_client_secret` is set but `oauth_client_id` is not.")
	} else {
		diag.Errorf("Either `username` and `password` or `oauth_client_id` and `oauth_client_secret` need to be set for acceptance tests")
	}

	client.Pagelen = 100
	client.MaxDepth = 10

	clients := &Clients{
		V1: v1Client,
		V2: client,
	}

	return clients, nil
}
