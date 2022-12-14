package bitbucket

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	v1 "github.com/brnck/terraform-provider-bitbucket/bitbucket/api/v1"
)

func resourceBitbucketGroupMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBitbucketGroupMemberCreate,
		ReadContext:   resourceBitbucketGroupMemberRead,
		DeleteContext: resourceBitbucketGroupMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBitbucketGroupMemberImport,
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
			"group": {
				Description: "The slug of the group.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"user": {
				Description: "The User's UUID.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceBitbucketGroupMemberCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V1

	_, err := client.GroupMembers.Create(
		&v1.GroupMemberOptions{
			OwnerUuid: resourceData.Get("workspace").(string),
			Slug:      resourceData.Get("group").(string),
			UserUuid:  resourceData.Get("user").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to create group member with error: %s", err))
	}

	return resourceBitbucketGroupMemberRead(ctx, resourceData, meta)
}

func resourceBitbucketGroupMemberRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V1

	group, err := client.Groups.Get(
		&v1.GroupOptions{
			OwnerUuid: resourceData.Get("workspace").(string),
			Slug:      resourceData.Get("group").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get group with error: %s", err))
	}

	groupMembers, err := client.GroupMembers.Get(
		&v1.GroupMemberOptions{
			OwnerUuid: resourceData.Get("workspace").(string),
			Slug:      resourceData.Get("group").(string),
			UserUuid:  resourceData.Get("user").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get group member with error: %s", err))
	}

	for _, member := range groupMembers {
		if strings.EqualFold(resourceData.Get("user").(string), member.UUID) {
			resourceData.SetId(generateGroupMemberId(group.Owner.Uuid, group.Slug, member.UUID))
		}
	}

	return nil
}

func resourceBitbucketGroupMemberDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V1

	err := client.GroupMembers.Delete(
		&v1.GroupMemberOptions{
			OwnerUuid: resourceData.Get("workspace").(string),
			Slug:      resourceData.Get("group").(string),
			UserUuid:  resourceData.Get("user").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete group member with error: %s", err))
	}

	resourceData.SetId("")

	return nil
}

func resourceBitbucketGroupMemberImport(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{resourceData}

	splitID := strings.Split(resourceData.Id(), "/")
	if len(splitID) < 3 {
		return ret, fmt.Errorf("invalid import ID. It must to be in this format \"<workspace-uuid>/<group-slug>/<user-uuid>\"")
	}

	_ = resourceData.Set("workspace", splitID[0])
	_ = resourceData.Set("group", splitID[1])
	_ = resourceData.Set("user", splitID[2])

	_ = resourceBitbucketGroupMemberRead(ctx, resourceData, meta)

	return ret, nil
}

func generateGroupMemberId(workspace string, group string, user string) string {
	return fmt.Sprintf("%s-%s-%s", workspace, group, user)
}
