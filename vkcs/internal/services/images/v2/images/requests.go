package images

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Create(client *gophercloud.ServiceClient, opts images.CreateOptsBuilder) images.CreateResult {
	r := images.Create(client, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(client *gophercloud.ServiceClient, id string) images.DeleteResult {
	r := images.Delete(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Get(client *gophercloud.ServiceClient, id string) images.GetResult {
	r := images.Get(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Update(client *gophercloud.ServiceClient, id string, opts images.UpdateOptsBuilder) images.UpdateResult {
	r := images.Update(client, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
