package clusters

import (
	"encoding/base64"

	"github.com/gophercloud/gophercloud"
)

type Response struct {
	ID                      string      `json:"id" required:"true"`
	Name                    string      `json:"name" required:"true"`
	Status                  string      `json:"status"`
	Info                    ClusterInfo `json:"info"`
	S3BucketName            string      `json:"s3_bucket_name"`
	DockerRegistryID        string      `json:"docker_registry_id" required:"true"`
	HistoryServerURL        string      `json:"history_server_url"`
	ControlInstanceID       string      `json:"control_instance_id"`
	InactiveMin             int         `json:"inactive_min"`
	SuspendAfterInactiveMin int         `json:"suspend_after_inactive_min"`
	DeleteAfterInactiveMin  int         `json:"delete_after_inactive_min"`
}

type ClusterInfo struct {
	Name                 string              `json:"name" required:"true"`
	AvailabilityZone     string              `json:"availability_zone" required:"true"`
	NetworkID            string              `json:"network_id" required:"true"`
	SubnetID             string              `json:"subnet_id,omitempty"`
	NodeGroups           []NodeGroupResponse `json:"node_groups" required:"true"`
	ClusterMode          string              `json:"cluster_mode" required:"true"`
	IPPool               string              `json:"ip_pool" required:"true"`
	Keypair              string              `json:"keypair,omitempty"`
	DeleteAfterDelay     int                 `json:"delete_after_delay,omitempty"`
	SuspendAfterDelay    int                 `json:"suspend_after_delay,omitempty"`
	SparkConfiguration   string              `json:"spark_configuration,omitempty"`
	EnvironmentVariables string              `json:"environment_variables,omitempty"`
}

type NodeGroupResponse struct {
	NodeCount          int    `json:"node_count,omitempty"`
	FlavorID           string `json:"flavor_id" required:"true"`
	AutoscalingEnabled bool   `json:"autoscaling_enabled" required:"true"`
	MinNodes           int    `json:"min_nodes,omitempty"`
	MaxNodes           int    `json:"max_nodes,omitempty"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*Response, error) {
	var res *Response
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}

	if res.Info.SparkConfiguration != "" {
		sparkConfigurationDecoded, err := base64.StdEncoding.DecodeString(res.Info.SparkConfiguration)
		if err != nil {
			return nil, err
		}
		res.Info.SparkConfiguration = string(sparkConfigurationDecoded)
	}

	if res.Info.EnvironmentVariables != "" {
		environmentVariablesDecoded, err := base64.StdEncoding.DecodeString(res.Info.EnvironmentVariables)
		if err != nil {
			return nil, err
		}
		res.Info.EnvironmentVariables = string(environmentVariablesDecoded)
	}

	return res, nil
}

type CreateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}
