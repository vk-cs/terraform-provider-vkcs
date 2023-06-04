package providers

import (
	"github.com/gophercloud/gophercloud"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// List returns information about backup providers
func List(client *gophercloud.ServiceClient) (r ListResult) {
	resp, err := client.Get(providersURL(client), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
