package nodegroups

import (
	"github.com/gophercloud/gophercloud"
)

type (
	NodeGroupDetail struct {
		NgSpec NodeGroup `json:"node_group"`
	}

	NodeGroup struct {
		ID                   string            `json:"id"`
		Name                 string            `json:"name"`
		CreatedAt            string            `json:"created_at"`
		Zones                []string          `json:"zones" required:"true"`
		ScaleSpec            ScaleSpec         `json:"scale_spec" required:"true"`
		Labels               map[string]string `json:"labels,omitempty"`
		Taints               []Taint           `json:"taints,omitempty"`
		VMEngine             VMEngine          `json:"vm_engine" required:"true"`
		ParallelUpgradeChunk int               `json:"parallel_upgrade_chunk,omitempty"`
		ClusterID            string            `json:"cluster_id"`
		DiskType             DiskType          `json:"disk_type,omitempty"`
		UUID                 string            `json:"uuid"`
	}

	NodeGroupID struct {
		ID string `json:"node_group_id"`
	}

	CreateResult struct {
		gophercloud.Result
	}

	GetResult struct {
		gophercloud.Result
	}
)

// Extract is a function that accepts a result and extracts a node group.
func (r CreateResult) Extract() (string, error) {
	if r.Err != nil {
		return "", r.Err
	}

	var id NodeGroupID
	err := r.ExtractInto(&id)
	return id.ID, err
}

// Extract is a function that accepts a result and extracts a node group.
func (r GetResult) Extract() (*NodeGroup, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var ng NodeGroupDetail
	err := r.ExtractInto(&ng)
	return &ng.NgSpec, err
}
