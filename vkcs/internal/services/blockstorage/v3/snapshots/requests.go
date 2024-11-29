package snapshots

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/snapshots"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Create(client *gophercloud.ServiceClient, opts snapshots.CreateOptsBuilder) snapshots.CreateResult {
	r := snapshots.Create(client, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Delete(client *gophercloud.ServiceClient, id string) snapshots.DeleteResult {
	r := snapshots.Delete(client, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Get(client *gophercloud.ServiceClient, id string) snapshots.GetResult {
	r := snapshots.Get(client, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Update(client *gophercloud.ServiceClient, id string, opts snapshots.UpdateOptsBuilder) snapshots.UpdateResult {
	r := snapshots.Update(client, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func UpdateMetadata(client *gophercloud.ServiceClient, id string, opts snapshots.UpdateMetadataOptsBuilder) snapshots.UpdateMetadataResult {
	r := snapshots.UpdateMetadata(client, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
