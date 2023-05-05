package configgroups

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type ConfigGroup struct {
	Configuration *CreateOpts `json:"configuration"`
}

type CreateOpts struct {
	Datastore   *datastores.DatastoreShort `json:"datastore"`
	Name        string                     `json:"name"`
	Values      map[string]interface{}     `json:"values"`
	Description string                     `json:"description,omitempty"`
}

type UpdateOpts struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Values      map[string]interface{} `json:"values,omitempty"`
}

type UpdateOpt struct {
	Configuration *UpdateOpts `json:"configuration"`
}

func (opts *ConfigGroup) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

func (opts *UpdateOpt) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(configGroupsURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	var result *http.Response
	result, r.Err = client.Get(configGroupURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Put(configGroupURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(configGroupURL(client, id), &gophercloud.RequestOpts{})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
