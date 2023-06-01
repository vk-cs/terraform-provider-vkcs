package transform

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonschema "github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
)

func DecodeValueType(input attr.Value) string {
	switch input.(type) {
	case types.Bool:
		return jsonschema.SchemaTypeBool
	case types.Float64:
		return jsonschema.SchemaTypeFloat
	case types.Int64:
		return jsonschema.SchemaTypeInt
	case types.List:
		return jsonschema.SchemaTypeList
	case types.Map:
		return jsonschema.SchemaTypeMap
	case types.Number:
		return jsonschema.SchemaTypeNumber
	case types.Object:
		return jsonschema.SchemaTypeObject
	case types.Set:
		return jsonschema.SchemaTypeSet
	case types.String:
		return jsonschema.SchemaTypeString
	}
	return ""
}
