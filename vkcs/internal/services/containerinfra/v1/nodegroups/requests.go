package nodegroups

import (
	"net/http"
	"time"

	"github.com/gophercloud/gophercloud"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type PatchOptsBuilder interface {
	PatchMap() ([]map[string]interface{}, error)
}

type Node struct {
	Name        string     `json:"name"`
	UUID        string     `json:"uuid"`
	NodeGroupID string     `json:"node_group_id"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type PatchOpts []PatchParams

type PatchParams struct {
	Path  string      `json:"path,omitempty"`
	Value interface{} `json:"value,omitempty"`
	Op    string      `json:"op,omitempty"`
}

type BatchAddParams struct {
	Action  string      `json:"action,omitempty"`
	Payload []NodeGroup `json:"payload,omitempty"`
}

type NodeGroup struct {
	Name               string    `json:"name,omitempty"`
	NodeCount          int       `json:"node_count,omitempty"`
	MaxNodes           int       `json:"max_nodes,omitempty"`
	MinNodes           int       `json:"min_nodes,omitempty"`
	VolumeSize         int       `json:"volume_size,omitempty"`
	VolumeType         string    `json:"volume_type,omitempty"`
	FlavorID           string    `json:"flavor_id,omitempty"`
	ImageID            string    `json:"image_id,omitempty"`
	Autoscaling        bool      `json:"autoscaling_enabled,omitempty"`
	ClusterID          string    `json:"cluster_id,omitempty"`
	UUID               string    `json:"uuid,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
	Nodes              []*Node   `json:"nodes,omitempty"`
	State              string    `json:"state,omitempty"`
	AvailabilityZones  []string  `json:"availability_zones"`
	MaxNodeUnavailable int       `json:"max_node_unavailable,omitempty"`
}

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

type Taint struct {
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
	Effect string `json:"effect,omitempty"`
}

// CreateOpts contains options to create node group.
type CreateOpts struct {
	ClusterID          string   `json:"cluster_id" required:"true"`
	Name               string   `json:"name"`
	Labels             []Label  `json:"labels,omitempty"`
	Taints             []Taint  `json:"taints,omitempty"`
	NodeCount          int      `json:"node_count,omitempty"`
	MaxNodes           int      `json:"max_nodes,omitempty"`
	MinNodes           int      `json:"min_nodes,omitempty"`
	VolumeSize         int      `json:"volume_size,omitempty"`
	VolumeType         string   `json:"volume_type,omitempty"`
	FlavorID           string   `json:"flavor_id,omitempty"`
	Autoscaling        bool     `json:"autoscaling_enabled,omitempty"`
	AvailabilityZones  []string `json:"availability_zones,omitempty"`
	MaxNodeUnavailable int      `json:"max_node_unavailable,omitempty"`
}

// ScaleOpts contains options to scale node group
type ScaleOpts struct {
	Delta    int    `json:"delta" required:"true"`
	Rollback string `json:"rollback,omitempty"`
}

// Map builds request params.
func (opts *CreateOpts) Map() (map[string]interface{}, error) {
	cluster, err := gophercloud.BuildRequestBody(*opts, "")
	return cluster, err
}

// Map builds request params.
func (opts *NodeGroup) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map builds request params.
func (opts *ScaleOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map builds request params.
func (opts *BatchAddParams) Map() (map[string]interface{}, error) {
	batch, err := gophercloud.BuildRequestBody(*opts, "")
	return batch, err
}

// Map builds request params.
func (opts *PatchParams) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// PatchMap collects all the params.
func (opts *PatchOpts) PatchMap() ([]map[string]interface{}, error) {
	var lb []map[string]interface{}
	for _, opt := range *opts {
		b, err := opt.Map()
		if err != nil {
			return nil, err
		}
		lb = append(lb, b)
	}
	return lb, nil
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(nodeGroupsURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	var result *http.Response
	result, r.Err = client.Get(nodeGroupURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func Patch(client *gophercloud.ServiceClient, id string, opts PatchOptsBuilder) (r PatchResult) {
	b, err := opts.PatchMap()
	if err != nil {
		r.Err = err
	}
	var result *http.Response
	result, r.Err = client.Patch(nodeGroupURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func Scale(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r ScaleResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Patch(scaleURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(nodeGroupURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
