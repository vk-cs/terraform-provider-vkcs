package provider

import (
	"fmt"

	jsonschema "github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
	"github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/transform/datasource"
	"github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/transform/resource"
	"github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/transform/sdk"
)

func WrappedProviderFromRaw(input *jsonschema.ProviderJSON, base, wrapper *jsonschema.ProviderWrapper) (*jsonschema.ProviderWrapper, error) {
	schema, err := ProviderFromRaw(input)
	if err != nil {
		return nil, err
	}

	if base != nil {
		addNewSince(base.ProviderSchema, schema, wrapper.ProviderVersion)
	}

	wrapper.ProviderSchema = schema

	return wrapper, nil
}

func ProviderFromRaw(input *jsonschema.ProviderJSON) (*jsonschema.ProviderSchemaJSON, error) {
	if input == nil {
		return nil, fmt.Errorf("provider was nil converting from raw")
	}

	result := &jsonschema.ProviderSchemaJSON{}

	providerSchema := make(map[string]jsonschema.SchemaJSON)
	resourceSchemas := make(map[string]jsonschema.ResourceJSON)
	dataSourceSchemas := make(map[string]jsonschema.ResourceJSON)

	for k, v := range input.SDKSchema() {
		providerSchema[k] = sdk.SchemaFromRaw(v)
	}

	for k, v := range input.SDKResourcesMap() {
		resource, err := sdk.ResourceFromRaw(v)
		if err != nil {
			return nil, err
		}
		resourceSchemas[k] = *resource
	}

	for k, v := range input.SDKDataSourcesMap() {
		dataSource, err := sdk.ResourceFromRaw(v)
		if err != nil {
			return nil, err
		}
		dataSourceSchemas[k] = *dataSource
	}

	for k, v := range input.ResourcesMap() {
		resource, err := resource.ResourceFromRaw(v)
		if err != nil {
			return nil, err
		}
		resourceSchemas[k] = *resource
	}

	for k, v := range input.DataSourcesMap() {
		dataSource, err := datasource.ResourceFromRaw(v)
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
