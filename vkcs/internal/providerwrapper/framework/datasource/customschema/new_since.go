package customschema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	jsonschema "github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
)

func addNewSinceToAttr(attr schema.Attribute, attrJSON jsonschema.SchemaJSON) string {
	if attr.GetDeprecationMessage() == "" && attrJSON.NewSince != "" {
		return attr.GetDescription() + fmt.Sprintf("_new_since_%s_.", attrJSON.NewSince)
	}
	return attr.GetDescription()
}

func addNewSinceToBlock(block schema.Block, blockJSON jsonschema.SchemaJSON) string {
	if block.GetDeprecationMessage() == "" && blockJSON.NewSince != "" {
		return block.GetDescription() + fmt.Sprintf("_new_since_%s_.", blockJSON.NewSince)
	}
	return block.GetDescription()
}
