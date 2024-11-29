package instances

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// Instance is used to send request to create database instance
type Instance struct {
	Instance *CreateOpts `json:"instance" required:"true"`
}

// CreateOpts represents parameters of creation of database instance
type CreateOpts struct {
	FlavorRef              string                     `json:"flavorRef,omitempty"`
	Volume                 *Volume                    `json:"volume" required:"true"`
	Name                   string                     `json:"name" required:"true"`
	Configuration          string                     `json:"configuration,omitempty"`
	Datastore              *datastores.DatastoreShort `json:"datastore" required:"true"`
	Nics                   []NetworkOpts              `json:"nics" required:"true"`
	ReplicaOf              string                     `json:"replica_of,omitempty"`
	AvailabilityZone       string                     `json:"availability_zone,omitempty"`
	FloatingIPEnabled      bool                       `json:"allow_remote_access,omitempty"`
	Keypair                string                     `json:"key_name,omitempty"`
	AutoExpand             *int                       `json:"volume_autoresize_enabled,omitempty"`
	MaxDiskSize            int                        `json:"volume_autoresize_max_size,omitempty"`
	Walvolume              *WalVolume                 `json:"wal_volume,omitempty"`
	Capabilities           []CapabilityOpts           `json:"capabilities,omitempty"`
	RestorePoint           *RestorePoint              `json:"restorePoint,omitempty"`
	BackupSchedule         *BackupSchedule            `json:"backup_schedule,omitempty"`
	CloudMonitoringEnabled bool                       `json:"cloud_monitoring_enabled,omitempty"`
	SecurityGroups         []string                   `json:"security_groups,omitempty"`
}

// NetworkOpts represents network parameters of database instance
type NetworkOpts struct {
	UUID      string `json:"net-id,omitempty"`
	Port      string `json:"port-id,omitempty"`
	V4FixedIP string `json:"v4-fixed-ip,omitempty" mapstructure:"fixed_ip_v4"`
	SubnetID  string `json:"subnet-id,omitempty" mapstructure:"subnet_id"`
}

type BackupScheduleOpts struct {
	BackupSchedule *BackupSchedule `json:"backup_schedule"`
}

// AutoExpandOpts represents autoresize parameters of volume of database instance
type AutoExpandOpts struct {
	AutoExpand  bool
	MaxDiskSize int `mapstructure:"max_disk_size"`
}

// DetachReplicaOpts represents parameters of request to detach replica of database instance
type DetachReplicaOpts struct {
	Instance struct {
		ReplicaOf string `json:"replica_of,omitempty"`
	} `json:"instance"`
}

// AttachConfigurationGroupOpts represents parameters of configuration group to be attached to database instance
type AttachConfigurationGroupOpts struct {
	RestartConfirmed *bool `json:"restart_confirmed"`
	Instance         struct {
		Configuration string `json:"configuration"`
	} `json:"instance"`
}

// DetachConfigurationGroupOpts represents parameters of configuration group to be detached from database instance
type DetachConfigurationGroupOpts struct {
	RestartConfirmed *bool `json:"restart_confirmed"`
	Instance         struct {
		Configuration string `json:"configuration"`
	} `json:"instance"`
}

// UpdateAutoExpandOpts represents parameters of request to update autoresize properties of volume of database instance
type UpdateAutoExpandOpts struct {
	Instance struct {
		VolumeAutoresizeEnabled int `json:"volume_autoresize_enabled"`
		VolumeAutoresizeMaxSize int `json:"volume_autoresize_max_size"`
	} `json:"instance"`
}

// ResizeVolumeOpts represents database instance volume resize parameters
type ResizeVolumeOpts struct {
	Resize struct {
		Volume struct {
			Size int `json:"size"`
		} `json:"volume"`
	} `json:"resize"`
}

// UpdateAutoExpandWalOpts represents parameters of request to update autoresize properties of wal volume of database instance
type UpdateAutoExpandWalOpts struct {
	Instance struct {
		WalVolume struct {
			VolumeAutoresizeEnabled int `json:"autoresize_enabled"`
			VolumeAutoresizeMaxSize int `json:"autoresize_max_size"`
		} `json:"wal_volume"`
	} `json:"instance"`
}

