package clusters

import (
	"encoding/base64"
	"net/http"

	"github.com/gophercloud/gophercloud"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type CreateOpts struct {
	Name                 string                `json:"name" required:"true"`
	AvailabilityZone     string                `json:"availability_zone" required:"true"`
	NetworkID            string                `json:"network_id" required:"true"`
	SubnetID             string                `json:"subnet_id,omitempty"`
	NodeGroups           []NodeGroupCreateOpts `json:"node_groups" required:"true"`
	ClusterMode          string                `json:"cluster_mode" required:"true"`
	Registry             RegistryCreateOpts    `json:"registry" required:"true"`
	IPPool               string                `json:"ip_pool" required:"true"`
	Keypair              string                `json:"keypair,omitempty"`
	DeleteAfterDelay     int                   `json:"delete_after_delay,omitempty"`
	SuspendAfterDelay    int                   `json:"suspend_after_delay,omitempty"`
	SparkConfiguration   string                `json:"spark_configuration,omitempty"`
	EnvironmentVariables string                `json:"environment_variables,omitempty"`
}

type NodeGroupCreateOpts struct {
	NodeCount          int    `json:"node_count,omitempty"`
	FlavorID           string `json:"flavor_id" required:"true"`
	AutoscalingEnabled bool   `json:"autoscaling_enabled" required:"true"`
	MinNodes           int    `json:"min_nodes,omitempty"`
	MaxNodes           int    `json:"max_nodes,omitempty"`
}

type RegistryCreateOpts struct {
	ExistingRegistryID string `json:"existing_registry_id" required:"true"`
}

func (opts *CreateOpts) Map() (map[string]interface{}, error) {
	if opts.SparkConfiguration != "" {
		opts.SparkConfiguration = base64.StdEncoding.EncodeToString([]byte(opts.SparkConfiguration))
	}
	if opts.EnvironmentVariables != "" {
		opts.EnvironmentVariables = base64.StdEncoding.EncodeToString([]byte(opts.EnvironmentVariables))
	}
	return gophercloud.BuildRequestBody(*opts, "")
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(clustersURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(clusterURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(clusterURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
