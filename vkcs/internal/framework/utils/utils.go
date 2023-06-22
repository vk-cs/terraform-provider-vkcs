package utils

import "github.com/hashicorp/terraform-plugin-framework/attr"

func IsKnown(v attr.Value) bool {
	return !v.IsNull() && !v.IsUnknown()
}
