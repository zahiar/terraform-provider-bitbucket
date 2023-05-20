package bitbucket

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketWorkspaceMembers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBitbucketWorkspaceMembersRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Description: "The slug of the workspace.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"members": {
				Description: "List of member UUID's (including the enclosing `{}`).",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The Member's UUID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"nickname": {
							Description: "The Member's nickname.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func dataSourceBitbucketWorkspaceMembersRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	membership, err := client.Workspaces.Members(resourceData.Get("workspace").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get workspace members with error: %s", err))
	}

	var members []interface{}
	for _, member := range membership.Members {
		members = append(members, map[string]interface{}{
			"id":       member.Uuid,
			"nickname": member.Nickname,
		})
	}
	_ = resourceData.Set("members", members)
	resourceData.SetId(resourceData.Get("workspace").(string))

	return nil
}
