package providerjson

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceFromRaw(input *schema.Resource) (*ResourceJSON, error) {
	if input == nil {
		return nil, fmt.Errorf("resource not found")
	}

	result := &ResourceJSON{
		Description: input.Description,
		Deprecated:  input.DeprecationMessage,
	}
	translatedSchema := make(map[string]SchemaJSON)

	for k, s := range input.Schema {
		translatedSchema[k] = SchemaFromRaw(s)
	}
	result.Schema = translatedSchema

	if input.Timeouts != nil {
		timeouts := &ResourceTimeoutJSON{}
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

func SchemaFromRaw(input *schema.Schema) SchemaJSON {
	return SchemaJSON{
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

func SchemaFromMap(input map[string]interface{}) SchemaJSON {
	result := SchemaJSON{}
	if t, ok := input["type"]; ok {
		result.Type = t.(string)
	}

	if t, ok := input["config_mode"]; ok {
		result.ConfigMode = t.(string)
	}

	if t, ok := input["optional"]; ok {
		result.Optional = t.(bool)
	}

	if t, ok := input["required"]; ok {
		result.Required = t.(bool)
	}

	if t, ok := input["default"]; ok {
		result.Default = t
	}

	if t, ok := input["description"]; ok {
		result.Description = t.(string)
	}

	if t, ok := input["computed"]; ok {
		result.Computed = t.(bool)
	}

	if t, ok := input["force_new"]; ok {
		result.ForceNew = t.(bool)
	}

	if t, ok := input["elem"]; ok {
		result.Elem = decodeElem(t)
	}

	if t, ok := input["min_items"]; ok {
		result.MinItems = int(t.(float64))
	}

	if t, ok := input["max_items"]; ok {
		result.MaxItems = int(t.(float64))
	}

	if t, ok := input["new_since"]; ok {
		result.NewSince = t.(string)
	}

	if t, ok := input["deprecated"]; ok {
		result.Deprecated = t.(string)
	}

	return result
}

func ResourceFromMap(input map[string]interface{}) ResourceJSON {
	result := ResourceJSON{
		Schema: make(map[string]SchemaJSON, 0),
	}
	for k, v := range input {
		result.Schema[k] = SchemaFromMap(v.(map[string]interface{}))
	}
	return result
}

func addNewSinceToResource(base, current *ResourceJSON, curVersion string) {
	if base == nil {
		current.NewSince = curVersion
		return
	}

	current.NewSince = base.NewSince
	if current.NewSince != "" {
		return
	}

	for k, curS := range current.Schema {
		var baseS *SchemaJSON
		if b, ok := base.Schema[k]; ok {
			baseS = &b
		}
		addNewSinceToSchema(baseS, &curS, curVersion)
		current.Schema[k] = curS
	}
}

func addNewSinceToSchema(base, current *SchemaJSON, curVersion string) {
	if base == nil {
		current.NewSince = curVersion
		return
	}

	current.NewSince = base.NewSince
	if current.NewSince != "" {
		return
	}

	switch el := current.Elem.(type) {
	case *SchemaJSON:
		if baseEl, ok := base.Elem.(SchemaJSON); ok {
			addNewSinceToSchema(&baseEl, el, curVersion)
		}
	case *ResourceJSON:
		if baseEl, ok := base.Elem.(ResourceJSON); ok {
			addNewSinceToResource(&baseEl, el, curVersion)
		}
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

func ProviderFromRaw(input *ProviderJSON) (*ProviderSchemaJSON, error) {
	if input == nil {
		return nil, fmt.Errorf("provider was nil converting from raw")
	}

	result := &ProviderSchemaJSON{}

	providerSchema := make(map[string]SchemaJSON)
	resourceSchemas := make(map[string]ResourceJSON)
	dataSourceSchemas := make(map[string]ResourceJSON)

	for k, v := range input.Schema {
		providerSchema[k] = SchemaFromRaw(v)
	}

	for k, v := range input.ResourcesMap {
		resource, err := ResourceFromRaw(v)
		if err != nil {
			return nil, err
		}
		resourceSchemas[k] = *resource
	}

	for k, v := range input.DataSourcesMap {
		dataSource, err := ResourceFromRaw(v)
		if err != nil {
			return nil, err
		}
		dataSourceSchemas[k] = *dataSource
	}

	result.Schema = providerSchema
	result.ResourcesMap = resourceSchemas
	result.DataSourcesMap = dataSourceSchemas
	return result, nil
}

func addNewSince(base, current *ProviderSchemaJSON, curVersion string) {
	for k, r := range current.ResourcesMap {
		var baseR *ResourceJSON
		if b, ok := base.ResourcesMap[k]; ok {
			baseR = &b
		}
		addNewSinceToResource(baseR, &r, curVersion)
		current.ResourcesMap[k] = r
	}

	for k, r := range current.DataSourcesMap {
		var baseR *ResourceJSON
		if b, ok := base.DataSourcesMap[k]; ok {
			baseR = &b
		}
		addNewSinceToResource(baseR, &r, curVersion)
		current.DataSourcesMap[k] = r
	}
}

func WrappedProvider(input *ProviderJSON, wrapper *ProviderWrapper) (*ProviderWrapper, error) {
	schema, err := ProviderFromRaw(input)
	if err != nil {
		return nil, err
	}

	wrapper.ProviderSchema = schema

	return wrapper, nil
}
