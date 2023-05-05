package backups

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

type CreateOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type Backup struct {
	Backup *BackupCreateOpts `json:"backup" required:"true"`
}

func (opts *Backup) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

type BackupCreateOpts struct {
	Name            string `json:"name" required:"true"`
	Description     string `json:"description,omitempty"`
	Instance        string `json:"instance,omitempty"`
	Cluster         string `json:"cluster,omitempty"`
	ContainerPrefix string `json:"container_prefix,omitempty"`
}

var dbBackupsAPIPath = "backups"

func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r GetResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(backupsURL(client, dbBackupsAPIPath), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes:      []int{202},
		JSONResponse: &r.Body,
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	var result *http.Response
	result, r.Err = client.Get(backupURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes:      []int{200},
		JSONResponse: &r.Body,
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(backupURL(client, id), &gophercloud.RequestOpts{})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
