package kms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceKeyEncrypt() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeyEncryptReadContext,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plaintext": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceKeyEncryptReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	diags := make([]diag.Diagnostic, 0)
	return diags
}
