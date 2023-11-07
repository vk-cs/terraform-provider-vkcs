package apioptions

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

func Get(client *gophercloud.ServiceClient) (r GetResult) {
	resp, err := client.Get(apiOptionsURL(client), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}
