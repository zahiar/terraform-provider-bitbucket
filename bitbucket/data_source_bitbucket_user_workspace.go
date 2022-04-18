package bitbucket

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketUserWorkspace() *schema.Resource {
	// Any changes made here must be made to `dataSourceBitbucketUser`
	return &schema.Resource{
		ReadContext: dataSourceBitbucketUserWorkspaceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The User's UUID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace": {
				Description: "The slug or UUID (including the enclosing `{}`) of the workspace this user belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"nickname": {
				Description: "The User's nickname.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"display_name": {
				Description: "The User's display name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"account_id": {
				Description: "The User's account ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"account_status": {
				Description: "The User's account status. Will be one of active, inactive or closed.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceBitbucketUserWorkspaceRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	workspaceUsers, err := client.Workspaces.Members(resourceData.Get("workspace").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get users with error: %s", err))
	}

	nickname := resourceData.Get("nickname").(string)
	for _, user := range workspaceUsers.Members {
		if nickname == user.Nickname {
			resourceData.SetId(user.Uuid)
			return dataSourceBitbucketUserRead(ctx, resourceData, meta)
		}
	}

	return diag.FromErr(fmt.Errorf("unable to find user"))
}
