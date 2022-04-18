package bitbucket

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBitbucketUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The User's UUID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"nickname": {
				Description: "The User's nickname.",
				Type:        schema.TypeString,
				Computed:    true,
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

func dataSourceBitbucketUserRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	user, err := client.Users.Get(resourceData.Get("id").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get user with error: %s", err))
	}

	_ = resourceData.Set("nickname", user.Nickname)
	_ = resourceData.Set("display_name", user.DisplayName)
	_ = resourceData.Set("account_id", user.AccountId)
	_ = resourceData.Set("account_status", user.AccountStatus)
	resourceData.SetId(user.Uuid)

	return nil
}
