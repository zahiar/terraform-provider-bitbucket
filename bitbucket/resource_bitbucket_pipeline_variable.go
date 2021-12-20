package bitbucket

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gobb "github.com/ktrysmt/go-bitbucket"
)

func resourceBitbucketPipelineVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBitbucketPipelineVariableCreate,
		ReadContext:   resourceBitbucketPipelineVariableRead,
		UpdateContext: resourceBitbucketPipelineVariableUpdate,
		DeleteContext: resourceBitbucketPipelineVariableDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBitbucketPipelineVariableImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the pipeline variable.",
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
			"key": {
				Description:      "The name of the variable (must consist of only ASCII letters, numbers, underscores & not begin with a number).",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateRepositoryVariableName,
			},
			"value": {
				Description: "The value of the variable.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"secured": {
				Description: "Whether this variable is considered secure/sensitive. If true, then it's value will not be exposed in any logs or API requests.",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
		},
	}
}

func resourceBitbucketPipelineVariableCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	pipelineVariable, err := client.Repositories.Repository.AddPipelineVariable(
		&gobb.RepositoryPipelineVariableOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
			Key:      resourceData.Get("key").(string),
			Value:    resourceData.Get("value").(string),
			Secured:  resourceData.Get("secured").(bool),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to create pipeline variable with error: %s", err))
	}

	resourceData.SetId(pipelineVariable.Uuid)

	return resourceBitbucketPipelineVariableRead(ctx, resourceData, meta)
}

func resourceBitbucketPipelineVariableRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	pipelineVariable, err := client.Repositories.Repository.GetPipelineVariable(
		&gobb.RepositoryPipelineVariableOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
			Uuid:     resourceData.Get("id").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get pipeline variable with error: %s", err))
	}

	_ = resourceData.Set("key", pipelineVariable.Key)

	if !pipelineVariable.Secured {
		_ = resourceData.Set("value", pipelineVariable.Value)
	} else {
		_ = resourceData.Set("value", resourceData.Get("value").(string))
	}

	_ = resourceData.Set("secured", pipelineVariable.Secured)

	resourceData.SetId(pipelineVariable.Uuid)

	return nil
}

func resourceBitbucketPipelineVariableUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	_, err := client.Repositories.Repository.UpdatePipelineVariable(
		&gobb.RepositoryPipelineVariableOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
			Uuid:     resourceData.Get("id").(string),
			Key:      resourceData.Get("key").(string),
			Value:    resourceData.Get("value").(string),
			Secured:  resourceData.Get("secured").(bool),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to update pipeline variable with error: %s", err))
	}

	return resourceBitbucketPipelineVariableRead(ctx, resourceData, meta)
}

func resourceBitbucketPipelineVariableDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	_, err := client.Repositories.Repository.DeletePipelineVariable(
		&gobb.RepositoryPipelineVariableDeleteOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
			Uuid:     resourceData.Id(),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete pipeline variable with error: %s", err))
	}

	resourceData.SetId("")

	return nil
}

func resourceBitbucketPipelineVariableImport(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{resourceData}

	splitID := strings.Split(resourceData.Id(), "/")
	if len(splitID) < 3 {
		return ret, fmt.Errorf("invalid import ID. It must to be in this format \"<workspace-slug|workspace-uuid>/<repository-slug|repository-uuid>/<pipeline-variable-uuid>\"")
	}

	_ = resourceData.Set("workspace", splitID[0])
	_ = resourceData.Set("repository", splitID[1])
	_ = resourceData.Set("id", splitID[2])

	_ = resourceBitbucketPipelineVariableRead(ctx, resourceData, meta)

	return ret, nil
}

func validateRepositoryVariableName(val interface{}, path cty.Path) diag.Diagnostics {
	match, _ := regexp.MatchString("^([a-zA-Z])[a-zA-Z0-9_]+$", val.(string))
	if !match {
		return diag.FromErr(fmt.Errorf("variable name must consist of only ASCII letters, numbers, underscores & not begin with a number (a-z, 0-9, _)"))
	}

	return diag.Diagnostics{}
}
