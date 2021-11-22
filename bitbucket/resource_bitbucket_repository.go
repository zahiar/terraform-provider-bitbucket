package bitbucket

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gobb "github.com/ktrysmt/go-bitbucket"
)

func resourceBitbucketRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBitbucketRepositoryCreate,
		ReadContext:   resourceBitbucketRepositoryRead,
		UpdateContext: resourceBitbucketRepositoryUpdate,
		DeleteContext: resourceBitbucketRepositoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBitbucketRepositoryImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The UUID of the repository.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace": {
				Description: "The slug or UUID (including the enclosing `{}`) of the workspace this repository belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description:      "The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateRepositoryName,
			},
			"project_key": {
				Description:      "The key of the project this repository belongs to.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateProjectKey,
			},
			"description": {
				Description: "The description of the repository.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"is_private": {
				Description: "A boolean to state if the repository is private or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"fork_policy": {
				Description:  "The name of the fork policy to apply to this repository.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "no_forks",
				ValidateFunc: validation.StringInSlice([]string{"allow_forks", "no_public_forks", "no_forks"}, false),
				DiffSuppressFunc: func(k, old, new string, resourceData *schema.ResourceData) bool {
					return !resourceData.Get("is_private").(bool)
				},
			},
		},
	}
}

func resourceBitbucketRepositoryCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	repository, err := client.Repositories.Repository.Create(
		&gobb.RepositoryOptions{
			Owner:       resourceData.Get("workspace").(string),
			RepoSlug:    resourceData.Get("name").(string),
			Description: resourceData.Get("description").(string),
			Project:     resourceData.Get("project_key").(string),
			IsPrivate:   strconv.FormatBool(resourceData.Get("is_private").(bool)),
			ForkPolicy:  resourceData.Get("fork_policy").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to create repository with error: %s", err))
	}

	resourceData.SetId(repository.Uuid)

	return resourceBitbucketRepositoryRead(ctx, resourceData, meta)
}

func resourceBitbucketRepositoryRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	repository, err := client.Repositories.Repository.Get(
		&gobb.RepositoryOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("name").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get repository with error: %s", err))
	}

	_ = resourceData.Set("description", repository.Description)
	_ = resourceData.Set("project_key", repository.Project.Key)
	_ = resourceData.Set("is_private", repository.Is_private)
	_ = resourceData.Set("fork_policy", repository.Fork_policy)

	resourceData.SetId(repository.Uuid)

	return nil
}

func resourceBitbucketRepositoryUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	_, err := client.Repositories.Repository.Update(
		&gobb.RepositoryOptions{
			Uuid:        resourceData.Id(),
			Owner:       resourceData.Get("workspace").(string),
			RepoSlug:    resourceData.Get("name").(string),
			Description: resourceData.Get("description").(string),
			Project:     resourceData.Get("project_key").(string),
			IsPrivate:   strconv.FormatBool(resourceData.Get("is_private").(bool)),
			ForkPolicy:  resourceData.Get("fork_policy").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to update repository with error: %s", err))
	}

	return resourceBitbucketRepositoryRead(ctx, resourceData, meta)
}

func resourceBitbucketRepositoryDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	_, err := client.Repositories.Repository.Delete(
		&gobb.RepositoryOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("name").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete repository with error: %s", err))
	}

	resourceData.SetId("")

	return nil
}

func resourceBitbucketRepositoryImport(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{resourceData}

	splitID := strings.Split(resourceData.Id(), "/")
	if len(splitID) < 2 {
		return ret, fmt.Errorf("invalid import ID. It must to be in this format \"<workspace-slug|workspace-uuid>/<repository-name>\"")
	}

	_ = resourceData.Set("workspace", splitID[0])
	_ = resourceData.Set("name", splitID[1])

	_ = resourceBitbucketRepositoryRead(ctx, resourceData, meta)

	return ret, nil
}

func validateRepositoryName(val interface{}, path cty.Path) diag.Diagnostics {
	match, _ := regexp.MatchString("^[a-z0-9_-]+$", val.(string))
	if !match {
		return diag.FromErr(fmt.Errorf("repository name must only consist of lowercase ASCII letters, numbers, underscores & hyphens (a-z, 0-9, _, -)"))
	}

	return diag.Diagnostics{}
}
