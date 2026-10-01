package keys

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Get(client *gophercloud.ServiceClient, name string) (r GetResult) {
	resp, err := client.Get(keyURL(client, name), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

type EncryptOpts struct {
	Plaintext string `json:"plaintext"`
}

func (opts EncryptOpts) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func Encrypt(client *gophercloud.ServiceClient, key string, opts EncryptOpts) (r EncryptResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(encryptURL(client, key), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}
