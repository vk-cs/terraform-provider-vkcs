package flavors

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
	paginationutil "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/pagination"
)

type ListOptsBuilder interface {
	ToListQuery() (string, error)
}

type ListOpts struct {
	CpuModel    *string `q:"cpuModel"`
	CpuCoresMin *int64  `q:"cpuCoresMin"`
	CpuCoresMax *int64  `q:"cpuCoresMax"`
	RamSizeMin  *int64  `q:"ramSizeMin"`
	RamSizeMax  *int64  `q:"ramSizeMax"`
	SsdSizeMin  *int64  `q:"ssdSizeMin"`
	SsdSizeMax  *int64  `q:"ssdSizeMax"`
	HddSizeMin  *int64  `q:"hddSizeMin"`
	HddSizeMax  *int64  `q:"hddSizeMax"`
}

// ToListQuery builds request params.
func (opts *ListOpts) ToListQuery() (string, error) {
	u, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}

	return u.String(), err
}

// Get retrieves a baremetal flavor by ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(flavorURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// List returns a pager to iterate over all baremetal flavors.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := flavorsURL(client)
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return Page{
			TokenPageBase: paginationutil.TokenPageBase{PageResult: r},
		}
	})
}
