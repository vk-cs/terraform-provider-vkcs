package subnets

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/gophercloud/pagination"
)

func ExtractSubnetInto(r subnets.GetResult, v interface{}) error {
	return r.ExtractIntoStructPtr(v, "subnet")
}

func ExtractSubnetsInto(r pagination.Page, v interface{}) error {
	return r.(subnets.SubnetPage).ExtractIntoSlicePtr(v, "subnets")
}
