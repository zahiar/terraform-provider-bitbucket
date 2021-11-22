package bitbucket

import (
	"context"
	"encoding/json"
	"errors"
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

	userResponse, err := client.Users.Get(resourceData.Get("id").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get user with error: %s", err))
	}

	user, err := decodeUserResponse(userResponse)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to decode user response with error: %s", err))
	}

	_ = resourceData.Set("nickname", user.Nickname)
	_ = resourceData.Set("display_name", user.DisplayName)
	_ = resourceData.Set("account_id", user.AccountId)
	_ = resourceData.Set("account_status", user.AccountStatus)
	resourceData.SetId(user.Uuid)

	return nil
}

type User struct {
	Uuid          string
	DisplayName   string `json:"display_name"`
	Nickname      string
	AccountId     string `json:"account_id"`
	AccountStatus string `json:"account_status"`
}

func decodeUserResponse(response interface{}) (*User, error) {
	userMap := response.(map[string]interface{})

	if userMap["type"] == "error" {
		return nil, errors.New("unable able to decode user API response")
	}

	user := &User{}
	jsonString, err := json.Marshal(userMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonString, &user)

	return user, err
}
