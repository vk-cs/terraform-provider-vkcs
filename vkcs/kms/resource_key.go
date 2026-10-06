package kms

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/keys"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func ResourceKey() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource representing KMS key",
		CreateContext: resourceKeyCreateContext,
		ReadContext:   resourceKeyReadContext,
		UpdateContext: resourceKeyUpdateContext,
		DeleteContext: resourceKeyDeleteContext,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			KeyFieldName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key name",
			},
			KeyFieldType: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Key type, choose one from list: `aes128-gcm96`, `aes256-gcm96` (default), `chacha20-poly1305`, `xchacha20-poly1305`",
				Default:     "aes256-gcm96",
			},
			KeyFieldDeletionAllowed: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A flag that represents if key can be deleted or not",
				Default:     false,
			},
		},
	}
}

func resourceKeyCreateContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	name, ok := d.Get(KeyFieldName).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, KeyFieldName)
	}

	t, ok := d.Get(KeyFieldType).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, KeyFieldType)
	}

	_, err = createKey(kmsV1Client, name, keys.CreateOpts{
		Type: t,
	})
	if err != nil {
		return diag.Errorf("Error creating key: %s", err)
	}

	deletionAllowed, ok := d.Get(KeyFieldDeletionAllowed).(bool)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, KeyFieldDeletionAllowed)
	}

	// If deletion_allowed is true - we need to wait for the key to be ready, then set this flag
	if deletionAllowed {
		var err error
		for range retriesCount {
			_, err = updateKey(kmsV1Client, name, keys.UpdateOpts{
				DeletionAllowed: deletionAllowed,
			})
			if err == nil {
				break
			}
			time.Sleep(defaultThreshold)
		}
		if err != nil {
			return diag.Errorf("Error updating key configuration: %s", err)
		}
	}

	d.SetId(name)
	return nil
}

func resourceKeyReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func resourceKeyUpdateContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	deletionAllowed, ok := d.Get(KeyFieldDeletionAllowed).(bool)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, KeyFieldDeletionAllowed)
	}

	_, err = updateKey(kmsV1Client, d.Id(), keys.UpdateOpts{
		DeletionAllowed: deletionAllowed,
	})
	if err != nil {
		return diag.Errorf("Error updating key configuration: %s", err)
	}

	return nil
}

func resourceKeyDeleteContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	name, ok := d.Get(KeyFieldName).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, KeyFieldName)
	}

	err = deleteKey(kmsV1Client, name)
	if err != nil {
		return diag.Errorf("Error deleting key: %s", err)
	}

	// Waiting for the key to be fully deleted
	for range retriesCount {
		var err error
		_, err = getKey(kmsV1Client, name)
		if err != nil {
			break
		}
		time.Sleep(defaultThreshold)
	}

	return nil
}
