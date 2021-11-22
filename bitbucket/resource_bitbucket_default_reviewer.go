package bitbucket

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	gobb "github.com/ktrysmt/go-bitbucket"
)

func resourceBitbucketDefaultReviewer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBitbucketDefaultReviewerCreate,
		ReadContext:   resourceBitbucketDefaultReviewerRead,
		DeleteContext: resourceBitbucketDefaultReviewerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBitbucketDefaultReviewerImport,
		},
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
				ForceNew:    true,
			},
			"repository": {
				Description: "The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"user": {
				Description: "The user's UUID.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceBitbucketDefaultReviewerCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	workspace := resourceData.Get("workspace").(string)
	repository := resourceData.Get("repository").(string)
	user := resourceData.Get("user").(string)

	_, err := client.Repositories.Repository.AddDefaultReviewer(
		&gobb.RepositoryDefaultReviewerOptions{
			Owner:    workspace,
			RepoSlug: repository,
			Username: user,
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to add default reviewer to repository with error: %s", err))
	}

	resourceData.SetId(generateDefaultReviewerResourceId(workspace, repository, user))

	return resourceBitbucketDefaultReviewerRead(ctx, resourceData, meta)
}

func resourceBitbucketDefaultReviewerRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

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
		// Handle a special case whereby if this returns this exact error, then we know it doesn't exist, so needs
		// re-creating rather than erroring out completely. This allows Terraform to re-create it/show that it needs
		// re-creating, and allows the user to decide how they wish to proceed.
		if strings.Contains(fmt.Sprint(err), "unable to get default reviewer: 404 Not Found") {
			_ = resourceData.Set("workspace", nil)
			_ = resourceData.Set("repository", nil)
			_ = resourceData.Set("user", nil)
		} else {
			return diag.FromErr(fmt.Errorf("unable to get default reviewer for repository with error: %s", err))
		}
	}

	return nil
}

func resourceBitbucketDefaultReviewerDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	workspace := resourceData.Get("workspace").(string)
	repository := resourceData.Get("repository").(string)
	user := resourceData.Get("user").(string)

	_, err := client.Repositories.Repository.DeleteDefaultReviewer(
		&gobb.RepositoryDefaultReviewerOptions{
			Owner:    workspace,
			RepoSlug: repository,
			Username: user,
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete default reviewer for repository with error: %s", err))
	}

	resourceData.SetId("")

	return nil
}

func resourceBitbucketDefaultReviewerImport(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{resourceData}

	splitID := strings.Split(resourceData.Id(), "/")
	if len(splitID) < 3 {
		return ret, fmt.Errorf("invalid import ID. It must to be in this format \"<workspace-slug|workspace-uuid>/<repository-slug|repository-uuid>/<user-uuid>\"")
	}

	workspace := splitID[0]
	repository := splitID[1]
	user := splitID[2]

	_ = resourceData.Set("workspace", workspace)
	_ = resourceData.Set("repository", repository)
	_ = resourceData.Set("user", user)

	client := meta.(*Clients).V2
	_, err := client.Repositories.Repository.GetDefaultReviewer(
		&gobb.RepositoryDefaultReviewerOptions{
			Owner:    workspace,
			RepoSlug: repository,
			Username: user,
		},
	)
	if err != nil {
		return ret, fmt.Errorf("unable to import default reviewer for repository with error: %s", err)
	}

	resourceData.SetId(generateDefaultReviewerResourceId(workspace, repository, user))

	return ret, nil
}

func generateDefaultReviewerResourceId(workspace string, repository string, user string) string {
	return fmt.Sprintf("%s-%s-%s", workspace, repository, user)
}
