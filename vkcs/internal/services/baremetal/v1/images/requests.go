package images

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
	paginationutil "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/pagination"
)

// Get retrieves a baremetal image by ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(imageURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// List returns a pager to iterate over all baremetal images.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	url := imagesURL(client)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return Page{
			TokenPageBase: paginationutil.TokenPageBase{PageResult: r},
		}
	})
}
