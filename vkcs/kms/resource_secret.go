package kms

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/secrets"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func ResourceSecret() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource representing KMS secret",
		CreateContext: resourceSecretCreateContext,
		ReadContext:   resourceSecretReadContext,
		UpdateContext: resourceSecretUpdateContext,
		DeleteContext: resourceSecretDeleteContext,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identifier of a secret",
			},
			"data_json": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Data stored in secret as a JSON object",
			},
			"created_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Time of secret's creation",
			},
			"version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Current version of secret",
			},
		},
	}
}

func resourceSecretCreateContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	path, ok := d.Get("path").(string)
	if !ok {
		return diag.Errorf("Error retrieving path from resource: %s", err)
	}

	data, ok := d.Get("data_json").(string)
	if !ok {
		return diag.Errorf("Error retrieving data_json from resource: %s", err)
	}

	secret, err := createOrUpdateSecret(kmsV1Client, path, secrets.CreateOrUpdateOpts{
		Data: data,
	})
	if err != nil {
		return diag.Errorf("Error creating secret: %s", err)
	}

	d.SetId(path)

	err = d.Set("created_time", secret.Data.CreatedTime.String())
	if err != nil {
		return diag.Errorf("Error setting created_time: %s", err)
	}

	err = d.Set("version", secret.Data.Version)
	if err != nil {
		return diag.Errorf("Error setting version: %s", err)
	}

	return nil
}

func resourceSecretReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	path, ok := d.Get("path").(string)
	if !ok {
		return diag.Errorf("Error retrieving path from resource: %s", err)
	}

	secret, err := getSecret(kmsV1Client, path)
	if err != nil {
		return diag.Errorf("Error getting secret: %s", err)
	}

	secretString, err := json.Marshal(secret.Data.Data)
	if err != nil {
		return diag.Errorf("Error creating data_json: %s", err)
	}

	err = d.Set("data_json", string(secretString))
	if err != nil {
		return diag.Errorf("Error setting data_json: %s", err)
	}

	err = d.Set("created_time", secret.Data.Metadata.CreatedTime.String())
	if err != nil {
		return diag.Errorf("Error setting created_time: %s", err)
	}

	err = d.Set("version", secret.Data.Metadata.Version)
	if err != nil {
		return diag.Errorf("Error setting version: %s", err)
	}

	return nil
}

func resourceSecretUpdateContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	path, ok := d.Get("path").(string)
	if !ok {
		return diag.Errorf("Error retrieving path from resource: %s", err)
	}

	data, ok := d.Get("data_json").(string)
	if !ok {
		return diag.Errorf("Error retrieving data_json from resource: %s", err)
	}

	secret, err := createOrUpdateSecret(kmsV1Client, path, secrets.CreateOrUpdateOpts{
		Data: data,
	})
	if err != nil {
		return diag.Errorf("Error creating secret: %s", err)
	}

	err = d.Set("created_time", secret.Data.CreatedTime.String())
	if err != nil {
		return diag.Errorf("Error setting created_time: %s", err)
	}

	err = d.Set("version", secret.Data.Version)
	if err != nil {
		return diag.Errorf("Error setting version: %s", err)
	}

	return nil
}

func resourceSecretDeleteContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	path, ok := d.Get("path").(string)
	if !ok {
		return diag.Errorf("Error retrieving path from resource: %s", err)
	}

	err = deleteSecret(kmsV1Client, path)
	if err != nil {
		return diag.Errorf("Error deleting secret: %s", err)
	}
	return nil
}
