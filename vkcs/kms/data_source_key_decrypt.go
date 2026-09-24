package kms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceKeyDecrypt() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeyDecryptReadContext,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ciphertext": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceKeyDecryptReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
