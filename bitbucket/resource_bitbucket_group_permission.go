package bitbucket

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	v1 "github.com/zahiar/terraform-provider-bitbucket/bitbucket/api/v1"
)

func resourceBitbucketGroupPermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBitbucketGroupPermissionCreate,
		ReadContext:   resourceBitbucketGroupPermissionRead,
		DeleteContext: resourceBitbucketGroupPermissionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBitbucketGroupPermissionImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the group permission.",
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
			"group": {
				Description: "The slug of the group.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"permission": {
				Description:  "The permission this group will haves. Must be one of 'read', 'write', 'admin'.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"read", "write", "admin"}, false),
				ForceNew:     true,
			},
		},
	}
}

func resourceBitbucketGroupPermissionCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V1

	_, err := client.GroupPrivileges.Create(
		&v1.GroupPrivilegeOptions{
			WorkspaceId: resourceData.Get("workspace").(string),
			RepoSlug:    resourceData.Get("repository").(string),
			GroupOwner:  resourceData.Get("workspace").(string),
			GroupSlug:   resourceData.Get("group").(string),
			Privilege:   resourceData.Get("permission").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to create group permission with error: %s", err))
	}

	return resourceBitbucketGroupPermissionRead(ctx, resourceData, meta)
}

func resourceBitbucketGroupPermissionRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V1

	groupPermission, err := client.GroupPrivileges.Get(
		&v1.GroupPrivilegeOptions{
			WorkspaceId: resourceData.Get("workspace").(string),
			RepoSlug:    resourceData.Get("repository").(string),
			GroupOwner:  resourceData.Get("workspace").(string),
			GroupSlug:   resourceData.Get("group").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get group permission with error: %s", err))
	}

	_ = resourceData.Set("workspace", groupPermission.Group.Owner.Uuid)
	_ = resourceData.Set("repository", groupPermission.Repository.Slug)
	_ = resourceData.Set("group", groupPermission.Group.Slug)
	_ = resourceData.Set("permission", groupPermission.Privilege)

	resourceData.SetId(generateGroupPermissionResourceId(groupPermission.Group.Owner.Uuid, groupPermission.Repository.Slug, groupPermission.Group.Slug))

	return nil
}

func resourceBitbucketGroupPermissionDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V1

	err := client.GroupPrivileges.Delete(
		&v1.GroupPrivilegeOptions{
			WorkspaceId: resourceData.Get("workspace").(string),
			RepoSlug:    resourceData.Get("repository").(string),
			GroupOwner:  resourceData.Get("workspace").(string),
			GroupSlug:   resourceData.Get("group").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete group permission with error: %s", err))
	}

	resourceData.SetId("")

	return nil
}

func resourceBitbucketGroupPermissionImport(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{resourceData}

	splitID := strings.Split(resourceData.Id(), "/")
	if len(splitID) < 3 {
		return ret, fmt.Errorf("invalid import ID. It must to be in this format \"<workspace-uuid>/<repo-slug>/<group-slug>\"")
	}

	_ = resourceData.Set("workspace", splitID[0])
	_ = resourceData.Set("repository", splitID[1])
	_ = resourceData.Set("group", splitID[2])

	_ = resourceBitbucketGroupPermissionRead(ctx, resourceData, meta)

	return ret, nil
}

func generateGroupPermissionResourceId(workspace string, repo string, group string) string {
	return fmt.Sprintf("%s-%s-%s", workspace, repo, group)
}
