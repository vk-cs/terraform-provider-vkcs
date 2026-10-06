package servergroups

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Create(client *gophercloud.ServiceClient, opts servergroups.CreateOptsBuilder) servergroups.CreateResult {
	r := servergroups.Create(client, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Get(client *gophercloud.ServiceClient, id string) servergroups.GetResult {
	r := servergroups.Get(client, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func List(client *gophercloud.ServiceClient, opts servergroups.ListOptsBuilder) pagination.Pager {
	return servergroups.List(client, opts)
}

func Delete(client *gophercloud.ServiceClient, id string) servergroups.DeleteResult {
	r := servergroups.Delete(client, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

// addRemoveMemberOpts is the request body payload for the add_member and
// remove_member server group actions.
type addRemoveMemberOpts struct {
	UUID string `json:"uuid"`
}

// AddMember adds a server to a server group using the add_member action
// (microversion 2.42).
func AddMember(client *gophercloud.ServiceClient, groupID, instanceID string) (r gophercloud.ErrResult) {
	body := map[string]interface{}{
		"add_member": addRemoveMemberOpts{UUID: instanceID},
	}

	resp, err := client.Post(client.ServiceURL("os-server-groups", groupID, "action"), body, nil, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// RemoveMember removes a server from a server group using the remove_member
// action (microversion 2.42).
func RemoveMember(client *gophercloud.ServiceClient, groupID, instanceID string) (r gophercloud.ErrResult) {
	body := map[string]interface{}{
		"remove_member": addRemoveMemberOpts{UUID: instanceID},
	}

	resp, err := client.Post(client.ServiceURL("os-server-groups", groupID, "action"), body, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}
