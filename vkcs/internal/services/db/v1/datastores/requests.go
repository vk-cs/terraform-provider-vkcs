package datastores

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// List will list all available datastores that instances can use.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, datastoresURL(client),
		func(r pagination.PageResult) pagination.Page {
			return Page{pagination.SinglePageBase(r)}
		})
}

// Get will retrieve the details of a specified datastore type.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(datastoreURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func ListCapabilities(client *gophercloud.ServiceClient, dsType string, versionID string) (r ListCapabilitiesResult) {
	url := datastoreCapabilitiesURL(client, dsType, versionID)
	resp, err := client.Get(url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// dataStoreListParameters will list all available configuration parameters
// for a specific version of a datastore.
func ListParameters(client *gophercloud.ServiceClient, dsType string, versionID string) (r ListParametersResult) {
	url := datastoreParametersURL(client, dsType, versionID)
	resp, err := client.Get(url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
