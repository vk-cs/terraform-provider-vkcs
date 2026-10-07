package secrets

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func List(client *gophercloud.ServiceClient) (r ListResult) {
	resp, err := client.Get(secretsMetadataURL(client), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Get(client *gophercloud.ServiceClient, path string) (r GetResult) {
	resp, err := client.Get(secretDataURL(client, path), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func GetDeleteProtection(client *gophercloud.ServiceClient, path string) (r GetDeleteProtectionResult) {
	resp, err := client.Get(secretDeleteProtectionURL(client, path), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

type CreateOrUpdateOpts struct {
	Data map[string]any `json:"data"`
}

func (opts CreateOrUpdateOpts) Map() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func CreateOrUpdate(client *gophercloud.ServiceClient, path string, opts CreateOrUpdateOpts) (r CreateOrUpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(secretDataURL(client, path), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

type SetDeleteProtectionOpts struct {
	DeleteProtection bool `json:"delete_protection"`
}

func (opts SetDeleteProtectionOpts) Map() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func SetDeleteProtection(client *gophercloud.ServiceClient, path string, opts SetDeleteProtectionOpts) (r SetDeleteProtectionResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	// This operation can return an empty success body; no JSON response is needed.
	resp, err := client.Put(secretDeleteProtectionURL(client, path), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Delete(client *gophercloud.ServiceClient, path string) (r DeleteResult) {
	resp, err := client.Delete(secretMetadataURL(client, path), &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}
