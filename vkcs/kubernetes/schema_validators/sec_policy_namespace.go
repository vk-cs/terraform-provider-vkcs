package schema_validators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	re = regexp.MustCompile("^[*]?$|[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
)

var (
	_ validator.String = (*SecPolicyNamespaceValidator)(nil)
)

type SecPolicyNamespaceValidator struct{}

func (v SecPolicyNamespaceValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v SecPolicyNamespaceValidator) MarkdownDescription(ctx context.Context) string {
	return "Value must be a valid Kubernetes namespace or '*' (wildcard)."
}

func (v SecPolicyNamespaceValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	namespace := req.ConfigValue.ValueString()

	if !re.MatchString(namespace) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Kubernetes namespace",
			fmt.Sprintf(
				"Value '%s' is not a valid namespace. "+
					"It must consist of lowercase alphanumeric characters or '-', "+
					"start and end with an alphanumeric character, or be '*' for all namespaces.",
				namespace,
			),
		)
	}
}
