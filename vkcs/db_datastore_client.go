package vkcs

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

type datastoreResp struct {
	Name               string              `json:"name"`
	ID                 string              `json:"id"`
	MinimumCPU         int                 `json:"minimum_cpu"`
	MinimumRAM         int                 `json:"minimum_ram"`
	Versions           *[]datastoreVersion `json:"versions"`
	VolumeTypes        *[]string           `json:"volume_types"`
	ClusterVolumeTypes *[]string           `json:"cluster_volume_types"`
}

type datastoreRespOpts struct {
	Datastore *datastoreResp `json:"datastore"`
}

type datastoresRespOpts struct {
	Datastores *[]datastoreResp `json:"datastores"`
}

type datastoreVersion struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type commonDatastoreResult struct {
	gophercloud.Result
}

type getDatastoreResult struct {
	commonDatastoreResult
}

type getDatastoresResult struct {
	commonDatastoreResult
}

func (r getDatastoreResult) extract() (*datastoreResp, error) {
	var d *datastoreRespOpts
	if err := r.ExtractInto(&d); err != nil {
		return nil, err
	}
	return d.Datastore, nil
}

func (r getDatastoresResult) extract() (*[]datastoreResp, error) {
	var d *datastoresRespOpts
	if err := r.ExtractInto(&d); err != nil {
		return nil, err
	}
	return d.Datastores, nil
}

var datastoresAPIPath = "datastores"

func datastoreGet(client databaseClient, filter string) (r getDatastoreResult) {
	reqOpts := getDBRequestOpts(200)
	var result *http.Response
	result, r.Err = client.Get(getURL(client, datastoresAPIPath, filter), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func datastoresGet(client databaseClient) (r getDatastoresResult) {
	reqOpts := getDBRequestOpts(200)
	var result *http.Response
	result, r.Err = client.Get(baseURL(client, datastoresAPIPath), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
