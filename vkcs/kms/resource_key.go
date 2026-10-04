package kms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func ResourceKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeyCreateContext,
		ReadContext:   resourceKeyReadContext,
		UpdateContext: resourceKeyUpdateContext,
		DeleteContext: resourceKeyDeleteContext,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deletion_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceKeyCreateContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		return diag.Errorf("Error listing secrets: %s", err)
	}

	d.SetId(name)

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
	return nil
}

func resourceKeyDeleteContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
