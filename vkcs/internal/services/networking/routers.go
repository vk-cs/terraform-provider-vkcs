package networking

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/pagination"
)

func ExtractRouterInto(r routers.GetResult, v interface{}) error {
	return r.ExtractIntoStructPtr(v, "router")
}

func ExtractRoutersInto(r pagination.Page, v interface{}) error {
	return r.(routers.RouterPage).Result.ExtractIntoSlicePtr(v, "routers")
}
