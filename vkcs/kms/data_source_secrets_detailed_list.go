package kms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecretsDetailedList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecretsDetailedListReadContext,
	}
}

func dataSourceSecretsDetailedListReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	diags := make([]diag.Diagnostic, 0)
	return diags
}
