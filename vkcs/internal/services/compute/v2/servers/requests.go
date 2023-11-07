package servers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Create(client *gophercloud.ServiceClient, opts servers.CreateOptsBuilder) servers.CreateResult {
	r := servers.Create(client, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(client *gophercloud.ServiceClient, id string) servers.DeleteResult {
	r := servers.Delete(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func ForceDelete(client *gophercloud.ServiceClient, id string) servers.ActionResult {
	r := servers.ForceDelete(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Get(client *gophercloud.ServiceClient, id string) servers.GetResult {
	r := servers.Get(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Update(client *gophercloud.ServiceClient, id string, opts servers.UpdateOptsBuilder) servers.UpdateResult {
	r := servers.Update(client, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func ChangeAdminPassword(client *gophercloud.ServiceClient, id, newPassword string) servers.ActionResult {
	r := servers.ChangeAdminPassword(client, id, newPassword)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Resize(client *gophercloud.ServiceClient, id string, opts servers.ResizeOptsBuilder) servers.ActionResult {
	r := servers.Resize(client, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func ConfirmResize(client *gophercloud.ServiceClient, id string) servers.ActionResult {
	r := servers.ConfirmResize(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Metadata(client *gophercloud.ServiceClient, id string) servers.GetMetadataResult {
	r := servers.Metadata(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func UpdateMetadata(client *gophercloud.ServiceClient, id string, opts servers.UpdateMetadataOptsBuilder) servers.UpdateMetadataResult {
	r := servers.UpdateMetadata(client, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func DeleteMetadatum(client *gophercloud.ServiceClient, id, key string) servers.DeleteMetadatumResult {
	r := servers.DeleteMetadatum(client, id, key)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
