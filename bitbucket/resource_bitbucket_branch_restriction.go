package bitbucket

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gobb "github.com/ktrysmt/go-bitbucket"
)

func resourceBitbucketBranchRestriction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBitbucketBranchRestrictionCreate,
		ReadContext:   resourceBitbucketBranchRestrictionRead,
		UpdateContext: resourceBitbucketBranchRestrictionUpdate,
		DeleteContext: resourceBitbucketBranchRestrictionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBitbucketBranchRestrictionImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the branch restriction.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace": {
				Description: "The slug or UUID (including the enclosing `{}`) of the workspace.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"repository": {
				Description:      "The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateRepositoryName,
			},
			"pattern": {
				Description: "The pattern to match against branches this restriction will apply to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"kind": {
				Description: "The type of restriction to apply.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"require_tasks_to_be_completed",
					"allow_auto_merge_when_builds_pass",
					"require_passing_builds_to_merge",
					"force",
					"require_all_dependencies_merged",
					"require_commits_behind",
					"restrict_merges",
					"enforce_merge_checks",
					"reset_pullrequest_changes_requested_on_change",
					"require_no_changes_requested",
					"smart_reset_pullrequest_approvals",
					"push",
					"require_approvals_to_merge",
					"require_default_reviewer_approvals_to_merge",
					"reset_pullrequest_approvals_on_change",
					"delete",
				}, false),
			},
			"value": {
				Description: "A configurable value used by the following restrictions: `require_passing_builds_to_merge` uses it to define the number of minimum number of passing builds, `require_approvals_to_merge` uses it to define the minimum number of approvals before the PR can be merged, `require_default_reviewer_approvals_to_merge` uses it to define the minimum number of approvals from default reviewers before the PR can be merged.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
			},
			"users": {
				Description: "A list of users (usernames or user's UUID) that are exempt from this branch restriction. Can only be set if restriction type (`kind`) is set to `push` or `restrict_merges`.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"groups": {
				Description: "A list of groups (group names or group UUID) that are exempt from this branch restriction. Can only be set if restriction type (`kind`) is set to `push` or `restrict_merges`.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceBitbucketBranchRestrictionCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	opts := &gobb.BranchRestrictionsOptions{
		Owner:    resourceData.Get("workspace").(string),
		RepoSlug: resourceData.Get("repository").(string),
		Pattern:  resourceData.Get("pattern").(string),
		Kind:     resourceData.Get("kind").(string),
		Value:    nil,
		Users:    nil,
		Groups:   nil,
	}
	value := resourceData.Get("value").(int)
	if value > 0 {
		opts.Value = value
	}
	users := parseBranchRestrictionUserFields(resourceData.Get("users").([]interface{}))
	if len(users) > 0 {
		opts.Users = users
	}
	groups := parseBranchRestrictionUserGroupFields(resourceData.Get("groups").([]interface{}))
	if len(groups) > 0 {
		opts.Groups = groups
	}

	branchRestriction, err := client.Repositories.BranchRestrictions.Create(opts)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to create branch restriction with error: %s", err))
	}

	resourceData.SetId(strconv.Itoa(branchRestriction.ID))

	return resourceBitbucketBranchRestrictionRead(ctx, resourceData, meta)
}

func resourceBitbucketBranchRestrictionRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	branchRestriction, err := client.Repositories.BranchRestrictions.Get(
		&gobb.BranchRestrictionsOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
			ID:       resourceData.Get("id").(string),
		},
	)
	if err != nil {
		// Handles a case whereby if the branch restrictions were deleted after being provisioned, Bitbucket's API
		// returns a 404, so we treat that as the item having been deleted, therefore Terraform will re-provision
		// if necessary.
		if err.Error() == "404 Not Found" {
			resourceData.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("unable to get branch restriction with error: %s", err))
	}

	_ = resourceData.Set("pattern", branchRestriction.Pattern)
	_ = resourceData.Set("kind", branchRestriction.Kind)
	_ = resourceData.Set("value", branchRestriction.Value)

	resourceData.SetId(strconv.Itoa(branchRestriction.ID))

	return nil
}

func resourceBitbucketBranchRestrictionUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	opts := &gobb.BranchRestrictionsOptions{
		Owner:    resourceData.Get("workspace").(string),
		RepoSlug: resourceData.Get("repository").(string),
		ID:       resourceData.Id(),
		Pattern:  resourceData.Get("pattern").(string),
		Kind:     resourceData.Get("kind").(string),
		Value:    nil,
		Users:    nil,
		Groups:   nil,
	}
	value := resourceData.Get("value").(int)
	if value > 0 {
		opts.Value = value
	}
	users := parseBranchRestrictionUserFields(resourceData.Get("users").([]interface{}))
	if len(users) > 0 {
		opts.Users = users
	}
	groups := parseBranchRestrictionUserGroupFields(resourceData.Get("groups").([]interface{}))
	if len(groups) > 0 {
		opts.Groups = groups
	}

	_, err := client.Repositories.BranchRestrictions.Update(opts)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to update branch restriction with error: %s", err))
	}

	return resourceBitbucketBranchRestrictionRead(ctx, resourceData, meta)
}

func resourceBitbucketBranchRestrictionDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	_, err := client.Repositories.BranchRestrictions.Delete(
		&gobb.BranchRestrictionsOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
			ID:       resourceData.Id(),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete branch restriction with error: %s", err))
	}

	resourceData.SetId("")

	return nil
}

func resourceBitbucketBranchRestrictionImport(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{resourceData}

	splitID := strings.Split(resourceData.Id(), "/")
	if len(splitID) < 2 {
		return ret, fmt.Errorf("invalid import ID. It must to be in this format \"<workspace-slug|workspace-uuid>/<repository-name>/<branch-restriction-id>\"")
	}

	_ = resourceData.Set("workspace", splitID[0])
	_ = resourceData.Set("repository", splitID[1])
	resourceData.SetId(splitID[2])

	_ = resourceBitbucketBranchRestrictionRead(ctx, resourceData, meta)

	return ret, nil
}

func parseBranchRestrictionUserFields(users []interface{}) []string {
	var usersArray []string
	for _, user := range users {
		usersArray = append(usersArray, user.(string))
	}

	return usersArray
}

func parseBranchRestrictionUserGroupFields(groups []interface{}) map[string]string {
	groupsMap := make(map[string]string)
	for _, group := range groups {
		group := group.(string)
		groupsMap[group] = group
	}

	return groupsMap
}
