package bitbucket

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gobb "github.com/ktrysmt/go-bitbucket"
)

func dataSourceBitbucketDefaultReviewer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBitbucketDefaultReviewerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the default reviewer.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace": {
				Description: "The slug or UUID (including the enclosing `{}`) of the workspace.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"repository": {
				Description: "The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"user": {
				Description: "The user's UUID.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceBitbucketDefaultReviewerRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	workspace := resourceData.Get("workspace").(string)
	repository := resourceData.Get("repository").(string)
	user := resourceData.Get("user").(string)

	_, err := client.Repositories.Repository.GetDefaultReviewer(
		&gobb.RepositoryDefaultReviewerOptions{
			Owner:    workspace,
			RepoSlug: repository,
			Username: user,
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get default reviewer for repository with error: %s", err))
	}

	resourceData.SetId(generateDefaultReviewerResourceId(workspace, repository, user))

	return nil
}
