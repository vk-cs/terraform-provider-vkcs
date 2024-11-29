package secrets

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Get(client *gophercloud.ServiceClient, id string) secrets.GetResult {
	r := secrets.Get(client, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func GetPayload(client *gophercloud.ServiceClient, id string, opts secrets.GetPayloadOptsBuilder) secrets.PayloadResult {
	r := secrets.GetPayload(client, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Create(client *gophercloud.ServiceClient, opts secrets.CreateOptsBuilder) secrets.CreateResult {
	r := secrets.Create(client, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Delete(client *gophercloud.ServiceClient, id string) secrets.DeleteResult {
	r := secrets.Delete(client, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Update(client *gophercloud.ServiceClient, id string, opts secrets.UpdateOptsBuilder) secrets.UpdateResult {
	r := secrets.Update(client, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func GetMetadata(client *gophercloud.ServiceClient, secretID string) secrets.MetadataResult {
	r := secrets.GetMetadata(client, secretID)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func CreateMetadata(client *gophercloud.ServiceClient, secretID string, opts secrets.CreateMetadataOptsBuilder) secrets.MetadataCreateResult {
	r := secrets.CreateMetadata(client, secretID, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func GetMetadatum(client *gophercloud.ServiceClient, secretID string, key string) secrets.MetadatumResult {
	r := secrets.GetMetadatum(client, secretID, key)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func CreateMetadatum(client *gophercloud.ServiceClient, secretID string, opts secrets.CreateMetadatumOptsBuilder) secrets.MetadatumCreateResult {
	r := secrets.CreateMetadatum(client, secretID, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func UpdateMetadatum(client *gophercloud.ServiceClient, secretID string, opts secrets.UpdateMetadatumOptsBuilder) secrets.MetadatumResult {
	r := secrets.UpdateMetadatum(client, secretID, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func DeleteMetadatum(client *gophercloud.ServiceClient, secretID string, key string) secrets.MetadatumDeleteResult {
	r := secrets.DeleteMetadatum(client, secretID, key)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
