package kms

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func DataSourceSecret() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecretReadContext,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"data_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func dataSourceSecretReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	path, ok := d.Get("path").(string)
	if !ok {
		return diag.Errorf("Error retrieving name from resource: %s", err)
	}

	secret, err := getSecret(kmsV1Client, path)
	if err != nil {
		return diag.Errorf("Error getting secret: %s", err)
	}

	d.SetId(path)

	err = d.Set("data", secret.Data.Data)
	if err != nil {
		return diag.Errorf("Error setting data_json: %s", err)
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
