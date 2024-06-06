package conntrackhelpers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type ConntrackHelperCreate struct {
	ConntrackHelper *CreateOpts `json:"dc_conntrack_helper"`
}

type CreateOpts struct {
	DCRouterID  string `json:"dc_router_id"`
	Protocol    string `json:"protocol"`
	Port        int    `json:"port"`
	Helper      string `json:"helper"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (opts *ConntrackHelperCreate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(conntrackHelpersURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(conntrackHelperURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}

type ConntrackHelperUpdate struct {
	ConntrackHelper *UpdateOpts `json:"dc_conntrack_helper"`
}

type UpdateOpts struct {
	Protocol    string `json:"protocol,omitempty"`
	Port        int    `json:"port,omitempty"`
	Helper      string `json:"helper,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (opts *ConntrackHelperUpdate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(conntrackHelperURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(conntrackHelperURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}
