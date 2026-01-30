package nodegroups

import (
	"github.com/gophercloud/gophercloud"
)

// CreateOpts represents options for creating a node group in V2 API
type CreateOpts struct {
	ClusterID string        `json:"cluster_id" required:"true"`
	Spec      NodeGroupSpec `json:"node_group_spec" required:"true"`
}

// UpdateOpts represents options for updating a node group in V2 API
type UpdateOpts struct {
	VMEngine             *VMEngine         `json:"vm_engine,omitempty"`
	ScaleSpec            *ScaleSpec        `json:"scale_spec,omitempty"`
	Labels               map[string]string `json:"labels,omitempty"`
	Taints               []Taint           `json:"taints,omitempty"`
	ParallelUpgradeChunk *int              `json:"parallel_upgrade_chunk,omitempty"`
}

// NodeGroupSpec represents node group specification
type NodeGroupSpec struct {
	Name                 string            `json:"name,omitempty"`
	VMEngine             VMEngine          `json:"vm_engine,omitempty"`
	Zones                []string          `json:"zones,omitempty"`
	ScaleSpec            ScaleSpec         `json:"scale_spec,omitempty"`
	Labels               map[string]string `json:"labels,omitempty"`
	Taints               []Taint           `json:"taints,omitempty"`
	ParallelUpgradeChunk int               `json:"parallel_upgrade_chunk,omitempty"`
	DiskType             DiskType          `json:"disk_type,omitempty"`
}

// VMEngine represents virtual machine engine configuration
type VMEngine struct {
	NovaEngine NovaEngine `json:"nova_engine,omitempty"`
}

// NovaEngine represents Nova-based engine configuration
type NovaEngine struct {
	FlavorID string `json:"flavor_id" required:"true"`
}

// ScaleSpec represents scaling specification
type ScaleSpec struct {
	FixedScale *FixedScale `json:"fixed_scale,omitempty"`
	AutoScale  *AutoScale  `json:"auto_scale,omitempty"`
}

// FixedScale represents fixed scaling configuration
type FixedScale struct {
	Size int `json:"size" required:"true"`
}

// AutoScale represents auto scaling configuration
type AutoScale struct {
	MinSize int `json:"min_size" required:"true"`
	MaxSize int `json:"max_size" required:"true"`
	Size    int `json:"size" required:"true"`
}

// Taint represents node taint configuration
type Taint struct {
	Key    string `json:"key" required:"true"`
	Value  string `json:"value" required:"true"`
	Effect string `json:"effect" required:"true"`
}

// DiskType represents disk type configuration
type DiskType struct {
	CinderVolumeType CinderVolumeType `json:"cinder_volume_type,omitempty"`
}

// CinderVolumeType represents Cinder volume type configuration
type CinderVolumeType struct {
	Type string `json:"type" required:"true"`
	Size int    `json:"size" required:"true"`
}

// ToNodeGroupCreateMap builds a request body from CreateOpts
func (opts CreateOpts) ToNodeGroupCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// ToNodeGroupUpdateMap builds a request body from UpdateOpts
func (opts UpdateOpts) ToNodeGroupUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create creates a new node group in the specified cluster.
func Create(client *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToNodeGroupCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Post(rootURL(client), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

// Get retrieves a specific node group based on its ID.
func Get(client *gophercloud.ServiceClient, nodeGroupID string) GetResult {
	var res GetResult

	_, res.Err = client.Get(resourceURL(client, nodeGroupID), &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

// Update updates a specific node group with new configuration.
func Scale(client *gophercloud.ServiceClient, nodeGroupID string, opts UpdateOpts) error {
	reqBody, err := opts.ToNodeGroupUpdateMap()
	if err != nil {
		return err
	}

	_, err = client.Patch(resourceURL(client, nodeGroupID), reqBody, nil, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return err
}

// Delete deletes a specific node group.
func Delete(client *gophercloud.ServiceClient, nodeGroupID string) error {
	_, err := client.Delete(resourceURL(client, nodeGroupID), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return err
}
