package diagutil

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Warningf creates a Diagnostics with a single Warning level Diagnostic entry.
func Warningf(format string, a ...interface{}) diag.Diagnostics {
	return diag.Diagnostics{diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  fmt.Sprintf(format, a...),
	}}
}
