package vkcs

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

type dbConfigGroup struct {
	Configuration *dbConfigGroupCreateOpts `json:"configuration"`
}

type dbConfigGroupCreateOpts struct {
	Datastore   *dataStore             `json:"datastore"`
	Name        string                 `json:"name"`
	Values      map[string]interface{} `json:"values"`
	Description string                 `json:"description,omitempty"`
}

type dbConfigGroupResp struct {
	ID                   string                 `json:"id"`
	DatastoreName        string                 `json:"datastore_name"`
	DatastoreVersionName string                 `json:"datastore_version_name"`
	Name                 string                 `json:"name"`
	Values               map[string]interface{} `json:"values"`
	Updated              string                 `json:"updated"`
	Created              string                 `json:"created"`
	Description          string                 `json:"description"`
}

type dbConfigGroupRespOpts struct {
	Configuration *dbConfigGroupResp `json:"configuration"`
}

type dbConfigGroupUpdateOpts struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Values      map[string]interface{} `json:"values,omitempty"`
}

type dbConfigGroupUpdateOpt struct {
	Configuration *dbConfigGroupUpdateOpts `json:"configuration"`
}

type dbDatastoreParametersResp struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type dbDatastoreParametersRespOpts struct {
	ConfigurationParameters []dbDatastoreParametersResp `json:"configuration-parameters"`
}

type commonDBConfigGroupResult struct {
	gophercloud.Result
}

type getDBConfigGroupResult struct {
	commonDBConfigGroupResult
}

type getDBDatastoreParametersResult struct {
	commonDBConfigGroupResult
}

type updateDBConfigGroupResult struct {
	gophercloud.ErrResult
}

type deleteDBConfigGroupResult struct {
	gophercloud.ErrResult
}

func (r getDBConfigGroupResult) extract() (*dbConfigGroupResp, error) {
	var c *dbConfigGroupRespOpts
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c.Configuration, nil
}

func (r getDBDatastoreParametersResult) extract() ([]dbDatastoreParametersResp, error) {
	var d *dbDatastoreParametersRespOpts
	if err := r.ExtractInto(&d); err != nil {
		return nil, err
	}
	return d.ConfigurationParameters, nil
}

func (opts *dbConfigGroup) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

func (opts *dbConfigGroupUpdateOpt) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

var dbConfigGroupsAPIPath = "configurations"
var dbDatastoresAPIPath = "datastores"

func dbConfigGroupCreate(client databaseClient, opts optsBuilder) (r getDBConfigGroupResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	reqOpts := getDBRequestOpts(200)
	result, r.Err = client.Post(baseURL(client, dbConfigGroupsAPIPath), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func dbConfigGroupGet(client databaseClient, id string) (r getDBConfigGroupResult) {
	reqOpts := getDBRequestOpts(200)
	var result *http.Response
	result, r.Err = client.Get(getURL(client, dbConfigGroupsAPIPath, id), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func dbConfigGroupUpdate(client databaseClient, id string, opts optsBuilder) (r updateDBConfigGroupResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getDBRequestOpts(202)
	var result *http.Response
	result, r.Err = client.Put(getURL(client, dbConfigGroupsAPIPath, id), b, nil, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func dbConfigGroupDelete(client databaseClient, id string) (r deleteDBConfigGroupResult) {
	reqOpts := getDBRequestOpts()
	var result *http.Response
	result, r.Err = client.Delete(getURL(client, dbConfigGroupsAPIPath, id), reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func dbDatastoreParametersGet(client databaseClient, dsType string, dsVersion string) (r getDBDatastoreParametersResult) {
	reqOpts := getDBRequestOpts(200)
	var result *http.Response
	result, r.Err = client.Get(datastoreParametersURL(client, dbDatastoresAPIPath, dsType, dsVersion), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
