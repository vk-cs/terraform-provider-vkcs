package kms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeyCreateContext,
		ReadContext:   resourceKeyReadContext,
		UpdateContext: resourceKeyUpdateContext,
		DeleteContext: resourceKeyDeleteContext,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deletion_protection": {
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
	return nil
}

func resourceKeyUpdateContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourceKeyDeleteContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
