package kms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecret() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecretReadContext,
	}
}

func dataSourceSecretReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	diags := make([]diag.Diagnostic, 0)
	return diags
}
