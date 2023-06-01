package sdk

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	jsonschema "github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
)

func ResourceFromRaw(input *schema.Resource) (*jsonschema.ResourceJSON, error) {
	if input == nil {
		return nil, fmt.Errorf("resource not found")
	}

	result := &jsonschema.ResourceJSON{
		Description: input.Description,
		Deprecated:  input.DeprecationMessage,
	}
	translatedSchema := make(map[string]jsonschema.SchemaJSON)

	for k, s := range input.Schema {
		translatedSchema[k] = SchemaFromRaw(s)
	}
	result.Schema = translatedSchema

	if input.Timeouts != nil {
		timeouts := &jsonschema.ResourceTimeoutJSON{}
		if t := input.Timeouts; t != nil {
			if t.Create != nil {
				timeouts.Create = int(t.Create.Minutes())
			}
			if t.Read != nil {
				timeouts.Read = int(t.Read.Minutes())
			}
			if t.Delete != nil {
				timeouts.Delete = int(t.Delete.Minutes())
			}
			if t.Update != nil {
				timeouts.Update = int(t.Update.Minutes())
			}
		}
		result.Timeouts = timeouts
	}

	return result, nil
}

func SchemaFromRaw(input *schema.Schema) jsonschema.SchemaJSON {
	return jsonschema.SchemaJSON{
		Type:        input.Type.String(),
		ConfigMode:  decodeConfigMode(input.ConfigMode),
		Optional:    input.Optional,
		Required:    input.Required,
		Default:     input.Default,
		Description: input.Description,
		Computed:    input.Computed,
		ForceNew:    input.ForceNew,
		Elem:        decodeElem(input.Elem),
		MaxItems:    input.MaxItems,
		MinItems:    input.MinItems,
		Deprecated:  input.Deprecated,
	}
}

func decodeConfigMode(input schema.SchemaConfigMode) (out string) {
	switch input {
	case 1:
		out = "Auto"
	case 2:
		out = "Block"
	case 4:
		out = "Attribute"
	}
	return
}

func decodeElem(input interface{}) interface{} {
	switch t := input.(type) {
	case bool:
		return t
	case string:
		return t
	case int:
		return t
	case float32, float64:
		return t
	case *schema.Schema:
		return SchemaFromRaw(t)
	case *schema.Resource:
		r, _ := ResourceFromRaw(t)
		return r
	}
	return nil
}
