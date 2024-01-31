package securitypolicies

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// CreateOpts contains options to create cluster security policy
type CreateOpts struct {
	ClusterID                string `json:"cluster_uuid" required:"true"`
	SecurityPolicyTemplateID string `json:"security_policy_uuid" required:"true"`
	PolicySettings           string `json:"policy_settings"`
	Namespace                string `json:"namespace"`
	Enabled                  bool   `json:"enabled"`
}

type UpdateOpts struct {
	PolicySettings string `json:"policy_settings"`
	Namespace      string `json:"namespace"`
	Enabled        bool   `json:"enabled"`
}

// Map builds request params.
func (opts *CreateOpts) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Map builds request params.
func (opts *UpdateOpts) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(securityPoliciesURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(securityPolicyURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}

func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(securityPolicyURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(securityPolicyURL(client, id), &gophercloud.RequestOpts{})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}
