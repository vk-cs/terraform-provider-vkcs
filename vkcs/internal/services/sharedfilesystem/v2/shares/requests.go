package shares

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Create(client *gophercloud.ServiceClient, opts shares.CreateOptsBuilder) shares.CreateResult {
	r := shares.Create(client, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(client *gophercloud.ServiceClient, id string) shares.DeleteResult {
	r := shares.Delete(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Get(client *gophercloud.ServiceClient, id string) shares.GetResult {
	r := shares.Get(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func ListExportLocations(client *gophercloud.ServiceClient, id string) shares.ListExportLocationsResult {
	r := shares.ListExportLocations(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func GrantAccess(client *gophercloud.ServiceClient, id string, opts shares.GrantAccessOptsBuilder) shares.GrantAccessResult {
	r := shares.GrantAccess(client, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func RevokeAccess(client *gophercloud.ServiceClient, id string, opts shares.RevokeAccessOptsBuilder) shares.RevokeAccessResult {
	r := shares.RevokeAccess(client, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func ListAccessRights(client *gophercloud.ServiceClient, id string) shares.ListAccessRightsResult {
	r := shares.ListAccessRights(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Extend(client *gophercloud.ServiceClient, id string, opts shares.ExtendOptsBuilder) shares.ExtendResult {
	r := shares.Extend(client, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Shrink(client *gophercloud.ServiceClient, id string, opts shares.ShrinkOptsBuilder) shares.ShrinkResult {
	r := shares.Shrink(client, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Update(client *gophercloud.ServiceClient, id string, opts shares.UpdateOptsBuilder) shares.UpdateResult {
	r := shares.Update(client, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