// WalVolumeOpts represents parameters for creation of database instance wal volume
type WalVolumeOpts struct {
	Size       int
	VolumeType string `mapstructure:"volume_type"`
}

// ResizeWalVolumeOpts represents database instance wal volume resize parameters
type ResizeWalVolumeOpts struct {
	Resize struct {
		Volume struct {
			Size int    `json:"size"`
			Kind string `json:"kind"`
		} `json:"volume"`
	} `json:"resize"`
}

// ResizeOpts represents database instance resize parameters
type ResizeOpts struct {
	Resize struct {
		FlavorRef string `json:"flavorRef"`
	} `json:"resize"`
}

// RootUserEnableOpts represents parameters of request to enable root user for database instance
type RootUserEnableOpts struct {
	Password string `json:"password,omitempty"`
}

// UpdateCloudMonitoringOpts represents parameters of request to update cloud monitoring options
type UpdateCloudMonitoringOpts struct {
	CloudMonitoring struct {
		Enable bool `json:"enable"`
	} `json:"cloud_monitoring"`
}

// CapabilityOpts represents parameters of database instance capabilities
type CapabilityOpts struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params,omitempty" mapstructure:"settings"`
}

// DatabaseCapability represents capability info from dbaas
type DatabaseCapability struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params,omitempty"`
	Status string            `json:"status"`
}

// ApplyCapabilityOpts is used to send request to apply capability to database instance
type ApplyCapabilityOpts struct {
	ApplyCapability struct {
		Capabilities []CapabilityOpts `json:"capabilities"`
	} `json:"apply_capability"`
}

type GetCapabilityOpts struct {
	Capabilities []DatabaseCapability `json:"capabilities"`
}

// Map converts opts to a map (for a request body)
func (opts *Instance) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

// Map converts opts to a map (for a request body)
func (opts *DetachReplicaOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
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
func (opts *RootUserEnableOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *UpdateCloudMonitoringOpts) Map() (map[string]interface{}, error) {
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

func (opts *BackupSchedule) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Create performs request to create database instance
func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(instancesURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// Get performs request to get database instance
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(instanceURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func GetCapabilities(client *gophercloud.ServiceClient, id string) (r GetCapabilitiesResult) {
	resp, err := client.Get(capabilitiesURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func GetBackupSchedule(client *gophercloud.ServiceClient, id string) (r GetBackupScheduleResult) {
	resp, err := client.Get(backupScheduleURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// DetachReplica performs request to detach replica of database instance
func DetachReplica(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r ActionResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(instanceURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// UpdateAutoExpand performs request to update database instance autoresize parameters
func UpdateAutoExpand(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r ActionResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(instanceURL(client, id), &b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// DetachConfigurationGroup performs request to detach configuration group from database instance
func DetachConfigurationGroup(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r ConfigurationResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(instanceURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// AttachConfigurationGroup performs request to attach configuration group to database instance
func AttachConfigurationGroup(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r ConfigurationResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(instanceURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// Action performs request to perform an action on the database instance
func Action(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r ActionResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// RootUserEnable performs request to enable root user on database instance
func RootUserEnable(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r CreateRootUserResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(rootUserURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// RootUserGet performs request to get root user of database instance
func RootUserGet(client *gophercloud.ServiceClient, id string) (r IsRootUserEnabledResult) {
	resp, err := client.Get(rootUserURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// RootUserDisable performs request to disable root user on database instance
func RootUserDisable(client *gophercloud.ServiceClient, id string) (r DeleteRootUserResult) {
	resp, err := client.Delete(rootUserURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func UpdateBackupSchedule(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateBackupScheduleResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(backupScheduleURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// Delete performs request to delete database instance
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(instanceURL(client, id), &gophercloud.RequestOpts{})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// List will list all database instances
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, instancesURL(client),
		func(r pagination.PageResult) pagination.Page {
			return Page{pagination.SinglePageBase(r)}
		})
}
