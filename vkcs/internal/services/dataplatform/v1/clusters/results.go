package clusters

import (
	"github.com/gophercloud/gophercloud"
)

// ClusterResp represents database cluster response
type Cluster struct {
	ID                string            `json:"id"`
	ClusterTemplateID string            `json:"cluster_template_id"`
	Configs           *ClusterConfig    `json:"configs"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	CreatedAt         string            `json:"created_at"`
	NetworkID         string            `json:"network_id"`
	PodGroups         []ClusterPodGroup `json:"pod_groups"`
	ProductName       string            `json:"product_name"`
	ProductType       string            `json:"product_type"`
	ProductVersion    string            `json:"product_version"`
	StackID           string            `json:"stack_id"`
	Status            string            `json:"status"`
	SubnetID          string            `json:"subnet_id"`
	Upgrades          []string          `json:"upgrades"`
	AvailabilityZone  string            `json:"availability_zone"`
	MultiAZ           bool              `json:"multi_az"`
	FloatingIPPool    string            `json:"floating_ip_pool"`
	Info              *ClusterInfo      `json:"info"`
}

type ClusterShortResp struct {
	ID string `json:"id"`
}

type commonClusterResult struct {
	gophercloud.Result
}

type commonShortClusterResult struct {
	gophercloud.Result
}

// CreateResult represents result of database cluster create
type CreateResult struct {
	commonShortClusterResult
}

type UpdateResult struct {
	commonShortClusterResult
}

// GetResult represents result of database cluster get
type GetResult struct {
	commonClusterResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}

// Extract is used to extract result into response struct
func (r commonClusterResult) Extract() (*Cluster, error) {
	var c *Cluster
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c, nil
}

// Extract is used to extract result into short response struct
func (r commonShortClusterResult) Extract() (*ClusterShortResp, error) {
	var c *ClusterShortResp
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c, nil
}
