package routers

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type RouterCreate struct {
	Router *CreateOpts `json:"dc_router"`
}

type CreateOpts struct {
	AvailabilityZone string `json:"availability_zone,omitempty"`
	Flavor           string `json:"flavor,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
}

func (opts *RouterCreate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(routersURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(routerURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type RouterUpdate struct {
	Router *UpdateOpts `json:"dc_router"`
}

type UpdateOpts struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (opts *RouterUpdate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(routerURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(routerURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
