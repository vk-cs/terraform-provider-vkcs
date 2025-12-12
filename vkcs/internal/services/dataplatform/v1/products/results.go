package products

import (
	"github.com/gophercloud/gophercloud"
)

type Products struct {
	Products []Product `json:"products"`
}

type Product struct {
	ProductName    string         `json:"product_name"`
	ProductVersion string         `json:"product_version"`
	Configs        *ProductConfig `json:"configs"`
}

type ProductConfig struct {
	Settings    []ProductConfigSetting    `json:"settings"`
	Connections []ProductConfigConnection `json:"connections"`
	UserRoles   []ProductConfigUserRole   `json:"user_roles"`
	Crontabs    []ProductConfigCrontabs   `json:"crontabs"`
}

type ProductConfigSetting struct {
	Alias           string   `json:"alias"`
	DefaultValue    string   `json:"default_value"`
	RegExp          string   `json:"regexp"`
	StringVariation []string `json:"string_variation"`
	IsRequired      bool     `json:"is_require"`
	IsSensitive     bool     `json:"is_sensitive"`
}

type ProductConfigConnection struct {
	Plug          string                           `json:"plug"`
	IsRequired    bool                             `json:"is_required"`
	Position      int                              `json:"position"`
	RequiredGroup string                           `json:"required_group"`
	Settings      []ProductConfigConnectionSetting `json:"settings"`
}

type ProductConfigConnectionSetting struct {
	Alias           string   `json:"alias"`
	DefaultValue    string   `json:"default_value"`
	RegExp          string   `json:"regexp"`
	StringVariation []string `json:"string_variation"`
	IsRequired      bool     `json:"is_require"`
	IsSensitive     bool     `json:"is_sensitive"`
}

type ProductConfigUserRole struct {
	Name string `json:"name"`
}

type ProductConfigCrontabs struct {
	Name     string                          `json:"name"`
	Start    string                          `json:"start"`
	Required bool                            `json:"required"`
	Settings []ProductConfigCrontabsSettings `json:"settings"`
}

type ProductConfigCrontabsSettings struct {
	Alias           string   `json:"alias"`
	DefaultValue    string   `json:"default_value"`
	RegExp          string   `json:"regexp"`
	StringVariation []string `json:"string_variation"`
	IsRequired      bool     `json:"is_require"`
	IsSensitive     bool     `json:"is_sensitive"`
}

type commonProductResult struct {
	gophercloud.Result
}

// ListResult represents result of dataplatform products get
type ListResult struct {
	commonProductResult
}

// Extract is used to extract result into response struct
func (r commonProductResult) Extract() (*Products, error) {
	var p *Products
	if err := r.ExtractInto(&p); err != nil {
		return nil, err
	}
	return p, nil
}
