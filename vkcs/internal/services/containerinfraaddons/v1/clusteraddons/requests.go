package clusteraddons

import (
	"encoding/base64"

	"github.com/gophercloud/gophercloud"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type UpgradeClusterAddonOpts struct {
	NewAddonID string `json:"new_addon_id"`
	Values     string `json:"values"`
}

func (opts UpgradeClusterAddonOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.Values != "" {
		valuesBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(opts.Values)))
		base64.StdEncoding.Encode(valuesBase64, []byte(opts.Values))
		b["values"] = string(valuesBase64)
	}

	return b, nil
}

type UpdateClusterAddonOpts struct {
	Values string `json:"values"`
}

func (opts UpdateClusterAddonOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.Values != "" {
		valuesBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(opts.Values)))
		base64.StdEncoding.Encode(valuesBase64, []byte(opts.Values))
		b["values"] = string(valuesBase64)
	}

	return b, nil
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(clusterAddonURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Upgrade(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(clusterAddonURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(clusterAddonURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(clusterAddonURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
