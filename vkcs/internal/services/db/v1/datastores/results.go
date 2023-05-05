package datastores

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Datastore represents a Datastore API resource.
type Datastore struct {
	Name               string    `json:"name"`
	ID                 string    `json:"id"`
	MinimumCPU         int       `json:"minimum_cpu"`
	MinimumRAM         int       `json:"minimum_ram"`
	Versions           []Version `json:"versions"`
	VolumeTypes        []string  `json:"volume_types"`
	ClusterVolumeTypes []string  `json:"cluster_volume_types"`
}

type DatastoreShort struct {
	Type    string `json:"type" required:"true"`
	Version string `json:"version" required:"true"`
}

type ParametersResp struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ParametersRespOpts struct {
	ConfigurationParameters []ParametersRespOpts `json:"configuration-parameters"`
}

// Version represents a version API resource
// Multiple versions belong to a dataStore.
type Version struct {
	ID    string
	Links []gophercloud.Link
	Name  string
}

// CapabilityParam represents a parameter of a Datastore capability.
type CapabilityParam struct {
	Required     bool        `json:"required"`
	Type         string      `json:"type"`
	ElementType  string      `json:"element_type"`
	EnumValues   []string    `json:"enum_values"`
	DefaultValue interface{} `json:"default_value"` // workaround since there's bug in API response
	MinValue     float64     `json:"min"`
	MaxValue     float64     `json:"max"`
	Regex        string      `json:"regex"`
	Masked       bool        `json:"masked"`
}

// Capability represents a Datastore capability.
type Capability struct {
	Name                   string                      `json:"name"`
	Description            string                      `json:"description"`
	Status                 string                      `json:"status"`
	Params                 map[string]*CapabilityParam `json:"params"`
	ShouldBeOnMaster       bool                        `json:"should_be_on_master"`
	AllowUpgradeFromBackup bool                        `json:"allow_upgrade_from_backup"`
	AllowMajorUpgrade      bool                        `json:"allow_major_upgrade"`
}

// Capabilities represents a object containing all datastore
// capabilities.
type Capabilities struct {
	Capabilities []Capability `json:"capabilities"`
}

// Param represents a configuration parameter supported by a datastore
type Param struct {
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	MinValue        float64 `json:"min"`
	MaxValue        float64 `json:"max"`
	RestartRequried bool    `json:"restart_required"`
}

// Params represents a object containing all datastore
// configuration parameters.
type Params struct {
	Params []Param `json:"configuration-parameters"`
}

// GetResult represents the result of a dataStoreGet operation.
type GetResult struct {
	gophercloud.Result
}

// Page represents a page of datastore resources.
type Page struct {
	pagination.SinglePageBase
}

// ListCapabilitiesResult represents the result of
// ListCapabilities operation.
type ListCapabilitiesResult struct {
	gophercloud.Result
}

// Extract retrieves a single dataStore struct from an operation result.
func (r GetResult) Extract() (*Datastore, error) {
	var s struct {
		Datastore *Datastore `json:"datastore"`
	}
	err := r.ExtractInto(&s)
	return s.Datastore, err
}

// IsEmpty indicates whether a Datastore collection is empty.
func (r Page) IsEmpty() (bool, error) {
	is, err := ExtractDatastores(r)
	return len(is) == 0, err
}

// ExtractDatastores retrieves a slice of dataStore structs from a paginated
// collection.
func ExtractDatastores(r pagination.Page) ([]Datastore, error) {
	var s struct {
		Datastores []Datastore `json:"datastores"`
	}
	err := (r.(Page)).ExtractInto(&s)
	return s.Datastores, err
}

// ListParametersResult represents the result of
// ListParameters operation.
type ListParametersResult struct {
	gophercloud.Result
}

func (r ListCapabilitiesResult) Extract() ([]Capability, error) {
	var dsCapabilities Capabilities
	err := r.ExtractInto(&dsCapabilities)
	return dsCapabilities.Capabilities, err
}

func (r ListParametersResult) Extract() ([]Param, error) {
	var dsParams Params
	err := r.ExtractInto(&dsParams)
	return dsParams.Params, err
}
