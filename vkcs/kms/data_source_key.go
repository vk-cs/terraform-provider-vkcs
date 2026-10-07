package kms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func DataSourceKey() *schema.Resource {
	return &schema.Resource{
		Description: "A data source containing a KMS key",
		ReadContext: dataSourceKeyReadContext,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			KeyFieldName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the KMS key to read.",
			},
			KeyFieldType: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Type of the key: `aes128-gcm96`, `aes256-gcm96`, `chacha20-poly1305`, or `xchacha20-poly1305`. Populated from KMS when the data source is read. This field does not filter the lookup.",
			},
			KeyFieldDeletionAllowed: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether KMS allows the key to be deleted, populated from KMS when the data source is read. This field does not change the deletion policy.",
			},
		},
	}
}

func dataSourceKeyReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	name, ok := d.Get(KeyFieldName).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, KeyFieldName)
	}

	key, err := getKey(kmsV1Client, name)
	if err != nil {
		return diag.Errorf("Error getting key: %s", err)
	}

	d.SetId(name)

	err = d.Set(KeyFieldName, key.Data.Name)
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, KeyFieldName, err)
	}

	err = d.Set(KeyFieldType, key.Data.Type)
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, KeyFieldType, err)
	}

	err = d.Set(KeyFieldDeletionAllowed, key.Data.DeletionAllowed)
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, KeyFieldDeletionAllowed, err)
	}

	return nil
}
