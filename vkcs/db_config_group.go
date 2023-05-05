package vkcs

import (
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
)

func extractDatabaseConfigGroupValues(rawValues map[string]interface{}, dsParameterTypes map[string]string) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	for name, value := range rawValues {
		if vType, ok := dsParameterTypes[name]; ok {
			switch vType {
			case "string":
				values[name] = value
			case "integer":
				if valueInt, err := strconv.Atoi(value.(string)); err == nil {
					values[name] = valueInt
				} else {
					return nil, fmt.Errorf("incorrect value type of parameter: %s", name)
				}
			case "boolean":
				if valueBool, err := strconv.ParseBool(value.(string)); err == nil {
					values[name] = valueBool
				} else {
					return nil, fmt.Errorf("incorrect value type of parameter: %s", name)
				}
			case "float":
				if valueFloat, err := strconv.ParseFloat(value.(string), 64); err == nil {
					values[name] = valueFloat
				} else {
					return nil, fmt.Errorf("incorrect value type of parameter: %s", name)
				}
			default:
				return nil, fmt.Errorf("unexpected type of parameter: %s", name)
			}
		} else {
			return nil, fmt.Errorf("incorrect parameter: %s", name)
		}
	}
	return values, nil
}

func flattenDatabaseConfigGroupValues(values map[string]interface{}) map[string]interface{} {
	rawValues := make(map[string]interface{})
	for name, value := range values {
		rawValues[name] = fmt.Sprintf("%v", value)
	}
	return rawValues
}

func getDSParameterTypesMap(dsParameters []datastores.Param) map[string]string {
	dsParameterTypes := make(map[string]string)
	for _, dsParameter := range dsParameters {
		dsParameterTypes[dsParameter.Name] = dsParameter.Type
	}
	return dsParameterTypes
}

func retrieveDatabaseConfigGroupValues(client *gophercloud.ServiceClient, datastore datastores.DatastoreShort, v map[string]interface{}) (map[string]interface{}, error) {
	dsParameters, err := datastores.ListParameters(client, datastore.Type, datastore.Version).Extract()
	if err != nil {
		return nil, fmt.Errorf("unable to determine vkcs_db_config_group parameter types")
	}
	dsParameterTypes := getDSParameterTypesMap(dsParameters)

	values, err := extractDatabaseConfigGroupValues(v, dsParameterTypes)
	if err != nil {
		return nil, fmt.Errorf("unable to determine vkcs_db_config_group values: %s", err)
	}
	return values, nil
}
