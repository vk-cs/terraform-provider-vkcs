package backups

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
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
	resp, err := client.Post(backupsURL(client, dbBackupsAPIPath), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes:      []int{202},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(backupURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes:      []int{200},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(backupURL(client, id), &gophercloud.RequestOpts{})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}
