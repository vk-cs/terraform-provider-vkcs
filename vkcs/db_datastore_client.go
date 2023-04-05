package vkcs

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// dataStoreVersion represents a version API resource
// Multiple versions belong to a dataStore.
type dataStoreVersion struct {
	ID    string
	Links []gophercloud.Link
	Name  string
}

// dataStore represents a Datastore API resource.
type dataStore struct {
	Name               string             `json:"name"`
	ID                 string             `json:"id"`
	MinimumCPU         int                `json:"minimum_cpu"`
	MinimumRAM         int                `json:"minimum_ram"`
	Versions           []dataStoreVersion `json:"versions"`
	VolumeTypes        []string           `json:"volume_types"`
	ClusterVolumeTypes []string           `json:"cluster_volume_types"`
}

// dataStoreGetResult represents the result of a dataStoreGet operation.
type dataStoreGetResult struct {
	gophercloud.Result
}

// dataStorePage represents a page of datastore resources.
type dataStorePage struct {
	pagination.SinglePageBase
}

// IsEmpty indicates whether a datastore collection is empty.
func (r dataStorePage) IsEmpty() (bool, error) {
	is, err := extractDatastores(r)
	return len(is) == 0, err
}

// extractDatastores retrieves a slice of dataStore structs from a paginated
// collection.
func extractDatastores(r pagination.Page) ([]dataStore, error) {
	var s struct {
		Datastores []dataStore `json:"datastores"`
	}
	err := (r.(dataStorePage)).ExtractInto(&s)
	return s.Datastores, err
}

// Extract retrieves a single dataStore struct from an operation result.
func (r dataStoreGetResult) Extract() (*dataStore, error) {
	var s struct {
		Datastore *dataStore `json:"datastore"`
	}
	err := r.ExtractInto(&s)
	return s.Datastore, err
}

var datastoresAPIPath = "datastores"

// dataStoreList will list all available datastores that instances can use.
func dataStoreList(client databaseClient) pagination.Pager {
	return pagination.NewPager(client.(*gophercloud.ServiceClient), datastoresURL(client, datastoresAPIPath),
		func(r pagination.PageResult) pagination.Page {
			return dataStorePage{pagination.SinglePageBase(r)}
		})
}

// dataStoreGet will retrieve the details of a specified datastore type.
func dataStoreGet(client databaseClient, datastoreID string) (r dataStoreGetResult) {
	resp, err := client.Get(datastoreURL(client, datastoresAPIPath, datastoreID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
