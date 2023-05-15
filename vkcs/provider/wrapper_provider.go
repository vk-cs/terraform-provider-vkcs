package provider

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson"
)

//go:generate go run ../internal/generate/providerschema/main.go -schemajson ../../.release/provider-schema.json

func wrapProvider(p *schema.Provider) (*schema.Provider, error) {
	var wrapperBase providerjson.ProviderWrapper
	err := json.Unmarshal([]byte(providerSchemaBase), &wrapperBase)
	if err != nil {
		return nil, err
	}

	base := wrapperBase.ProviderSchema

	for resource, rs := range p.ResourcesMap {
		bRs := base.ResourcesMap[resource]
		if bRs.NewSince != "" {
			rs.Description += fmt.Sprintf(" **New since %s**.", bRs.NewSince)
		}

		for propertyName, propertySchema := range rs.Schema {
			// Get the same from the base json
			bS := bRs.Schema[propertyName]
			customizeSchema(bS, propertySchema, propertyName)
		}
	}

	for resource, rs := range p.DataSourcesMap {
		bRs := base.DataSourcesMap[resource]
		if bRs.NewSince != "" {
			rs.Description += fmt.Sprintf(" **New since %s**.", bRs.NewSince)
		}

		for propertyName, propertySchema := range rs.Schema {
			// Get the same from the base json
			bS := bRs.Schema[propertyName]
			customizeSchema(bS, propertySchema, propertyName)
		}
	}

	return p, nil
}

func customizeSchema(baseS providerjson.SchemaJSON, s *schema.Schema, nodeName string) {
	if s.Deprecated != "" && !strings.Contains(strings.ToLower(s.Description), "deprecated") {
		s.Description += fmt.Sprintf(" **Deprecated** %s.", strings.TrimSuffix(s.Deprecated, "."))
	}

	if s.Deprecated == "" && baseS.NewSince != "" {
		s.Description += fmt.Sprintf(" **New since %s**.", baseS.NewSince)
	}

	if nodeIsBlock(s) {
		current := s.Elem.(*schema.Resource).Schema
		var base map[string]providerjson.SchemaJSON
		if providerjson.NodeIsBlock(baseS) {
			base = baseS.Elem.(providerjson.ResourceJSON).Schema
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
