package customschema

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func addDeprecatedToAttr(attr schema.Attribute) string {
	if attr.GetDeprecationMessage() != "" && !strings.Contains(strings.ToLower(attr.GetDescription()), "deprecated") {
		return attr.GetDescription() + fmt.Sprintf(" **Deprecated** %s.", strings.TrimSuffix(attr.GetDeprecationMessage(), "."))
	}
	return attr.GetDescription()
}

func addDeprecatedToBlock(block schema.Block) string {
	if block.GetDeprecationMessage() != "" && !strings.Contains(strings.ToLower(block.GetDescription()), "deprecated") {
		return block.GetDescription() + fmt.Sprintf(" **Deprecated** %s.", strings.TrimSuffix(block.GetDeprecationMessage(), "."))
	}
	return block.GetDescription()
}
