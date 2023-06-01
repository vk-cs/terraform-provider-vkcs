package provider

import (
	"github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
)

func addNewSince(base, current *schema.ProviderSchemaJSON, curVersion string) {
	for k, r := range current.ResourcesMap {
		var baseR *schema.ResourceJSON
		if b, ok := base.ResourcesMap[k]; ok {
			baseR = &b
		}
		addNewSinceToResource(baseR, &r, curVersion)
		current.ResourcesMap[k] = r
	}

	for k, r := range current.DataSourcesMap {
		var baseR *schema.ResourceJSON
		if b, ok := base.DataSourcesMap[k]; ok {
			baseR = &b
		}
		addNewSinceToResource(baseR, &r, curVersion)
		current.DataSourcesMap[k] = r
	}
}

func addNewSinceToResource(base, current *schema.ResourceJSON, curVersion string) {
	if base == nil {
		current.NewSince = curVersion
		return
	}

	current.NewSince = base.NewSince

	for k, curS := range current.Schema {
		if k == "id" {
			continue
		}
		var baseS *schema.SchemaJSON
		if b, ok := base.Schema[k]; ok {
			baseS = &b
		}
		addNewSinceToSchema(baseS, &curS, curVersion)
		current.Schema[k] = curS
	}
}

func addNewSinceToSchema(base, current *schema.SchemaJSON, curVersion string) {
	if base == nil {
		current.NewSince = curVersion
		return
	}

	current.NewSince = base.NewSince

	switch el := current.Elem.(type) {
	case *schema.SchemaJSON:
		if baseEl, ok := base.Elem.(schema.SchemaJSON); ok {
			addNewSinceToSchema(&baseEl, el, curVersion)
		}
	case *schema.ResourceJSON:
		if baseEl, ok := base.Elem.(schema.ResourceJSON); ok {
			addNewSinceToResource(&baseEl, el, curVersion)
		}
	}
}
