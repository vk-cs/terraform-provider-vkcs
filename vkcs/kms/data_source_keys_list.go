package kms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceKeysList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeysListReadContext, Schema: map[string]*schema.Schema{
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceKeysListReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	diags := make([]diag.Diagnostic, 0)
	return diags
}
