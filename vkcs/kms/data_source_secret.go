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
		Schema: map[string]*schema.Schema{
			"name": {
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
		},
	}
}

func dataSourceSecretReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	name, ok := d.Get("name").(string)
	if !ok {
		return diag.Errorf("Error retrieving name from resource: %s", err)
	}

	secret, err := getSecret(kmsV1Client, name)
	if err != nil {
		return diag.Errorf("Error listing secrets: %s", err)
	}

	d.SetId(name)

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

	return nil
}
