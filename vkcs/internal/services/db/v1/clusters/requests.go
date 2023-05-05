package clusters

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/instances"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// Cluster is used to send proper request for cluster creation
type Cluster struct {
	Cluster *CreateOpts `json:"cluster" required:"true"`
}

// CreateOpts represents database cluster creation parameters
type CreateOpts struct {
	Name                   string                     `json:"name" required:"true"`
	Datastore              *datastores.DatastoreShort `json:"datastore" required:"true"`
	FloatingIPEnabled      bool                       `json:"allow_remote_access,omitempty"`
	AutoExpand             int                        `json:"volume_autoresize_enabled,omitempty"`
	MaxDiskSize            int                        `json:"volume_autoresize_max_size,omitempty"`
	WalAutoExpand          int                        `json:"wal_autoresize_enabled,omitempty"`
	WalMaxDiskSize         int                        `json:"wal_autoresize_max_size,omitempty"`
	Instances              []InstanceCreateOpts       `json:"instances"`
	Capabilities           []instances.CapabilityOpts `json:"capabilities,omitempty"`
	RestorePoint           *instances.RestorePoint    `json:"restorePoint,omitempty"`
	BackupSchedule         *instances.BackupSchedule  `json:"backup_schedule,omitempty"`
	CloudMonitoringEnabled bool                       `json:"cloud_monitoring_enabled,omitempty"`
}

// InstanceCreateOpts represents database cluster instance creation parameters
type InstanceCreateOpts struct {
	Keypair          string                  `json:"key_name,omitempty"`
	AvailabilityZone string                  `json:"availability_zone,omitempty"`
	FlavorRef        string                  `json:"flavorRef,omitempty" mapstructure:"flavor_id"`
	Nics             []instances.NetworkOpts `json:"nics" required:"true"`
	Volume           *instances.Volume       `json:"volume" required:"true"`
	Walvolume        *instances.WalVolume    `json:"wal_volume,omitempty"`
	ShardID          string                  `json:"shard_id,omitempty"`
	SecurityGroups   []string                `json:"security_groups,omitempty"`
}

// AttachConfigurationGroupOpts represents parameters of configuration group to be attached to database cluster
type AttachConfigurationGroupOpts struct {
	ConfigurationAttach struct {
		ConfigurationID string `json:"configuration_id"`
	} `json:"configuration_attach"`
}

// DetachConfigurationGroupOpts represents parameters of configuration group to be detached from database cluster
type DetachConfigurationGroupOpts struct {
	ConfigurationDetach struct {
		ConfigurationID string `json:"configuration_id"`
	} `json:"configuration_detach"`
}

// ResizeVolumeOpts represents parameters of volume resize of database cluster
type ResizeVolumeOpts struct {
	Resize struct {
		Volume struct {
			Size int `json:"size"`
		} `json:"volume"`
		ShardID string `json:"shard_id,omitempty"`
	} `json:"resize"`
}

// ResizeWalVolumeOpts represents parameters of wal volume resize of database cluster
type ResizeWalVolumeOpts struct {
	Resize struct {
		Volume struct {
			Size int    `json:"size"`
			Kind string `json:"kind"`
		} `json:"volume"`
		ShardID string `json:"shard_id,omitempty"`
	} `json:"resize"`
}

// ResizeOpts represents database cluster resize parameters
type ResizeOpts struct {
	Resize struct {
		FlavorRef string `json:"flavorRef"`
		ShardID   string `json:"shard_id,omitempty"`
	} `json:"resize"`
}

// UpdateAutoExpandOpts represents autoresize parameters of volume of database cluster
type UpdateAutoExpandOpts struct {
	Cluster struct {
		VolumeAutoresizeEnabled int `json:"volume_autoresize_enabled"`
		VolumeAutoresizeMaxSize int `json:"volume_autoresize_max_size"`
	} `json:"cluster"`
}

// UpdateAutoExpandWalOpts represents autoresize parameters of wal volume of database cluster
type UpdateAutoExpandWalOpts struct {
	Cluster struct {
		WalVolume struct {
			VolumeAutoresizeEnabled int `json:"autoresize_enabled"`
			VolumeAutoresizeMaxSize int `json:"autoresize_max_size"`
		} `json:"wal_volume"`
	} `json:"cluster"`
}

// ApplyCapabilityOpts represents parameters of capabilities to be applied to database cluster
type ApplyCapabilityOpts struct {
	ApplyCapability struct {
		Capabilities []instances.CapabilityOpts `json:"capabilities"`
	} `json:"apply_capability"`
}

type BackupScheduleOpts struct {
	BackupSchedule *instances.BackupSchedule `json:"backup_schedule"`
}

// UpdateCloudMonitoringOpts represents parameters of request to update cloud monitoring options
type UpdateCloudMonitoringOpts struct {
	CloudMonitoring struct {
		Enable bool `json:"enable"`
	} `json:"cloud_monitoring"`
}

// GrowClusterOpts is used to send proper request to grow cluster
type GrowClusterOpts struct {
	Grow []GrowOpts `json:"grow"`
}

// GrowOpts represents parameters of growing cluster
type GrowOpts struct {
	Keypair          string               `json:"key_name"`
	AvailabilityZone string               `json:"availability_zone" required:"true"`
	FlavorRef        string               `json:"flavorRef" required:"true"`
	Volume           *instances.Volume    `json:"volume" required:"true"`
	Walvolume        *instances.WalVolume `json:"wal_volume,omitempty"`
	ShardID          string               `json:"shard_id,omitempty"`
}

// ShrinkClusterOpts is used to send proper request to shrink database cluster
type ShrinkClusterOpts struct {
	Shrink []ShrinkOpts `json:"shrink"`
}

// ClusterShrinkOpts represents parameters of shrinking database cluster
type ShrinkOpts struct {
	ID string `json:"id" required:"true"`
}

// Map converts opts to a map (for a request body)
func (opts Cluster) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Map converts opts to a map (for a request body)
func (opts *DetachConfigurationGroupOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *AttachConfigurationGroupOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ResizeVolumeOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ResizeWalVolumeOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ResizeOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *UpdateAutoExpandOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *UpdateAutoExpandWalOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ApplyCapabilityOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *UpdateCloudMonitoringOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *GrowClusterOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ShrinkClusterOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Create performs request to create database cluster
func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(clustersURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// Get performs request to get database cluster
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	var result *http.Response
	result, r.Err = client.Get(clusterURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func GetBackupSchedule(client *gophercloud.ServiceClient, id string) (r GetBackupScheduleResult) {
	var result *http.Response
	result, r.Err = client.Get(backupScheduleURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func GetCapabilities(client *gophercloud.ServiceClient, id string) (r GetCapabilitiesResult) {
	var result *http.Response
	result, r.Err = client.Get(capabilitiesURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// ClusterAction performs request to perform an action on the database cluster
func ClusterAction(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r ActionResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(clusterURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// UpdateAutoExpand performs request to update database cluster autoresize parameters
func UpdateAutoExpand(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r ActionResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Patch(clusterURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func UpdateBackupSchedule(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateBackupScheduleResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Put(backupScheduleURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(clusterURL(client, id), &gophercloud.RequestOpts{})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
