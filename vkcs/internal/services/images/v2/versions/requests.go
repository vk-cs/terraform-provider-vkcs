package versions

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Get(client *gophercloud.ServiceClient) VersionsResult {
	var r VersionsResult
	_, r.Err = client.Get(versionsURL(client), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 300},
	})
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
