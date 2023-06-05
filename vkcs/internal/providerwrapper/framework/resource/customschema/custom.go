package customschema

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	jsonschema "github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
)

func CustomizeSchema(sJSON jsonschema.ResourceJSON, s schema.Schema) schema.Schema {
	if s.DeprecationMessage != "" && !strings.Contains(strings.ToLower(s.Description), "deprecated") {
		s.Description += fmt.Sprintf(" **Deprecated** %s.", strings.TrimSuffix(s.DeprecationMessage, "."))
	}

	if s.DeprecationMessage == "" && sJSON.NewSince != "" {
		s.Description += fmt.Sprintf(" **New since %s**.", sJSON.NewSince)
	}

	for name, attr := range s.Attributes {
		attrJSON := sJSON.Schema[name]
		s.Attributes[name] = customizeAttribute(attrJSON, attr, name)
	}

	for name, block := range s.Blocks {
		blockJSON := sJSON.Schema[name]
		customizeBlock(blockJSON, block, name)
	}

	return s
}

func customizeBlock(blockJSON jsonschema.SchemaJSON, block schema.Block, name string) schema.Block {
	switch b := block.(type) {
	case schema.SingleNestedBlock:
		b.Description = addNewSinceToBlock(b, blockJSON)
		b.Description = addDeprecatedToBlock(b)

		var nestedAttrsJSONs map[string]jsonschema.SchemaJSON
		if jsonschema.NodeIsBlock(blockJSON) {
			nestedAttrsJSONs = blockJSON.Elem.(jsonschema.ResourceJSON).Schema
		}

		for nestedName, nestedAttr := range b.Attributes {
			b.Attributes[nestedName] = customizeAttribute(nestedAttrsJSONs[nestedName], nestedAttr, nestedName)
		}

		return b
	case schema.ListNestedBlock:
		b.Description = addNewSinceToBlock(b, blockJSON)
		b.Description = addDeprecatedToBlock(b)

		var nestedAttrsJSONs map[string]jsonschema.SchemaJSON
		if jsonschema.NodeIsBlock(blockJSON) {
			nestedAttrsJSONs = blockJSON.Elem.(jsonschema.ResourceJSON).Schema
		}

		for nestedName, nestedAttr := range b.NestedObject.Attributes {
			b.NestedObject.Attributes[nestedName] = customizeAttribute(nestedAttrsJSONs[nestedName], nestedAttr, nestedName)
		}

		for nestedName, nestedBlock := range b.NestedObject.Blocks {
			b.NestedObject.Blocks[nestedName] = customizeBlock(nestedAttrsJSONs[nestedName], nestedBlock, nestedName)
		}

		return b
	case schema.SetNestedBlock:
		b.Description = addNewSinceToBlock(b, blockJSON)
		b.Description = addDeprecatedToBlock(b)

		var nestedAttrsJSONs map[string]jsonschema.SchemaJSON
		if jsonschema.NodeIsBlock(blockJSON) {
			nestedAttrsJSONs = blockJSON.Elem.(jsonschema.ResourceJSON).Schema
		}

		for nestedName, nestedAttr := range b.NestedObject.Attributes {
			b.NestedObject.Attributes[nestedName] = customizeAttribute(nestedAttrsJSONs[nestedName], nestedAttr, nestedName)
		}

		for nestedName, nestedBlock := range b.NestedObject.Blocks {
			b.NestedObject.Blocks[nestedName] = customizeBlock(nestedAttrsJSONs[nestedName], nestedBlock, nestedName)
		}

		return b
	}

	return block
}

func customizeAttribute(attrJSON jsonschema.SchemaJSON, attr schema.Attribute, name string) schema.Attribute {
	switch a := attr.(type) {
	case schema.BoolAttribute:
		a.Description = addNewSinceToAttr(a, attrJSON)
		a.Description = addDeprecatedToAttr(a)
		return a
	case schema.Float64Attribute:
		a.Description = addNewSinceToAttr(a, attrJSON)
		a.Description = addDeprecatedToAttr(a)
		return a
	case schema.Int64Attribute:
		a.Description = addNewSinceToAttr(a, attrJSON)
		a.Description = addDeprecatedToAttr(a)
		return a
	case schema.ListAttribute:
		a.Description = addNewSinceToAttr(a, attrJSON)
		a.Description = addDeprecatedToAttr(a)
		return a
	case schema.MapAttribute:
		a.Description = addNewSinceToAttr(a, attrJSON)
		a.Description = addDeprecatedToAttr(a)
		return a
	case schema.NumberAttribute:
		a.Description = addNewSinceToAttr(a, attrJSON)
		a.Description = addDeprecatedToAttr(a)
		return a
	case schema.ObjectAttribute:
		a.Description = addNewSinceToAttr(a, attrJSON)
		a.Description = addDeprecatedToAttr(a)
		return a
	case schema.SetAttribute:
		a.Description = addNewSinceToAttr(a, attrJSON)
		a.Description = addDeprecatedToAttr(a)
		return a
	case schema.StringAttribute:
		a.Description = addNewSinceToAttr(a, attrJSON)
		a.Description = addDeprecatedToAttr(a)
		return a
	case schema.ListNestedAttribute:
		a.Description = addNewSinceToAttr(a, attrJSON)
		a.Description = addDeprecatedToAttr(a)

		var nestedAttrsJSONs map[string]jsonschema.SchemaJSON
		if jsonschema.NodeIsBlock(attrJSON) {
			nestedAttrsJSONs = attrJSON.Elem.(jsonschema.ResourceJSON).Schema
		}

		for nestedName, nestedAttr := range a.NestedObject.Attributes {
			a.NestedObject.Attributes[nestedName] = customizeAttribute(nestedAttrsJSONs[nestedName], nestedAttr, nestedName)
		}

		return a
	case schema.SetNestedAttribute:
		a.Description = addNewSinceToAttr(a, attrJSON)
		a.Description = addDeprecatedToAttr(a)

		var nestedAttrsJSONs map[string]jsonschema.SchemaJSON
		if jsonschema.NodeIsBlock(attrJSON) {
			nestedAttrsJSONs = attrJSON.Elem.(jsonschema.ResourceJSON).Schema
		}

		for nestedName, nestedAttr := range a.NestedObject.Attributes {
			a.NestedObject.Attributes[nestedName] = customizeAttribute(nestedAttrsJSONs[nestedName], nestedAttr, nestedName)
		}

		return a
	case schema.MapNestedAttribute:
		a.Description = addNewSinceToAttr(a, attrJSON)
		a.Description = addDeprecatedToAttr(a)

		var nestedAttrsJSONs map[string]jsonschema.SchemaJSON
		if jsonschema.NodeIsBlock(attrJSON) {
			nestedAttrsJSONs = attrJSON.Elem.(jsonschema.ResourceJSON).Schema
		}

		for nestedName, nestedAttr := range a.NestedObject.Attributes {
			a.NestedObject.Attributes[nestedName] = customizeAttribute(nestedAttrsJSONs[nestedName], nestedAttr, nestedName)
		}

		return a
	}

	return attr
}
