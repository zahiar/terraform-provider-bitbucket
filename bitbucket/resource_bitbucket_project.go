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

func resourceBitbucketProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBitbucketProjectCreate,
		ReadContext:   resourceBitbucketProjectRead,
		UpdateContext: resourceBitbucketProjectUpdate,
		DeleteContext: resourceBitbucketProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBitbucketProjectImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The UUID of the project.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace": {
				Description: "The slug or UUID (including the enclosing `{}`) of the workspace the project belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of the project.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"key": {
				Description:      "The key of the project.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateProjectKey,
			},
			"description": {
				Description: "The description of the project.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"is_private": {
				Description: "A boolean to state if the project is private or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceBitbucketProjectCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	project, err := client.Workspaces.CreateProject(
		&gobb.ProjectOptions{
			Owner:       resourceData.Get("workspace").(string),
			Name:        resourceData.Get("name").(string),
			Key:         resourceData.Get("key").(string),
			Description: resourceData.Get("description").(string),
			IsPrivate:   resourceData.Get("is_private").(bool),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to create project with error: %s", err))
	}

	resourceData.SetId(project.Uuid)

	return resourceBitbucketProjectRead(ctx, resourceData, meta)
}

func resourceBitbucketProjectRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	project, err := client.Workspaces.GetProject(
		&gobb.ProjectOptions{
			Owner: resourceData.Get("workspace").(string),
			Key:   resourceData.Get("key").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get project with error: %s", err))
	}

	_ = resourceData.Set("name", project.Name)
	_ = resourceData.Set("key", project.Key)
	_ = resourceData.Set("description", project.Description)
	_ = resourceData.Set("is_private", project.Is_private)

	resourceData.SetId(project.Uuid)

	return nil
}

func resourceBitbucketProjectUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	_, err := client.Workspaces.UpdateProject(
		&gobb.ProjectOptions{
			Uuid:        resourceData.Id(),
			Owner:       resourceData.Get("workspace").(string),
			Name:        resourceData.Get("name").(string),
			Key:         resourceData.Get("key").(string),
			Description: resourceData.Get("description").(string),
			IsPrivate:   resourceData.Get("is_private").(bool),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to update project with error: %s", err))
	}

	return resourceBitbucketProjectRead(ctx, resourceData, meta)
}

func resourceBitbucketProjectDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	_, err := client.Workspaces.DeleteProject(
		&gobb.ProjectOptions{
			Owner: resourceData.Get("workspace").(string),
			Key:   resourceData.Get("key").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete project with error: %s", err))
	}

	resourceData.SetId("")

	return nil
}

func resourceBitbucketProjectImport(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{resourceData}

	splitID := strings.Split(resourceData.Id(), "/")
	if len(splitID) < 2 {
		return ret, fmt.Errorf("invalid import ID. It must to be in this format \"<workspace-slug|workspace-uuid>/<project-key>\"")
	}

	_ = resourceData.Set("workspace", splitID[0])
	_ = resourceData.Set("key", splitID[1])

	_ = resourceBitbucketProjectRead(ctx, resourceData, meta)

	return ret, nil
}

func validateProjectKey(val interface{}, path cty.Path) diag.Diagnostics {
	match, _ := regexp.MatchString("^[A-Za-z][A-Za-z0-9_]+$", val.(string))
	if !match {
		return diag.FromErr(fmt.Errorf("project keys must start with a letter and may only consist of ASCII letters, numbers and underscores (A-Z, a-z, 0-9, _)"))
	}

	return diag.Diagnostics{}
}
