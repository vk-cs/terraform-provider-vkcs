package keys

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func List(client *gophercloud.ServiceClient) (r ListResult) {
	resp, err := client.Get(keysURL(client), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Get(client *gophercloud.ServiceClient, name string) (r GetResult) {
	resp, err := client.Get(keyURL(client, name), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

type CreateOpts struct {
	Type string `json:"type"`
}

func (opts CreateOpts) Map() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func Create(client *gophercloud.ServiceClient, name string, opts CreateOpts) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(keyURL(client, name), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

type UpdateOpts struct {
	DeletionAllowed bool `json:"deletion_allowed"`
}

func (opts UpdateOpts) Map() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func Update(client *gophercloud.ServiceClient, name string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(keyConfigURL(client, name), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Delete(client *gophercloud.ServiceClient, name string) (r DeleteResult) {
	resp, err := client.Delete(keyURL(client, name), &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

type EncryptOpts struct {
	Plaintext string `json:"plaintext"`
}

func (opts EncryptOpts) Map() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func Encrypt(client *gophercloud.ServiceClient, key string, opts EncryptOpts) (r EncryptResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(encryptURL(client, key), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

type DecryptOpts struct {
	Ciphertext string `json:"ciphertext"`
}

func (opts DecryptOpts) Map() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func Decrypt(client *gophercloud.ServiceClient, key string, opts DecryptOpts) (r DecryptResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(decryptURL(client, key), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}
