package bitbucket

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gobb "github.com/ktrysmt/go-bitbucket"
)

func resourceBitbucketPipelineKeyPair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBitbucketPipelineKeyPairCreate,
		ReadContext:   resourceBitbucketPipelineKeyPairRead,
		DeleteContext: resourceBitbucketPipelineKeyPairDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the pipeline key pair.",
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
			"public_key": {
				Description: "The public SSH key part of this pipeline key pair.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"private_key": {
				Description: "The private SSH key part of this pipeline key pair.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceBitbucketPipelineKeyPairCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	_, err := client.Repositories.Repository.AddPipelineKeyPair(
		&gobb.RepositoryPipelineKeyPairOptions{
			Owner:      resourceData.Get("workspace").(string),
			RepoSlug:   resourceData.Get("repository").(string),
			PrivateKey: resourceData.Get("private_key").(string),
			PublicKey:  resourceData.Get("public_key").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to create pipeline key pair with error: %s", err))
	}

	resourceData.SetId(generatePipelineKeyPairId(resourceData.Get("workspace").(string), resourceData.Get("repository").(string)))

	return resourceBitbucketPipelineKeyPairRead(ctx, resourceData, meta)
}

func resourceBitbucketPipelineKeyPairRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	pipelineKeyPair, err := client.Repositories.Repository.GetPipelineKeyPair(
		&gobb.RepositoryPipelineKeyPairOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get pipeline key pair with error: %s", err))
	}

	_ = resourceData.Set("public_key", pipelineKeyPair.Public_key)

	resourceData.SetId(generatePipelineKeyPairId(resourceData.Get("workspace").(string), resourceData.Get("repository").(string)))

	return nil
}

func resourceBitbucketPipelineKeyPairDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).V2

	_, err := client.Repositories.Repository.DeletePipelineKeyPair(
		&gobb.RepositoryPipelineKeyPairOptions{
			Owner:    resourceData.Get("workspace").(string),
			RepoSlug: resourceData.Get("repository").(string),
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete pipeline key pair with error: %s", err))
	}

	resourceData.SetId("")

	return nil
}

func generatePipelineKeyPairId(workspace string, slug string) string {
	return fmt.Sprintf("%s-%s", workspace, slug)
}
