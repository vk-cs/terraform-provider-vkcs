package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ListEmpty(elementType attr.Type) basetypes.ListValue {
	return basetypes.NewListValueMust(elementType, []attr.Value{})
}
