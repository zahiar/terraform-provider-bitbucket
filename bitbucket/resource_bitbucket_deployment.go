package bitbucket

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gobb "github.com/ktrysmt/go-bitbucket"
)

func resourceBitbucketDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBitbucketDeploymentCreate,
		ReadContext:   resourceBitbucketDeploymentRead,
		DeleteContext: resourceBitbucketDeploymentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBitbucketDeploymentImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the deployment.",
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
				Description:      "The name of the repository (must consist of only lowercase ASCII letters, numbers, underscores and hyphens).",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateRepositoryName,
			},
			"name": {
				Description: "The name of the deployment environment.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"environment": {
				Description:  "The environment of the deployment (must be either 'Test', 'Staging' or 'Production').",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{gobb.Test.String(), gobb.Staging.String(), gobb.Production.String()}, false),
			},
		},
	}
}

func resourceBitbucketDeploymentCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	deploymentEnvironment, err := getDeploymentEnvironmentIntValue(resourceData.Get("environment").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	deployment, err := client.Repositories.Repository.AddEnvironment(
		&gobb.RepositoryEnvironmentOptions{
			Owner:           resourceData.Get("workspace").(string),
			RepoSlug:        resourceData.Get("repository").(string),
			Name:            resourceData.Get("name").(string),
			EnvironmentType: deploymentEnvironment,
			Rank:            int(deploymentEnvironment),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to create deployment enviroment with error: %s", err))
	}

	resourceData.SetId(deployment.Uuid)

	return resourceBitbucketDeploymentRead(ctx, resourceData, meta)
}

func resourceBitbucketDeploymentRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	deployment, err := client.Repositories.Repository.GetEnvironment(
		&gobb.RepositoryEnvironmentOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
			Uuid:     resourceData.Get("id").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get deployment environment with error: %s", err))
	}

	_ = resourceData.Set("name", deployment.Name)
	_ = resourceData.Set("environment", gobb.RepositoryEnvironmentTypeOption(deployment.Rank).String())
	resourceData.SetId(deployment.Uuid)

	return nil
}

func resourceBitbucketDeploymentReadByNameOrId(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := resourceData.Get("id")
	if id != nil && id != "" {
		return resourceBitbucketDeploymentRead(ctx, resourceData, meta)
	}

	name := resourceData.Get("name")
	if name != nil && name != "" {
		return resourceBitbucketDeploymentReadByName(ctx, resourceData, meta)
	}

	return diag.Errorf("Either name or id must be provided")
}

func resourceBitbucketDeploymentReadByName(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Artificial sleep due to Bitbucket's API taking time to return newly created deployments
	time.Sleep(3 * time.Second)

	client := meta.(*Clients).V2

	deployments, err := client.Repositories.Repository.ListEnvironments(
		&gobb.RepositoryEnvironmentsOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get deployment variable with error: %s", err))
	}

	name := resourceData.Get("name").(string)

	var deployment *gobb.Environment
	for _, item := range deployments.Environments {
		if strings.EqualFold(item.Name, name) {
			deployment = &item
			break
		}
	}

	if deployment == nil {
		return diag.FromErr(errors.New("unable to get deployment, Bitbucket API did not return it"))
	}

	_ = resourceData.Set("name", deployment.Name)
	_ = resourceData.Set("environment", gobb.RepositoryEnvironmentTypeOption(deployment.Rank).String())
	resourceData.SetId(deployment.Uuid)

	return nil
}

func resourceBitbucketDeploymentDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	_, err := client.Repositories.Repository.DeleteEnvironment(
		&gobb.RepositoryEnvironmentDeleteOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
			Uuid:     resourceData.Get("id").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete deployment environment with error: %s", err))
	}

	resourceData.SetId("")

	return nil
}

func resourceBitbucketDeploymentImport(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{resourceData}

	splitID := strings.Split(resourceData.Id(), "/")
	if len(splitID) < 3 {
		return ret, fmt.Errorf("invalid import ID. It must to be in this format \"<workspace-slug|workspace-uuid>/<repository-slug|repository-uuid>/<deployment-uuid>\"")
	}

	_ = resourceData.Set("workspace", splitID[0])
	_ = resourceData.Set("repository", splitID[1])
	_ = resourceData.Set("id", splitID[2])

	_ = resourceBitbucketDeploymentRead(ctx, resourceData, meta)

	return ret, nil
}

func getDeploymentEnvironmentIntValue(environment string) (gobb.RepositoryEnvironmentTypeOption, error) {
	switch environment {
	case "Test":
		return gobb.Test, nil
	case "Staging":
		return gobb.Staging, nil
	case "Production":
		return gobb.Production, nil
	}

	return -1, errors.New("invalid deployment environment")
}
