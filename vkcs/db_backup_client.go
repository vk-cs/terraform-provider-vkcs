package vkcs

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

type dbBackup struct {
	Backup *dbBackupCreateOpts `json:"backup" required:"true"`
}

type dbBackupCreateOpts struct {
	Name            string `json:"name" required:"true"`
	Description     string `json:"description,omitempty"`
	Instance        string `json:"instance,omitempty"`
	Cluster         string `json:"cluster,omitempty"`
	ContainerPrefix string `json:"container_prefix,omitempty"`
}

type dbBackupResp struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	LocationRef string     `json:"location_ref"`
	InstanceID  string     `json:"instance_id"`
	ClusterID   string     `json:"cluster_id"`
	Created     string     `json:"created"`
	Updated     string     `json:"updated"`
	Size        float64    `json:"size"`
	WalSize     float64    `json:"wal_size"`
	Status      string     `json:"status"`
	Datastore   *dataStore `json:"datastore"`
	Meta        string     `json:"meta"`
}

type dbBackupRespOpts struct {
	Backup *dbBackupResp `json:"backup"`
}

type commonDBBackupResult struct {
	gophercloud.Result
}

type getDBBackupResult struct {
	commonDBBackupResult
}

type deleteDBBackupResult struct {
	gophercloud.ErrResult
}

func (r getDBBackupResult) extract() (*dbBackupResp, error) {
	var b *dbBackupRespOpts
	if err := r.ExtractInto(&b); err != nil {
		return nil, err
	}
	return b.Backup, nil
}

func (opts *dbBackup) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

var dbBackupsAPIPath = "backups"

func dbBackupCreate(client databaseClient, opts optsBuilder) (r getDBBackupResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	reqOpts := getDBRequestOpts(202)
	result, r.Err = client.Post(baseURL(client, dbBackupsAPIPath), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func dbBackupGet(client databaseClient, id string) (r getDBBackupResult) {
	reqOpts := getDBRequestOpts(200)
	var result *http.Response
	result, r.Err = client.Get(getURL(client, dbBackupsAPIPath, id), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func dbBackupDelete(client databaseClient, id string) (r deleteDBBackupResult) {
	reqOpts := getDBRequestOpts()
	var result *http.Response
	result, r.Err = client.Delete(getURL(client, dbBackupsAPIPath, id), reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
