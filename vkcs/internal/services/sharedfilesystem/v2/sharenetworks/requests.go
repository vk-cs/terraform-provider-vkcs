package sharenetworks

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/sharenetworks"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Create(client *gophercloud.ServiceClient, opts sharenetworks.CreateOptsBuilder) sharenetworks.CreateResult {
	r := sharenetworks.Create(client, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Delete(client *gophercloud.ServiceClient, id string) sharenetworks.DeleteResult {
	r := sharenetworks.Delete(client, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Get(client *gophercloud.ServiceClient, id string) sharenetworks.GetResult {
	r := sharenetworks.Get(client, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Update(client *gophercloud.ServiceClient, id string, opts sharenetworks.UpdateOptsBuilder) sharenetworks.UpdateResult {
	r := sharenetworks.Update(client, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func AddSecurityService(client *gophercloud.ServiceClient, id string, opts sharenetworks.AddSecurityServiceOptsBuilder) sharenetworks.UpdateResult {
	r := sharenetworks.AddSecurityService(client, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func RemoveSecurityService(client *gophercloud.ServiceClient, id string, opts sharenetworks.RemoveSecurityServiceOptsBuilder) sharenetworks.UpdateResult {
	r := sharenetworks.RemoveSecurityService(client, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
