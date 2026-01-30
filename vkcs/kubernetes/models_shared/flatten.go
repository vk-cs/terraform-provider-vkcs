package models_shared

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FlattenStringSet(strSet []string) (types.Set, diag.Diagnostics) {
	if len(strSet) == 0 {
		return types.SetNull(types.StringType), nil
	}

	resList := make([]attr.Value, len(strSet))
	for i, reg := range strSet {
		resList[i] = types.StringValue(reg)
	}

	return types.SetValue(types.StringType, resList)
}

func FlattenStringMap(strMap map[string]string) (types.Map, diag.Diagnostics) {
	if len(strMap) == 0 {
		return types.MapNull(types.StringType), nil
	}

	resMap := make(map[string]attr.Value, len(strMap))
	for k, v := range strMap {
		resMap[k] = types.StringValue(v)
	}

	return types.MapValue(types.StringType, resMap)
}
