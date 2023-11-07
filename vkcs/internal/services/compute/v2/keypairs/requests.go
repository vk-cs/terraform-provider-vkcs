package keypairs

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Create(client *gophercloud.ServiceClient, opts keypairs.CreateOptsBuilder) keypairs.CreateResult {
	r := keypairs.Create(client, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Get(client *gophercloud.ServiceClient, name string, opts keypairs.GetOptsBuilder) keypairs.GetResult {
	r := keypairs.Get(client, name, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(client *gophercloud.ServiceClient, name string, opts keypairs.DeleteOptsBuilder) keypairs.DeleteResult {
	r := keypairs.Delete(client, name, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
