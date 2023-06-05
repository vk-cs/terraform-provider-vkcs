package datasource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	jsonschema "github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
	"github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/transform"
)

func ResourceFromRaw(input *schema.Schema) (*jsonschema.ResourceJSON, error) {
	if input == nil {
		return nil, fmt.Errorf("data source not found")
	}

	result := &jsonschema.ResourceJSON{
		Description: input.Description,
		Deprecated:  input.DeprecationMessage,
	}
	translatedSchema := make(map[string]jsonschema.SchemaJSON)
	for k, a := range input.Attributes {
		translatedSchema[k] = SchemaFromAttribute(a)
	}
	for k, b := range input.Blocks {
		translatedSchema[k] = SchemaFromBlock(b)
	}
	result.Schema = translatedSchema

	return result, nil
}

func SchemaFromAttribute(input schema.Attribute) jsonschema.SchemaJSON {
	return jsonschema.SchemaJSON{
		Type:        decodeAttrType(input),
		Optional:    input.IsOptional(),
		Required:    input.IsRequired(),
		Description: input.GetDescription(),
		Computed:    input.IsComputed(),
		Deprecated:  input.GetDeprecationMessage(),
		Elem:        decodeAttrElem(input),
	}
}

func SchemaFromBlock(input schema.Block) jsonschema.SchemaJSON {
	var maxItems int
	if _, ok := input.(schema.SingleNestedBlock); ok {
		maxItems = 1
	}

	return jsonschema.SchemaJSON{
		Type:        decodeBlockType(input),
		Description: input.GetDescription(),
		Deprecated:  input.GetDeprecationMessage(),
		Elem:        decodeBlockElem(input),
		MaxItems:    maxItems,
	}
}

func decodeAttrType(input schema.Attribute) string {
	switch input.(type) {
	case schema.BoolAttribute:
		return jsonschema.SchemaTypeBool
	case schema.StringAttribute:
		return jsonschema.SchemaTypeString
	case schema.Int64Attribute:
		return jsonschema.SchemaTypeInt
	case schema.Float64Attribute:
		return jsonschema.SchemaTypeFloat
	case schema.ListAttribute:
		return jsonschema.SchemaTypeList
	case schema.ListNestedAttribute:
		return jsonschema.SchemaTypeList
	case schema.SetAttribute:
		return jsonschema.SchemaTypeSet
	case schema.SetNestedAttribute:
		return jsonschema.SchemaTypeSet
	case schema.MapAttribute:
		return jsonschema.SchemaTypeMap
	case schema.MapNestedAttribute:
		return jsonschema.SchemaTypeMap
	}

	return ""
}

func decodeBlockType(input schema.Block) string {
	switch input.(type) {
	case schema.SingleNestedBlock:
		return jsonschema.SchemaTypeList
	case schema.ListNestedBlock:
		return jsonschema.SchemaTypeList
	case schema.SetNestedBlock:
		return jsonschema.SchemaTypeSet
	}

	return ""
}

func decodeAttrElem(input schema.Attribute) interface{} {
	switch t := input.(type) {
	case schema.ListAttribute:
		return &jsonschema.SchemaJSON{Type: transform.DecodeValueType(t.ElementType.ValueType(context.Background()))}
	case schema.SetAttribute:
		return &jsonschema.SchemaJSON{Type: transform.DecodeValueType(t.ElementType.ValueType(context.Background()))}
	case schema.MapAttribute:
		return &jsonschema.SchemaJSON{Type: transform.DecodeValueType(t.ElementType.ValueType(context.Background()))}
	case schema.ListNestedAttribute:
		return decodeNestedAttributeObject(t.NestedObject)
	case schema.SetNestedAttribute:
		return decodeNestedAttributeObject(t.NestedObject)
	case schema.MapNestedAttribute:
		return decodeNestedAttributeObject(t.NestedObject)
	}

	return nil
}

func decodeBlockElem(input schema.Block) *jsonschema.ResourceJSON {
	switch t := input.(type) {
	case schema.SingleNestedBlock:
		m := make(map[string]jsonschema.SchemaJSON)
		for k, a := range t.Attributes {
			m[k] = SchemaFromAttribute(a)
		}
		for k, b := range t.Blocks {
			m[k] = SchemaFromBlock(b)
		}
		return &jsonschema.ResourceJSON{Schema: m}
	case schema.ListNestedBlock:
		return decodeNestedBlockObject(t.NestedObject)
	case schema.SetNestedBlock:
		return decodeNestedBlockObject(t.NestedObject)
	}

	return nil
}

func decodeNestedAttributeObject(input schema.NestedAttributeObject) *jsonschema.ResourceJSON {
	m := make(map[string]jsonschema.SchemaJSON)
	for k, a := range input.Attributes {
		m[k] = SchemaFromAttribute(a)
	}
	return &jsonschema.ResourceJSON{Schema: m}
}

func decodeNestedBlockObject(input schema.NestedBlockObject) *jsonschema.ResourceJSON {
	m := make(map[string]jsonschema.SchemaJSON)
	for k, a := range input.Attributes {
		m[k] = SchemaFromAttribute(a)
	}
	for k, b := range input.Blocks {
		m[k] = SchemaFromBlock(b)
	}
	return &jsonschema.ResourceJSON{Schema: m}
}
