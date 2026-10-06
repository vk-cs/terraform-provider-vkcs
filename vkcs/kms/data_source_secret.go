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
		Description: "A data source containing a KMS secret",
		ReadContext: dataSourceSecretReadContext,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			SecretFieldPath: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifier of a secret",
			},
			SecretFieldDataJSON: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Data stored in secret as a JSON object",
			},
			SecretFieldCreatedTime: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Time of secret's creation",
			},
			SecretFieldVersion: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Current version of secret",
			},
			SecretFieldDeleteProtection: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A flag that represents if secret can be deleted or not",
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

	path, ok := d.Get(SecretFieldPath).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, SecretFieldPath)
	}

	secret, err := getSecret(kmsV1Client, path)
	if err != nil {
		return diag.Errorf("Error getting secret: %s", err)
	}

	deleteProtection, err := getSecretDeleteProtection(kmsV1Client, path)
	if err != nil {
		return diag.Errorf("Error getting secret delete protection: %s", err)
	}

	d.SetId(path)

	secretString, err := json.Marshal(secret.Data.Data)
	if err != nil {
		return diag.Errorf("Error creating data_json: %s", err)
	}

	err = d.Set(SecretFieldDataJSON, string(secretString))
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, SecretFieldDataJSON, err)
	}

	err = d.Set(SecretFieldCreatedTime, secret.Data.Metadata.CreatedTime.String())
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, SecretFieldCreatedTime, err)
	}

	err = d.Set(SecretFieldVersion, secret.Data.Metadata.Version)
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, SecretFieldVersion, err)
	}

	err = d.Set(SecretFieldDeleteProtection, deleteProtection)
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, SecretFieldDeleteProtection, err)
	}

	return nil
}
