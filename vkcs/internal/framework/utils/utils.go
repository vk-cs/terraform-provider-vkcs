package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func IsKnown(v attr.Value) bool {
	return !v.IsNull() && !v.IsUnknown()
}

func GetFirstNotEmptyValue(values ...types.String) string {
	for _, value := range values {
		if len(value.ValueString()) != 0 {
			return value.ValueString()
		}
	}

	return ""
}
