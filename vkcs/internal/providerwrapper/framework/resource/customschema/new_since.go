package customschema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	jsonschema "github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
)

func addNewSinceToAttr(attr schema.Attribute, attrJSON jsonschema.SchemaJSON) string {
	if attr.GetDeprecationMessage() == "" && attrJSON.NewSince != "" {
		return attr.GetDescription() + fmt.Sprintf(" **New since %s**.", attrJSON.NewSince)
	}
	return attr.GetDescription()
}

func addNewSinceToBlock(block schema.Block, blockJSON jsonschema.SchemaJSON) string {
	if block.GetDeprecationMessage() == "" && blockJSON.NewSince != "" {
		return block.GetDescription() + fmt.Sprintf(" **New since %s**.", blockJSON.NewSince)
	}
	return block.GetDescription()
}
