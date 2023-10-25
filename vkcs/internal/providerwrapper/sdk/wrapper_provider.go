package sdk

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	jsonschema "github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/providerwrapper"
)

func WrapProvider(p *schema.Provider) (*schema.Provider, error) {
	var wrapperJSON jsonschema.ProviderWrapper
	err := json.Unmarshal([]byte(providerwrapper.ProviderSchemaJSON), &wrapperJSON)
	if err != nil {
		return nil, err
	}

	providerJSON := wrapperJSON.ProviderSchema

	for name, rs := range p.ResourcesMap {
		rsJSON := providerJSON.ResourcesMap[name]

		if rs.DeprecationMessage != "" && !strings.Contains(strings.ToLower(rs.Description), "deprecated") {
			rs.Description += fmt.Sprintf(" **Deprecated** %s.", strings.TrimSuffix(rs.DeprecationMessage, "."))
		}

		if rs.DeprecationMessage == "" && rsJSON.NewSince != "" {
			rs.Description += fmt.Sprintf("_new_since_%s_.", rsJSON.NewSince)
		}

		for propertyName, propertySchema := range rs.Schema {
			// Get the same from the base json
			propertyJSON := rsJSON.Schema[propertyName]
			customizeSchema(propertyJSON, propertySchema, propertyName)
		}
	}

	for name, ds := range p.DataSourcesMap {
		dsJSON := providerJSON.DataSourcesMap[name]

		if ds.DeprecationMessage != "" && !strings.Contains(strings.ToLower(ds.Description), "deprecated") {
			ds.Description += fmt.Sprintf(" **Deprecated** %s.", strings.TrimSuffix(ds.DeprecationMessage, "."))
		}

		if ds.DeprecationMessage == "" && dsJSON.NewSince != "" {
			ds.Description += fmt.Sprintf("_new_since_%s_.", dsJSON.NewSince)
		}

		for propertyName, propertySchema := range ds.Schema {
			// Get the same from the base json
			propertyJSON := dsJSON.Schema[propertyName]
			customizeSchema(propertyJSON, propertySchema, propertyName)
		}
	}

	return p, nil
}

func customizeSchema(sJSON jsonschema.SchemaJSON, s *schema.Schema, nodeName string) {
	if s.Deprecated != "" && !strings.Contains(strings.ToLower(s.Description), "deprecated") {
		s.Description += fmt.Sprintf(" **Deprecated** %s.", strings.TrimSuffix(s.Deprecated, "."))
	}

	if s.Deprecated == "" && sJSON.NewSince != "" {
		s.Description += fmt.Sprintf("_new_since_%s_.", sJSON.NewSince)
	}

	if nodeIsBlock(s) {
		current := s.Elem.(*schema.Resource).Schema
		var base map[string]jsonschema.SchemaJSON
		if jsonschema.NodeIsBlock(sJSON) {
			base = sJSON.Elem.(jsonschema.ResourceJSON).Schema
		}
		for k, c := range current {
			b := base[k]
			customizeSchema(b, c, k)
		}
	}
}

func nodeIsBlock(input *schema.Schema) bool {
	if input.Type == schema.TypeList || input.Type == schema.TypeSet {
		if _, ok := input.Elem.(*schema.Resource); ok {
			return true
		}
	}

	return false
}
