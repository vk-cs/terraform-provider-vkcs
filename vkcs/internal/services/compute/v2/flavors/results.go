package flavors

import (
	"encoding/json"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/pagination"
)

type FlavorExtraFields struct {
	ExtraSpecs map[string]interface{} `json:"extra_specs"`
}

// FlavorWithExtraFields needs for extract FlavorExtraFields from flavors.FlavorPage
type FlavorWithExtraFields struct {
	flavors.Flavor
	FlavorExtraFields
}

func (f *FlavorWithExtraFields) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &f.Flavor); err != nil {
		return err
	}

	if err := json.Unmarshal(data, &f.FlavorExtraFields); err != nil {
		return err
	}

	return nil
}

func ExtractFlavorWithExtraSpecs(r pagination.Page) ([]FlavorWithExtraFields, error) {
	var s struct {
		Flavors []FlavorWithExtraFields `json:"flavors"`
	}
	err := (r.(flavors.FlavorPage)).ExtractInto(&s)

	return s.Flavors, err
}
