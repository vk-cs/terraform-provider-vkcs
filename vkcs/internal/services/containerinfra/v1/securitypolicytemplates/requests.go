package securitypolicytemplates

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func List(client *gophercloud.ServiceClient) (r securityPolicyTemplatesResult) {
	resp, err := client.Get(securityPolicyTemplatesURL(client), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}
