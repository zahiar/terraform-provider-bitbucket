package bitbucket

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ktrysmt/go-bitbucket"
)

func resourceBitbucketUserPermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBitbucketUserPermissionCreate,
		ReadContext:   resourceBitbucketUserPermissionRead,
		DeleteContext: resourceBitbucketUserPermissionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBitbucketUserPermissionImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the user permission.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace": {
				Description: "The UUID (including the enclosing `{}`) of the workspace.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"repository": {
				Description:      "The slug of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateRepositoryName,
				ForceNew:         true,
			},
			"user": {
				Description: "The UUID (including the enclosing `{}`) of the user.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"permission": {
				Description:  "The permission this user will have. Must be one of 'read', 'write', 'admin'.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"read", "write", "admin"}, false),
				ForceNew:     true,
			},
		},
	}
}

func resourceBitbucketUserPermissionCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	_, err := client.Repositories.Repository.SetUserPermissions(&bitbucket.RepositoryUserPermissionsOptions{
		Owner:      resourceData.Get("workspace").(string),
		RepoSlug:   resourceData.Get("repository").(string),
		User:       resourceData.Get("user").(string),
		Permission: resourceData.Get("permission").(string),
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to create user permission with error: %s", err))
	}

	return resourceBitbucketUserPermissionRead(ctx, resourceData, meta)
}

func resourceBitbucketUserPermissionRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	workspace := resourceData.Get("workspace").(string)
	repository := resourceData.Get("repository").(string)

	userPermission, err := client.Repositories.Repository.GetUserPermissions(&bitbucket.RepositoryUserPermissionsOptions{
		Owner:      workspace,
		RepoSlug:   repository,
		User:       resourceData.Get("user").(string),
		Permission: resourceData.Get("permission").(string),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get user permission with error: %s", err))
	}

	_ = resourceData.Set("user", userPermission.User.Uuid)
	_ = resourceData.Set("permission", userPermission.Permission)

	resourceData.SetId(generateUserPermissionResourceId(workspace, repository, userPermission.User.Uuid))

	return nil
}

func resourceBitbucketUserPermissionDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	_, err := client.Repositories.Repository.DeleteUserPermissions(&bitbucket.RepositoryUserPermissionsOptions{
		Owner:      resourceData.Get("workspace").(string),
		RepoSlug:   resourceData.Get("repository").(string),
		User:       resourceData.Get("user").(string),
		Permission: resourceData.Get("permission").(string),
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete user permission with error: %s", err))
	}

	resourceData.SetId("")

	return nil
}

func resourceBitbucketUserPermissionImport(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{resourceData}

	splitID := strings.Split(resourceData.Id(), "/")
	if len(splitID) < 3 {
		return ret, fmt.Errorf("invalid import ID. It must to be in this format \"<workspace-uuid>/<repo-slug>/<user-uuid>\"")
	}

	_ = resourceData.Set("workspace", splitID[0])
	_ = resourceData.Set("repository", splitID[1])
	_ = resourceData.Set("user", splitID[2])

	_ = resourceBitbucketUserPermissionRead(ctx, resourceData, meta)

	return ret, nil
}

func generateUserPermissionResourceId(workspace string, repo string, user string) string {
	return fmt.Sprintf("%s-%s-%s", workspace, repo, user)
}
