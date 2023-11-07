package pools

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Create(c *gophercloud.ServiceClient, opts pools.CreateOptsBuilder) pools.CreateResult {
	r := pools.Create(c, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Get(c *gophercloud.ServiceClient, id string) pools.GetResult {
	r := pools.Get(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Update(c *gophercloud.ServiceClient, id string, opts pools.UpdateOptsBuilder) pools.UpdateResult {
	r := pools.Update(c, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(c *gophercloud.ServiceClient, id string) pools.DeleteResult {
	r := pools.Delete(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func CreateMember(c *gophercloud.ServiceClient, poolID string, opts pools.CreateMemberOptsBuilder) pools.CreateMemberResult {
	r := pools.CreateMember(c, poolID, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func GetMember(c *gophercloud.ServiceClient, poolID string, memberID string) pools.GetMemberResult {
	r := pools.GetMember(c, poolID, memberID)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func UpdateMember(c *gophercloud.ServiceClient, poolID string, memberID string, opts pools.UpdateMemberOptsBuilder) pools.UpdateMemberResult {
	r := pools.UpdateMember(c, poolID, memberID, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func BatchUpdateMembers(c *gophercloud.ServiceClient, poolID string, opts []pools.BatchUpdateMemberOpts) pools.UpdateMembersResult {
	r := pools.BatchUpdateMembers(c, poolID, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func DeleteMember(c *gophercloud.ServiceClient, poolID string, memberID string) pools.DeleteMemberResult {
	r := pools.DeleteMember(c, poolID, memberID)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
