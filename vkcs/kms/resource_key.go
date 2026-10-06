package kms

import (
	"context"

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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key name",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Key type, take a look at type parameter in OpenBao documentation: https://openbao.org/docs/api/secret/transit/#parameters",
				Default:     "aes256-gcm96",
			},
			"deletion_allowed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A flag that shows if key can be deleted or not",
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

	name, ok := d.Get("name").(string)
	if !ok {
		return diag.Errorf("Error retrieving name from resource: %s", err)
	}

	t, ok := d.Get("type").(string)
	if !ok {
		return diag.Errorf("Error retrieving type from resource: %s", err)
	}

	_, err = createKey(kmsV1Client, keys.CreateOpts{
		Name: name,
		Type: t,
	})
	if err != nil {
		return diag.Errorf("Error creating key: %s", err)
	}

	deletionAllowed, ok := d.Get("deletion_allowed").(bool)
	if !ok {
		return diag.Errorf("Error retrieving deletion_allowed from resource: %s", err)
	}

	_, err = updateKey(kmsV1Client, name, keys.UpdateOpts{
		DeletionAllowed: deletionAllowed,
	})
	if err != nil {
		return diag.Errorf("Error updating key configuration: %s", err)
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

	name, ok := d.Get("name").(string)
	if !ok {
		return diag.Errorf("Error retrieving name from resource: %s", err)
	}

	key, err := getKey(kmsV1Client, name)
	if err != nil {
		return diag.Errorf("Error getting key: %s", err)
	}

	err = d.Set("name", key.Data.Name)
	if err != nil {
		return diag.Errorf("Error setting name: %s", err)
	}

	err = d.Set("type", key.Data.Type)
	if err != nil {
		return diag.Errorf("Error setting type: %s", err)
	}

	err = d.Set("deletion_allowed", key.Data.DeletionAllowed)
	if err != nil {
		return diag.Errorf("Error setting deletion_allowed: %s", err)
	}

	return nil
}

func resourceKeyUpdateContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	deletionAllowed, ok := d.Get("deletion_allowed").(bool)
	if !ok {
		return diag.Errorf("Error retrieving deletion_allowed from resource: %s", err)
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

	name, ok := d.Get("name").(string)
	if !ok {
		return diag.Errorf("Error retrieving name from resource: %s", err)
	}

	err = deleteKey(kmsV1Client, name)
	if err != nil {
		return diag.Errorf("Error deleting key: %s", err)
	}
	return nil
}
