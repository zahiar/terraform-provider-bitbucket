package bitbucket

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBitbucketWorkspaceProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBitbucketWorkspaceProjectsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Description: "The slug of the workspace.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"projects": {
				Description: "List of Projects.",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The Project's UUID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "The Project's name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"key": {
							Description: "The Project's key.",
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

func dataSourceBitbucketWorkspaceProjectsRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	membership, err := client.Workspaces.Projects(resourceData.Get("workspace").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get workspace projects with error: %s", err))
	}

	var projects []interface{}
	for _, project := range membership.Items {
		projects = append(projects, map[string]interface{}{
			"id":   project.Uuid,
			"name": project.Name,
			"key":  project.Key,
		})
	}
	_ = resourceData.Set("projects", projects)
	resourceData.SetId(resourceData.Get("workspace").(string))

	return nil
}
