package bitbucket

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketWorkspace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBitbucketWorkspaceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The slug of the workspace.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description: "The type of the workspace.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"uuid": {
				Description: "The UUID of the workspace.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_private": {
				Description: "A boolean to state if the project is private or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"name": {
				Description: "The name of the workspace.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceBitbucketWorkspaceRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	workspace, err := client.Workspaces.Get(resourceData.Get("id").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get workspace with error: %s", err))
	}

	_ = resourceData.Set("type", workspace.Type)
	_ = resourceData.Set("uuid", workspace.UUID)
	_ = resourceData.Set("is_private", workspace.Is_Private)
	_ = resourceData.Set("name", workspace.Name)
	resourceData.SetId(workspace.Slug)

	return nil
}
