package acls

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/acls"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func GetContainerACL(client *gophercloud.ServiceClient, containerID string) acls.ACLResult {
	r := acls.GetContainerACL(client, containerID)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func GetSecretACL(client *gophercloud.ServiceClient, secretID string) acls.ACLResult {
	r := acls.GetSecretACL(client, secretID)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func SetContainerACL(client *gophercloud.ServiceClient, containerID string, opts acls.SetOptsBuilder) acls.ACLRefResult {
	r := acls.SetContainerACL(client, containerID, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func SetSecretACL(client *gophercloud.ServiceClient, secretID string, opts acls.SetOptsBuilder) acls.ACLRefResult {
	r := acls.SetSecretACL(client, secretID, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func UpdateContainerACL(client *gophercloud.ServiceClient, containerID string, opts acls.SetOptsBuilder) acls.ACLRefResult {
	r := acls.UpdateContainerACL(client, containerID, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func UpdateSecretACL(client *gophercloud.ServiceClient, secretID string, opts acls.SetOptsBuilder) acls.ACLRefResult {
	r := acls.UpdateSecretACL(client, secretID, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
