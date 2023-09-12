package api_options

import (
	"github.com/gophercloud/gophercloud"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

func Get(client *gophercloud.ServiceClient) (r GetResult) {
	resp, err := client.Get(apiOptionsURL(client), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
