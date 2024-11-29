package ipportforwardings

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type IPPortForwardingCreate struct {
	IPPortForwarding *CreateOpts `json:"dc_ip_port_forwarding"`
}

type CreateOpts struct {
	DCInterfaceID string  `json:"dc_interface_id"`
	Protocol      string  `json:"protocol"`
	Source        *string `json:"source"`
	Destination   *string `json:"destination"`
	Port          *int64  `json:"port"`
	ToDestination string  `json:"to_destination"`
	ToPort        *int64  `json:"to_port"`
	Name          string  `json:"name,omitempty"`
	Description   string  `json:"description,omitempty"`
}

func (opts *IPPortForwardingCreate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ipPortForwardingsURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ipPortForwardingURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

type IPPortForwardingUpdate struct {
	IPPortForwarding *UpdateOpts `json:"dc_ip_port_forwarding"`
}

type UpdateOpts struct {
	Protocol      string  `json:"protocol,omitempty"`
	Source        *string `json:"source"`
	Destination   *string `json:"destination"`
	Port          *int64  `json:"port"`
	ToDestination string  `json:"to_destination,omitempty"`
	ToPort        *int64  `json:"to_port"`
	Name          string  `json:"name,omitempty"`
	Description   string  `json:"description,omitempty"`
}

func (opts *IPPortForwardingUpdate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ipPortForwardingURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ipPortForwardingURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}
