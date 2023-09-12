package interfaces

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type InterfaceCreate struct {
	Interface *CreateOpts `json:"dc_interface"`
}

type CreateOpts struct {
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	DCRouterID         string `json:"dc_router_id"`
	NetworkID          string `json:"network_id"`
	SubnetID           string `json:"subnet_id,omitempty"`
	BGPAnnounceEnabled *bool  `json:"bgp_announce_enabled,omitempty"`
}

func (opts *InterfaceCreate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(interfacesURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(interfaceURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type InterfaceUpdate struct {
	Interface *UpdateOpts `json:"dc_interface"`
}

type UpdateOpts struct {
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	BGPAnnounceEnabled *bool  `json:"bgp_announce_enabled,omitempty"`
}

func (opts *InterfaceUpdate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(interfaceURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(interfaceURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
