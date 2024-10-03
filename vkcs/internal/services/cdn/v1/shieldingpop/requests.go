package shieldingpop

import (
	"github.com/gophercloud/gophercloud"
)

// List returns a list of CDN origin shielding points of precense available in a project.
func List(client *gophercloud.ServiceClient, projectID string) (r ListResult) {
	resp, err := client.Get(shieldingPopsURL(client, projectID), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
