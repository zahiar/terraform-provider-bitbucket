package bitbucket

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	gobb "github.com/ktrysmt/go-bitbucket"
)

func dataSourceBitBucketWorkspace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBitBucketWorkspaceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The UUID of the workspace.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "The type of the workspace.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"slug": {
				Description: "The slug of the workspace.",
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
				Required:    true,
			},
		},
	}
}

func dataSourceBitBucketWorkspaceRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	workspace, err := client.Workspaces.Get(resourceData.Get("name").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get workspace with error: %s", err))
	}

	_ = resourceData.Set("type", workspace.Type)
	_ = resourceData.Set("slug", workspace.Slug)
	_ = resourceData.Set("is_private", workspace.Is_Private)
	_ = resourceData.Set("name", workspace.Name)

	resourceData.SetId(workspace.UUID)

	return nil
}
