package templates

import (
	"github.com/gophercloud/gophercloud"
)

type ClusterTemplates struct {
	ClusterTemplates []ClusterTemplate `json:"cluster_templates"`
}

type ClusterTemplate struct {
	ID             string                    `json:"id"`
	Name           string                    `json:"name"`
	ProductName    string                    `json:"product_name"`
	ProductVersion string                    `json:"product_version"`
	PodGroups      []ClusterTemplatePodgroup `json:"pod_groups"`
}

type ClusterTemplatePodgroup struct {
	ID       string                                   `json:"id"`
	Name     string                                   `json:"name"`
	Resource ClusterTemplatePodgroupResource          `json:"resource"`
	Volumes  map[string]ClusterTemplatePodgroupVolume `json:"volumes"`
	Count    int                                      `json:"count"`
}

type ClusterTemplatePodgroupResource struct {
	CpuRequest string  `json:"cpu_request"`
	CpuMargin  float64 `json:"cpu_margin"`
	RamRequest string  `json:"ram_request"`
	RamMargin  float64 `json:"ram_margin"`
}

type ClusterTemplatePodgroupVolume struct {
	Count            int    `json:"count"`
	StorageClassName string `json:"storageClassName"`
	Storage          string `json:"storage"`
}

type commonClusterTemplateResult struct {
	gophercloud.Result
}

// GetResult represents result of database cluster get
type GetResult struct {
	commonClusterTemplateResult
}

// Extract is used to extract result into response struct
func (r commonClusterTemplateResult) Extract() (*ClusterTemplates, error) {
	var c *ClusterTemplates
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c, nil
}
