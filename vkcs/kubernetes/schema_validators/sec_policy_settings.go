package schema_validators

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ validator.String = (*SecPolicySettingsValidator)(nil)
)

type SecPolicySettingsValidator struct{}

func (v SecPolicySettingsValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v SecPolicySettingsValidator) MarkdownDescription(ctx context.Context) string {
	return "Value must be a valid JSON object."
}

func (v SecPolicySettingsValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	policySettings := req.ConfigValue.ValueString()

	var pSettings map[string]interface{}
	if err := json.Unmarshal([]byte(policySettings), &pSettings); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid security policy settings JSON",
			fmt.Sprintf(
				"Value is not valid JSON: %s. "+
					"Ensure it is a properly formatted JSON object (e.g. {\"key\": \"value\"}).",
				err.Error(),
			),
		)
	}
}
