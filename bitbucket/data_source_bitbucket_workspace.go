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
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_private": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBitBucketWorkspaceRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

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
