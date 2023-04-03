package vkcs

import (
	"context"
	"errors"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabaseCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	if diff.Id() != "" && diff.HasChange("cloud_monitoring_enabled") {
		t, exists := diff.GetOk("datastore.0.type")
		if !exists {
			return errors.New("datastore.0.type is not found")
		}
		if exists && isOperationNotSupported(t.(string), Redis, MongoDB) {
			return diff.ForceNew("cloud_monitoring_enabled")
		}
	}
	return nil
}

func checkDBNetworks(
	rawNetworks []interface{}, path cty.Path, diags diag.Diagnostics,
) diag.Diagnostics {
	if len(rawNetworks) > 1 {
		p := path
		p = append(p,
			cty.GetAttrStep{Name: "network"},
			cty.IndexStep{Key: cty.NumberIntVal(1)},
		)
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Warning,
			Summary:       "Multiple networks are deprecated",
			Detail:        "Multiple networks are deprecated and won't be supported in next major release.",
			AttributePath: p,
		})
	}
	return diags
}
