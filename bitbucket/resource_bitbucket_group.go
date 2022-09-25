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

func resourceBitbucketGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBitbucketGroupCreate,
		ReadContext:   resourceBitbucketGroupRead,
		UpdateContext: resourceBitbucketGroupUpdate,
		DeleteContext: resourceBitbucketGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBitbucketGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace": {
				Description: "The UUID (including the enclosing `{}`) of the workspace this group belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"slug": {
				Description: "The slug of the group (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"permission": {
				Description:  "The global permission this group will have over all repositories. Must be one of 'none', 'read', 'write', 'admin'.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "none",
				ValidateFunc: validation.StringInSlice([]string{"none", "read", "write", "admin"}, false),
			},
		},
	}
}

func resourceBitbucketGroupCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V1

	group, err := client.Groups.Create(
		&v1.GroupOptions{
			OwnerUuid: resourceData.Get("workspace").(string),
			Name:      resourceData.Get("name").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to create group with error: %s", err))
	}

	_ = resourceData.Set("slug", group.Slug)

	// We do an update as well, as that's where the other options like permissions are set.
	// The POST endpoint only accepts a name, and just creates a group without setting any other options.
	return resourceBitbucketGroupUpdate(ctx, resourceData, meta)
}

func resourceBitbucketGroupRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V1

	group, err := client.Groups.Get(
		&v1.GroupOptions{
			OwnerUuid: resourceData.Get("workspace").(string),
			Slug:      resourceData.Get("slug").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get group with error: %s", err))
	}

	_ = resourceData.Set("name", group.Name)
	_ = resourceData.Set("permission", group.Permission)

	resourceData.SetId(generateGroupResourceId(group.Owner.Uuid, group.Slug))

	return nil
}

func resourceBitbucketGroupUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V1

	group, err := client.Groups.Update(
		&v1.GroupOptions{
			OwnerUuid:  resourceData.Get("workspace").(string),
			Slug:       resourceData.Get("slug").(string),
			Name:       resourceData.Get("name").(string),
			Permission: resourceData.Get("permission").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to update group with error: %s", err))
	}

	_ = resourceData.Set("slug", group.Slug)

	return resourceBitbucketGroupRead(ctx, resourceData, meta)
}

func resourceBitbucketGroupDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V1

	err := client.Groups.Delete(
		&v1.GroupOptions{
			OwnerUuid: resourceData.Get("workspace").(string),
			Slug:      resourceData.Get("slug").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete group with error: %s", err))
	}

	resourceData.SetId("")

	return nil
}

func resourceBitbucketGroupImport(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{resourceData}

	splitID := strings.Split(resourceData.Id(), "/")
	if len(splitID) < 2 {
		return ret, fmt.Errorf("invalid import ID. It must to be in this format \"<workspace-uuid>/<group-slug>\"")
	}

	_ = resourceData.Set("workspace", splitID[0])
	_ = resourceData.Set("slug", splitID[1])

	_ = resourceBitbucketGroupRead(ctx, resourceData, meta)

	return ret, nil
}

func generateGroupResourceId(workspace string, slug string) string {
	return fmt.Sprintf("%s-%s", workspace, slug)
}
