package dnsnames

import (
	"github.com/gophercloud/gophercloud"
)

func Get(client *gophercloud.ServiceClient) (r GetResult) {
	resp, err := client.Get(dnsNamesURL(client), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
