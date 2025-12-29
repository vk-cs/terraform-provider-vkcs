package flavors

import (
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/pagination"
)

type FlavorExtExtraSpecs struct {
	ExtraSpecs map[string]interface{} `json:"extra_specs"`
}

func ExtractFlavorsInto(r pagination.Page, to interface{}) error {
	return (r.(flavors.FlavorPage)).ExtractIntoSlicePtr(to, "flavors")
}
