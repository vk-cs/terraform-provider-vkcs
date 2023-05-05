package instances

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
)

// InstanceResp represents result of database instance get
type InstanceResp struct {
	СomputeInstanceID string                     `json:"compute_instance_id"`
	Configuration     *Configuration             `json:"configuration"`
	ConfigurationID   string                     `json:"configuration_id"`
	ID                string                     `json:"id"`
	Created           db.DateTimeWithoutTZFormat `json:"created"`
	Updated           db.DateTimeWithoutTZFormat `json:"updated"`
	DataStore         *datastores.DatastoreShort `json:"datastore"`
	Flavor            *Links                     `json:"flavor"`
	GaVersion         string                     `json:"ga_version"`
	HealthStatus      string                     `json:"health_status"`
	IP                *[]string                  `json:"ip"`
	Links             *[]Link                    `json:"links"`
	Name              string                     `json:"name"`
	Region            string                     `json:"region"`
	Status            string                     `json:"status"`
	Volume            *Volume                    `json:"volume"`
	ReplicaOf         *Links                     `json:"replica_of"`
	AutoExpand        int                        `json:"volume_autoresize_enabled"`
	MaxDiskSize       int                        `json:"volume_autoresize_max_size"`
	WalVolume         *WalVolume                 `json:"wal_volume"`
}

type InstanceShortResp struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// InstanceRespOpts is used to get instance response
type InstanceRespOpts struct {
	Instance *InstanceResp `json:"instance"`
}

type InstanceShortRespOpts struct {
	Instance *InstanceShortResp `json:"instance"`
}

// volume represents database instance volume
type Volume struct {
	Size       *int     `json:"size" required:"true"`
	Used       *float32 `json:"used,omitempty"`
	VolumeID   string   `json:"volume_id,,omitempty"`
	VolumeType string   `json:"type,,omitempty" required:"true"`
}

// walVolume represents database instance wal volume
type WalVolume struct {
	Size        *int     `json:"size" required:"true"`
	Used        *float32 `json:"used,omitempty"`
	VolumeID    string   `json:"volume_id,,omitempty"`
	VolumeType  string   `json:"type,,omitempty" required:"true"`
	AutoExpand  int      `json:"autoresize_enabled,omitempty"`
	MaxDiskSize int      `json:"autoresize_max_size,omitempty"`
}

// links represents database instance links
type Links struct {
	ID    string  `json:"id"`
	Links *[]Link `json:"links"`
}

// link represents database instance link
type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

// configuration represents database instance configuration
type Configuration struct {
	ID    string  `json:"id"`
	Links *[]Link `json:"links"`
	Name  string  `json:"name"`
}

type RestorePoint struct {
	BackupRef string `json:"backupRef" required:"true" mapstructure:"backup_id"`
	Target    string `json:"target,omitempty"`
}

type BackupSchedule struct {
	Name          string `json:"name"`
	StartHours    int    `json:"start_hours"`
	StartMinutes  int    `json:"start_minutes"`
	IntervalHours int    `json:"interval_hours"`
	KeepCount     int    `json:"keep_count"`
}

type commonInstanceResult struct {
	gophercloud.Result
}

type commonInstanceShortResult struct {
	gophercloud.Result
}

type commonCapabilitiesResult struct {
	gophercloud.Result
}

type commonBackupScheduleResult struct {
	gophercloud.Result
}

// GetResult represents result of database instance get
type GetResult struct {
	commonInstanceResult
}

type CreateResult struct {
	commonInstanceShortResult
}

type GetCapabilitiesResult struct {
	commonCapabilitiesResult
}

type GetBackupScheduleResult struct {
	commonBackupScheduleResult
}

// RootUserResp represents parameters of root user response
type RootUserResp struct {
	Password string `json:"password"`
	Name     string `json:"name"`
}

// RootUserRespOpts is used to get root user response
type RootUserRespOpts struct {
	User *RootUserResp `json:"user"`
}

type commonRootUserResult struct {
	gophercloud.Result
}

type commonRootUserErrResult struct {
	gophercloud.ErrResult
}

// CreateRootUserResult represents result of root user create
type CreateRootUserResult struct {
	commonRootUserResult
}

// DeleteRootUserResult represents result of root user delete
type DeleteRootUserResult struct {
	commonRootUserErrResult
}

// ConfigurationResult represents result of configuration attach and detach
type ConfigurationResult struct {
	gophercloud.ErrResult
}

// ActionResult represents result of database instance action
type ActionResult struct {
	gophercloud.ErrResult
}

// IsRootUserEnabledResult represents result of getting root user status
type IsRootUserEnabledResult struct {
	commonRootUserResult
}

type UpdateBackupScheduleResult struct {
	gophercloud.ErrResult
}

// DeleteResult represents result of database instance delete
type DeleteResult struct {
	gophercloud.ErrResult
}

// Extract is used to extract result into response struct
func (r commonInstanceResult) Extract() (*InstanceResp, error) {
	var i *InstanceRespOpts
	if err := r.ExtractInto(&i); err != nil {
		return nil, err
	}
	return i.Instance, nil
}

func (r commonInstanceShortResult) Extract() (*InstanceShortResp, error) {
	var i *InstanceShortRespOpts
	if err := r.ExtractInto(&i); err != nil {
		return nil, err
	}
	return i.Instance, nil
}

// Extract is used to extract result into response struct
func (r commonRootUserResult) Extract() (*RootUserResp, error) {
	var u *RootUserRespOpts
	if err := r.ExtractInto(&u); err != nil {
		return nil, err
	}
	return u.User, nil
}

// Extract is used to extract result into response struct
func (r IsRootUserEnabledResult) Extract() (bool, error) {
	return r.Body.(map[string]interface{})["rootEnabled"].(bool), r.Err
}

func (r commonCapabilitiesResult) Extract() ([]DatabaseCapability, error) {
	var c *GetCapabilityOpts
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c.Capabilities, nil
}

func (r commonBackupScheduleResult) Extract() (*BackupSchedule, error) {
	var b *BackupScheduleOpts
	if r.Body == nil {
		return nil, nil
	}
	if err := r.ExtractInto(&b); err != nil {
		return nil, err
	}
	return b.BackupSchedule, nil
}
