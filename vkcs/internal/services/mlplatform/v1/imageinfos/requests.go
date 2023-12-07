package imageinfos

import (
	"github.com/gophercloud/gophercloud"
)

func Get(client *gophercloud.ServiceClient, instanceType string) (r GetResult) {
	resp, err := client.Get(imageInfosURL(client, instanceType), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
