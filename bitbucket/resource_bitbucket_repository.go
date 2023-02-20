package bitbucket

import (
	"context"
	"fmt"
	"reflect"
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
				Description:      "The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores, hyphens and periods).",
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
			"has_wiki": {
				Description: "A boolean to state if the repository includes a wiki or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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
			"enable_pipelines": {
				Description: "A boolean to state if pipelines have been enabled for this repository.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceBitbucketRepositoryCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	workspace := resourceData.Get("workspace").(string)

	repository, err := client.Repositories.Repository.Create(
		&gobb.RepositoryOptions{
			Owner:       workspace,
			RepoSlug:    resourceData.Get("name").(string),
			Description: resourceData.Get("description").(string),
			Project:     resourceData.Get("project_key").(string),
			IsPrivate:   strconv.FormatBool(resourceData.Get("is_private").(bool)),
			HasWiki:     strconv.FormatBool(resourceData.Get("has_wiki").(bool)),
			ForkPolicy:  resourceData.Get("fork_policy").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to create repository with error: %s", err))
	}

	resourceData.SetId(repository.Uuid)

	_, err = client.Repositories.Repository.UpdatePipelineConfig(
		&gobb.RepositoryPipelineOptions{
			Owner:    workspace,
			RepoSlug: resourceData.Get("name").(string),
			Enabled:  resourceData.Get("enable_pipelines").(bool),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to enable pipelines for repository with error: %s", err))
	}

	// If cache has been warmed up, we need to add entry to the list to ensure repo presence is known
	// during runtime
	if len(repositoriesCache[workspace]) != 0 {
		repositoriesCache[workspace][repository.Slug] = *repository
	}

	return resourceBitbucketRepositoryRead(ctx, resourceData, meta)
}

func resourceBitbucketRepositoryRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	workspace := resourceData.Get("workspace").(string)

	if len(repositoriesCache[workspace]) == 0 {
		if err := warmUpRepositoriesCacheInTheWorkspace(client, workspace); err != nil {
			return diag.FromErr(fmt.Errorf("unable to get repositories with error: %s", err))
		}
	}

	repository := repositoriesCache[workspace][resourceData.Get("name").(string)]
	if reflect.DeepEqual(repository, gobb.Repository{}) {
		return diag.FromErr(fmt.Errorf("repository does not exist or was removed without using provider"))
	}

	_ = resourceData.Set("description", repository.Description)
	_ = resourceData.Set("project_key", repository.Project.Key)
	_ = resourceData.Set("is_private", repository.Is_private)
	_ = resourceData.Set("has_wiki", repository.Has_wiki)
	_ = resourceData.Set("fork_policy", repository.Fork_policy)

	resourceData.SetId(repository.Uuid)

	repositoryPipelineConfig, err := client.Repositories.Repository.GetPipelineConfig(
		&gobb.RepositoryPipelineOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("name").(string),
		},
	)
	if err != nil {
		// Underlying go-bitbucket library only returns an error object, so this is our best way to check for a 404.
		// This specifically addresses an issue whereby if you import a Bitbucket repository that has never had its
		// pipelines enabled, Bitbucket's API returns a 404.
		if err.Error() != "unable to get pipeline config: 404 Not Found" {
			return diag.FromErr(fmt.Errorf("unable to get pipeline configuration for repository with error: %s", err))
		}

		_ = resourceData.Set("enable_pipelines", false)
	} else {
		_ = resourceData.Set("enable_pipelines", repositoryPipelineConfig.Enabled)
	}

	return nil
}

func resourceBitbucketRepositoryUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	workspace := resourceData.Get("workspace").(string)

	repository, err := client.Repositories.Repository.Update(
		&gobb.RepositoryOptions{
			Uuid:        resourceData.Id(),
			Owner:       workspace,
			RepoSlug:    resourceData.Get("name").(string),
			Description: resourceData.Get("description").(string),
			Project:     resourceData.Get("project_key").(string),
			IsPrivate:   strconv.FormatBool(resourceData.Get("is_private").(bool)),
			HasWiki:     strconv.FormatBool(resourceData.Get("has_wiki").(bool)),
			ForkPolicy:  resourceData.Get("fork_policy").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to update repository with error: %s", err))
	}

	_, err = client.Repositories.Repository.UpdatePipelineConfig(
		&gobb.RepositoryPipelineOptions{
			Owner:    workspace,
			RepoSlug: resourceData.Get("name").(string),
			Enabled:  resourceData.Get("enable_pipelines").(bool),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to update pipeline configuration for repository with error: %s", err))
	}

	// If cache has been warmed up, we need to update entry in the list to ensure repo current state is known
	// during runtime
	if len(repositoriesCache[workspace]) != 0 {
		repositoriesCache[workspace][repository.Slug] = *repository
	}

	return resourceBitbucketRepositoryRead(ctx, resourceData, meta)
}

func resourceBitbucketRepositoryDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	workspace := resourceData.Get("workspace").(string)
	repoSlug := resourceData.Get("name").(string)

	_, err := client.Repositories.Repository.Delete(
		&gobb.RepositoryOptions{
			Owner:    workspace,
			RepoSlug: repoSlug,
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete repository with error: %s", err))
	}

	resourceData.SetId("")

	// If cache has been warmed up, we need to remove entry from the list
	if len(repositoriesCache[workspace]) != 0 {
		delete(repositoriesCache[workspace], repoSlug)
	}

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
	match, _ := regexp.MatchString("^[a-z0-9\\._-]+$", val.(string))
	if !match {
		return diag.FromErr(fmt.Errorf("repository name must only consist of lowercase ASCII letters, numbers, underscores & hyphens (a-z, 0-9, _, -)"))
	}

	return diag.Diagnostics{}
}

var repositoriesCache = make(map[string]map[string]gobb.Repository)

func warmUpRepositoriesCacheInTheWorkspace(client *gobb.Client, workspace string) error {
	repositories, err := client.Repositories.ListForAccount(&gobb.RepositoriesOptions{Owner: workspace})
	if err != nil {
		return err
	}

	tempRepositories := make(map[string]gobb.Repository)
	for _, item := range repositories.Items {
		tempRepositories[item.Slug] = item
	}
	repositoriesCache[workspace] = tempRepositories

	return nil
}
