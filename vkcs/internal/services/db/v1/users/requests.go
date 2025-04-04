package users

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/databases"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// BatchCreateOpts is used to send request to create database users
type BatchCreateOpts struct {
	Users []CreateOpts `json:"users"`
}

// CreateOpts represents parameters of creation of database user
type CreateOpts struct {
	Name      string                 `json:"name" required:"true"`
	Password  string                 `json:"password" required:"true"`
	Databases []databases.CreateOpts `json:"databases,omitempty"`
	Host      string                 `json:"host,omitempty"`
}

// UpdateOpts represents parameters of update of database user
type UpdateOpts struct {
	User struct {
		Name     string `json:"name,omitempty"`
		Password string `json:"password,omitempty"`
		Host     string `json:"host,omitempty"`
	} `json:"user"`
}

// UpdateDatabasesOpts represents parameters of request to update users databases
type UpdateDatabasesOpts struct {
	Databases []map[string]string `json:"databases"`
}

// Map converts opts to a map (for a request body)
func (opts *BatchCreateOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *UpdateOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *UpdateDatabasesOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Create performs request to create database user
func Create(client *gophercloud.ServiceClient, id string, opts OptsBuilder, dbmsType string) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(usersURL(client, dbmsType, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List performs request to get list of database users
func List(client *gophercloud.ServiceClient, id string, dbmsType string) pagination.Pager {
	return pagination.NewPager(client, usersURL(client, dbmsType, id), func(r pagination.PageResult) pagination.Page {
		return Page{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// Update performs request to update database user
func Update(client *gophercloud.ServiceClient, id string, name string, opts OptsBuilder, dbmsType string) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(userURL(client, dbmsType, id, name), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// UpdateDatabases performs request to update database user databases
func UpdateDatabases(client *gophercloud.ServiceClient, id string, name string, opts OptsBuilder, dbmsType string) (r UpdateDatabasesResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(userDatabasesURL(client, dbmsType, id, name), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// DeleteDatabase performs request to delete database user
func DeleteDatabase(client *gophercloud.ServiceClient, id string, userName string, dbName string, dbmsType string) (r DeleteDatabaseResult) {
	resp, err := client.Delete(userDatabaseURL(client, dbmsType, id, userName, dbName), &gophercloud.RequestOpts{})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// Delete performs request to delete database user
func Delete(client *gophercloud.ServiceClient, id string, userName string, dbmsType string) (r DeleteResult) {
	resp, err := client.Delete(userURL(client, dbmsType, id, userName), &gophercloud.RequestOpts{})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}
