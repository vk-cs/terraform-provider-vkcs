package volumes

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Create(client *gophercloud.ServiceClient, opts volumes.CreateOptsBuilder) volumes.CreateResult {
	r := volumes.Create(client, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Delete(client *gophercloud.ServiceClient, id string, opts volumes.DeleteOptsBuilder) volumes.DeleteResult {
	r := volumes.Delete(client, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Get(client *gophercloud.ServiceClient, id string) volumes.GetResult {
	r := volumes.Get(client, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

type UpdateOpts struct {
	Name        *string           `json:"name,omitempty"`
	Description *string           `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata"`
}

func (opts UpdateOpts) ToVolumeUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "volume")
}

func Update(client *gophercloud.ServiceClient, id string, opts volumes.UpdateOptsBuilder) volumes.UpdateResult {
	r := volumes.Update(client, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
