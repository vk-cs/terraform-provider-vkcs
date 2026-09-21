package kms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecretCreateContext,
		ReadContext:   resourceSecretReadContext,
		UpdateContext: resourceSecretUpdateContext,
		DeleteContext: resourceSecretDeleteContext,
	}
}

func resourceSecretCreateContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	diags := make([]diag.Diagnostic, 0)
	return diags
}

func resourceSecretReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	diags := make([]diag.Diagnostic, 0)
	return diags
}

func resourceSecretUpdateContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	diags := make([]diag.Diagnostic, 0)
	return diags
}

func resourceSecretDeleteContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	diags := make([]diag.Diagnostic, 0)
	return diags
}
