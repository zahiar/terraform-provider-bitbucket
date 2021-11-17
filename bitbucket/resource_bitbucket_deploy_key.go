package bitbucket

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gobb "github.com/ktrysmt/go-bitbucket"
)

func resourceBitbucketDeployKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBitbucketDeployKeyCreate,
		ReadContext:   resourceBitbucketDeployKeyRead,
		DeleteContext: resourceBitbucketDeployKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBitbucketDeployKeyImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the deploy key.",
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
			"label": {
				Description: "The label of the deploy key.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"key": {
				Description: "The public SSH key to attach to this deploy key.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceBitbucketDeployKeyCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	deployKey, err := client.Repositories.DeployKeys.Create(
		&gobb.DeployKeyOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
			Label:    resourceData.Get("label").(string),
			Key:      resourceData.Get("key").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to create deploy key with error: %s", err))
	}

	resourceData.SetId(strconv.Itoa(deployKey.Id))

	return resourceBitbucketDeployKeyRead(ctx, resourceData, meta)
}

func resourceBitbucketDeployKeyRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	deployKeyId, _ := strconv.Atoi(resourceData.Get("id").(string))
	deployKey, err := client.Repositories.DeployKeys.Get(
		&gobb.DeployKeyOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
			Id:       deployKeyId,
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get deploy key with error: %s", err))
	}

	_ = resourceData.Set("label", deployKey.Label)

	if deployKey.Comment != "" {
		_ = resourceData.Set("key", fmt.Sprintf("%s %s", deployKey.Key, deployKey.Comment))
	} else {
		_ = resourceData.Set("key", deployKey.Key)
	}

	resourceData.SetId(strconv.Itoa(deployKey.Id))

	return nil
}

func resourceBitbucketDeployKeyDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gobb.Client)

	deployKeyId, _ := strconv.Atoi(resourceData.Id())
	_, err := client.Repositories.DeployKeys.Delete(
		&gobb.DeployKeyOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
			Id:       deployKeyId,
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete deploy key with error: %s", err))
	}

	resourceData.SetId("")

	return nil
}

func resourceBitbucketDeployKeyImport(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{resourceData}

	splitID := strings.Split(resourceData.Id(), "/")
	if len(splitID) < 2 {
		return ret, fmt.Errorf("invalid import ID. It must to be in this format \"<workspace-slug|workspace-uuid>/<repository-name>/<deploy-key-id>\"")
	}

	_ = resourceData.Set("workspace", splitID[0])
	_ = resourceData.Set("repository", splitID[1])
	resourceData.SetId(splitID[2])

	_ = resourceBitbucketDeployKeyRead(ctx, resourceData, meta)

	return ret, nil
}
