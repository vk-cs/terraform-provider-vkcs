package securityservices

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/securityservices"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Create(client *gophercloud.ServiceClient, opts securityservices.CreateOptsBuilder) securityservices.CreateResult {
	r := securityservices.Create(client, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(client *gophercloud.ServiceClient, id string) securityservices.DeleteResult {
	r := securityservices.Delete(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Get(client *gophercloud.ServiceClient, id string) securityservices.GetResult {
	r := securityservices.Get(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Update(client *gophercloud.ServiceClient, id string, opts securityservices.UpdateOptsBuilder) securityservices.UpdateResult {
	r := securityservices.Update(client, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
