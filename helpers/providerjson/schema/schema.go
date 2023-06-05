package schema

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	SchemaTypeSet    = "TypeSet"
	SchemaTypeList   = "TypeList"
	SchemaTypeMap    = "TypeMap"
	SchemaTypeInt    = "TypeInt"
	SchemaTypeString = "TypeString"
	SchemaTypeBool   = "TypeBool"
	SchemaTypeFloat  = "TypeFloat"
	SchemaTypeNumber = "TypeNumber"
	SchemaTypeObject = "TypeObject"
)

type ProviderJSON struct {
	SDKProvider *sdkschema.Provider
	Provider    provider.Provider
}

type SchemaJSON struct {
	Type        string      `json:"type,omitempty"`
	ConfigMode  string      `json:"config_mode,omitempty"`
	Optional    bool        `json:"optional,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description,omitempty"`
	Computed    bool        `json:"computed,omitempty"`
	ForceNew    bool        `json:"force_new,omitempty"`
	Elem        interface{} `json:"elem,omitempty"`
	MaxItems    int         `json:"max_items,omitempty"`
	MinItems    int         `json:"min_items,omitempty"`
	NewSince    string      `json:"new_since,omitempty"`
	Deprecated  string      `json:"deprecated,omitempty"`
}

func (b *SchemaJSON) UnmarshalJSON(body []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	b.Type, _ = m["type"].(string)
	b.ConfigMode, _ = m["config_mode"].(string)
	b.Optional, _ = m["optional"].(bool)
	b.Required, _ = m["required"].(bool)
	b.Description, _ = m["description"].(string)
	b.Computed, _ = m["computed"].(bool)
	b.ForceNew, _ = m["force_new"].(bool)
	if max, ok := m["max_items"].(float64); ok {
		b.MaxItems = int(max)
	}
	if min, ok := m["min_items"].(float64); ok {
		b.MinItems = int(min)
	}

	if def, ok := m["default"]; ok && def != nil {
		switch def.(type) {
		case string:
			b.Default = def
		case bool:
			b.Default = def
		case int:
			b.Default = def
		case float32:
			b.Default = def
		case float64:
			b.Default = def
		}
	}

	if e, ok := m["elem"]; ok && e != nil {
		elem := e.(map[string]interface{})
		if _, ok := elem["schema"]; ok {
			elemJSON, _ := json.Marshal(m["elem"])
			var e ResourceJSON
			_ = json.Unmarshal(elemJSON, &e)
			b.Elem = e
		} else if t, ok := elem["type"]; ok {
			b.Elem = t.(string)
		}
	}

	b.NewSince, _ = m["new_since"].(string)
	b.Deprecated, _ = m["deprecated"].(string)

	return nil
}

type ResourceJSON struct {
	Schema      map[string]SchemaJSON `json:"schema"`
	Timeouts    *ResourceTimeoutJSON  `json:"timeouts,omitempty"`
	Description string                `json:"description,omitempty"`
	NewSince    string                `json:"new_since,omitempty"`
	Deprecated  string                `json:"deprecated,omitempty"`
}

type ResourceTimeoutJSON struct {
	Create int `json:"create,omitempty"`
	Read   int `json:"read,omitempty"`
	Delete int `json:"delete,omitempty"`
	Update int `json:"update,omitempty"`
}

type ProviderSchemaJSON struct {
	Schema         map[string]SchemaJSON   `json:"schema"`
	ResourcesMap   map[string]ResourceJSON `json:"resources,omitempty"`
	DataSourcesMap map[string]ResourceJSON `json:"data_sources,omitempty"`
}

type ProviderWrapper struct {
	ProviderName    string              `json:"provider_name"`
	ProviderVersion string              `json:"provider_version"`
	SchemaVersion   string              `json:"schema_version"`
	ProviderSchema  *ProviderSchemaJSON `json:"provider_schema,omitempty"`
}
