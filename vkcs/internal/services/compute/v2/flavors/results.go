package flavors

import (
	"encoding/json"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/pagination"
)

type FlavorMissingFields struct {
	ExtraSpecs map[string]interface{} `json:"extra_specs"`
}

// FlavorWithExtraSpecs needs for extract ExtraSpecs from flavors.FlavorPage
type FlavorWithExtraSpecs struct {
	flavors.Flavor
	FlavorMissingFields
}

func (f *FlavorWithExtraSpecs) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &f.Flavor); err != nil {
		return err
	}

	if err := json.Unmarshal(data, &f.FlavorMissingFields); err != nil {
		return err
	}

	return nil
}

func ExtractFlavorWithExtraSpecs(r pagination.Page) ([]FlavorWithExtraSpecs, error) {
	var s struct {
		Flavors []FlavorWithExtraSpecs `json:"flavors"`
	}
	err := (r.(flavors.FlavorPage)).ExtractInto(&s)

	return s.Flavors, err
}

func (f *FlavorWithExtraSpecs) ToFlavor() *flavors.Flavor {
	return &flavors.Flavor{
		ID:         f.ID,
		Disk:       f.Disk,
		RAM:        f.RAM,
		Name:       f.Name,
		RxTxFactor: f.RxTxFactor,
		Swap:       f.Swap,
		VCPUs:      f.VCPUs,
		IsPublic:   f.IsPublic,
		Ephemeral:  f.Ephemeral,
	}
}
