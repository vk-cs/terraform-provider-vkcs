package clusters

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/instances"
)

// ClusterResp represents database cluster response
type ClusterResp struct {
	ConfigurationID string                     `json:"configuration_id"`
	Created         db.DateTimeWithoutTZFormat `json:"created"`
	DataStore       *datastores.DatastoreShort `json:"datastore"`
	HealthStatus    string                     `json:"health_status"`
	ID              string                     `json:"id"`
	Instances       []ClusterInstanceResp      `json:"instances"`
	Links           *[]instances.Link          `json:"links"`
	LoadbalancerID  string                     `json:"loadbalancer_id"`
	Name            string                     `json:"name"`
	Task            Task                       `json:"task"`
	Updated         db.DateTimeWithoutTZFormat `json:"updated"`
	AutoExpand      int                        `json:"volume_autoresize_enabled"`
	MaxDiskSize     int                        `json:"volume_autoresize_max_size"`
	VRRPPortID      string                     `json:"vrrp_port_id"`
	WalAutoExpand   int                        `json:"wal_autoresize_enabled"`
	WalMaxDiskSize  int                        `json:"wal_autoresize_max_size"`
}

// ClusterInstanceResp represents database cluster instance response
type ClusterInstanceResp struct {
	СomputeInstanceID string               `json:"compute_instance_id"`
	Flavor            *instances.Links     `json:"flavor"`
	GaVersion         string               `json:"ga_version"`
	ID                string               `json:"id"`
	IP                *[]string            `json:"ip"`
	Links             *[]instances.Link    `json:"links"`
	Name              string               `json:"name"`
	Role              string               `json:"role"`
	Status            string               `json:"status"`
	Type              string               `json:"type"`
	Volume            *instances.Volume    `json:"volume"`
	WalVolume         *instances.WalVolume `json:"wal_volume"`
	ShardID           string               `json:"shard_id"`
}

type ClusterShortResp struct {
	ID string `json:"id"`
}

// ClusterRespOpts is used to properly extract database cluster response
type ClusterRespOpts struct {
	Cluster *ClusterResp `json:"cluster"`
}

type ClusterShortRespOpts struct {
	Cluster *ClusterShortResp `json:"cluster"`
}

// Task represents database cluster task
type Task struct {
	Description string `json:"description"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
}

type GetCapabilityOpts struct {
	Capabilities []instances.DatabaseCapability `json:"capabilities"`
}

type commonClusterResult struct {
	gophercloud.Result
}

type commonShortClusterResult struct {
	gophercloud.Result
}

type commonBackupScheduleResult struct {
	gophercloud.Result
}

type commonCapabilitiesResult struct {
	gophercloud.Result
}

// CreateResult represents result of database cluster create
type CreateResult struct {
	commonShortClusterResult
}

// GetResult represents result of database cluster get
type GetResult struct {
	commonClusterResult
}

// Page represents a page of database clusters
type Page struct {
	pagination.SinglePageBase
}

// IsEmpty indicates whether a database cluster collection is empty.
func (r Page) IsEmpty() (bool, error) {
	is, err := ExtractClusters(r)
	return len(is) == 0, err
}

// ExtractClusters retrieves a slice of database clusterResp structs from a paginated
// collection.
func ExtractClusters(r pagination.Page) ([]ClusterResp, error) {
	var s struct {
		Clusters []ClusterResp `json:"clusters"`
	}
	err := (r.(Page)).ExtractInto(&s)
	return s.Clusters, err
}

type DeleteResult struct {
	gophercloud.ErrResult
}

// ActionResult represents result of database cluster action
type ActionResult struct {
	gophercloud.ErrResult
}

type GetBackupScheduleResult struct {
	commonBackupScheduleResult
}

type UpdateBackupScheduleResult struct {
	gophercloud.ErrResult
}

type GetCapabilitiesResult struct {
	commonCapabilitiesResult
}

// Extract is used to extract result into response struct
func (r commonClusterResult) Extract() (*ClusterResp, error) {
	var c *ClusterRespOpts
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c.Cluster, nil
}

// Extract is used to extract result into short response struct
func (r commonShortClusterResult) Extract() (*ClusterShortResp, error) {
	var c *ClusterShortRespOpts
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c.Cluster, nil
}

// Extract is used to extract result into response struct
func (r commonBackupScheduleResult) Extract() (*instances.BackupSchedule, error) {
	var b *BackupScheduleOpts
	if r.Body == nil {
		return nil, nil
	}
	if err := r.ExtractInto(&b); err != nil {
		return nil, err
	}
	return b.BackupSchedule, nil
}

func (r commonCapabilitiesResult) Extract() ([]instances.DatabaseCapability, error) {
	var c *GetCapabilityOpts
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c.Capabilities, nil
}
