package vrrps

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type VRRPCreate struct {
	VRRP *CreateOpts `json:"dc_vrrp"`
}

type CreateOpts struct {
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	GroupID        int    `json:"group_id"`
	NetworkID      string `json:"network_id"`
	SubnetID       string `json:"subnet_id,omitempty"`
	AdvertInterval int    `json:"advert_interval,omitempty"`
	Enabled        *bool  `json:"enabled,omitempty"`
}

func (opts *VRRPCreate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(vrrpsURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(vrrpURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type VRRPUpdate struct {
	VRRP *UpdateOpts `json:"dc_vrrp"`
}

type UpdateOpts struct {
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	GroupID        int    `json:"group_id,omitempty"`
	AdvertInterval int    `json:"advert_interval,omitempty"`
	Enabled        *bool  `json:"enabled,omitempty"`
}

func (opts *VRRPUpdate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(vrrpURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(vrrpURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
