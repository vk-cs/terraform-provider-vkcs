package clusters

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/common"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// ClusterCreate represents dataplatform cluster creation parameters
type ClusterCreate struct {
	Name              string            `json:"name" required:"true"`
	ClusterTemplateID string            `json:"cluster_template_id" required:"true"`
	NetworkID         string            `json:"network_id" required:"true"`
	SubnetID          string            `json:"subnet_id" required:"true"`
	ProductName       string            `json:"product_name" required:"true"`
	ProductVersion    string            `json:"product_version" required:"true"`
	AvailabilityZone  string            `json:"availability_zone" required:"true"`
	Configs           *ClusterConfig    `json:"configs" required:"true"`
	PodGroups         []ClusterPodGroup `json:"pod_groups" required:"true"`
	Description       string            `json:"description,omitempty"`
}

type ClusterUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ClusterConfig struct {
	Settings    []ClusterConfigSetting    `json:"settings" required:"true"`
	Maintenance *ClusterConfigMaintenance `json:"maintenance" required:"true"`
	Warehouses  []ClusterConfigWarehouse  `json:"warehouses" required:"true"`
}

type ClusterConfigFeature struct {
	VolumeAutoresize *ClusterConfigFeatureVolumeAutoresize `json:"volume_autoresize,omitempty"`
}

type ClusterConfigFeatureVolumeAutoresize struct {
	Data *ClusterConfigFeatureVolumeAutoresizeObj `json:"data,omitempty"`
	Wal  *ClusterConfigFeatureVolumeAutoresizeObj `json:"wal,omitempty"`
}

type ClusterConfigFeatureVolumeAutoresizeObj struct {
	ScaleStepSize      int  `json:"scale_step_size,omitempty"`
	MaxScaleSize       int  `json:"max_scale_size,omitempty"`
	SizeScaleThreshold int  `json:"size_scale_threshold,omitempty"`
	Enabled            bool `json:"enabled,omitempty"`
}

type ClusterConfigSetting struct {
	Alias string `json:"alias" required:"true"`
	Value string `json:"value" required:"true"`
}

type ClusterConfigMaintenance struct {
	Start    string                             `json:"start" required:"true"`
	Backup   *ClusterConfigMaintenanceBackup    `json:"backup,omitempty"`
	CronTabs []ClusterConfigMaintenanceCronTabs `json:"cron_tabs,omitempty"`
}

type ClusterConfigMaintenanceBackup struct {
	Full         *ClusterConfigMaintenanceBackupObj `json:"full,omitempty"`
	Incremental  *ClusterConfigMaintenanceBackupObj `json:"incremental,omitempty"`
	Differential *ClusterConfigMaintenanceBackupObj `json:"differential,omitempty"`
}

type ClusterConfigMaintenanceBackupObj struct {
	Enabled   bool   `json:"enabled,omitempty"`
	Start     string `json:"start" required:"true"`
	KeepCount int    `json:"keep_count,omitempty"`
	KeepTime  int    `json:"keep_time,omitempty"`
}

type ClusterConfigMaintenanceCronTabs struct {
	Required bool                   `json:"required" required:"true"`
	Name     string                 `json:"name" required:"true"`
	Start    string                 `json:"start" required:"true"`
	Settings []ClusterConfigSetting `json:"settings,omitempty"`
}

type ClusterConfigWarehouse struct {
	ID          string                             `json:"id,omitempty"`
	Name        string                             `json:"name,omitempty"`
	Connections []ClusterConfigWarehouseConnection `json:"connections" required:"true"`
	Extensions  []ClusterConfigWarehouseExtension  `json:"extensions,omitempty"`
}

type ClusterConfigWarehouseConnection struct {
	Name      string                 `json:"name" required:"true"`
	Plug      string                 `json:"plug" required:"true"`
	Settings  []ClusterConfigSetting `json:"settings" required:"true"`
	ID        string                 `json:"id,omitempty"`
	CreatedAt string                 `json:"created_at,omitempty"`
}

type ClusterConfigWarehouseExtension struct {
	Type      string                 `json:"type" required:"true"`
	Version   string                 `json:"version,omitempty"`
	Settings  []ClusterConfigSetting `json:"settings,omitempty"`
	ID        string                 `json:"id,omitempty"`
	CreatedAt string                 `json:"created_at,omitempty"`
	Name      string                 `json:"name,omitempty"`
}

type ClusterPodGroup struct {
	ID                 string                           `json:"id,omitempty"`
	Name               string                           `json:"name,omitempty"`
	Count              int                              `json:"count" required:"true"`
	Resource           *ClusterPodGroupResource         `json:"resource" required:"true"`
	PodGroupTemplateID string                           `json:"pod_group_template_id" required:"true"`
	Volumes            map[string]ClusterPodGroupVolume `json:"volumes,omitempty"`
	FloatingIPPool     string                           `json:"floating_ip_pool,omitempty"`
	AvailabilityZone   string                           `json:"availability_zone,omitempty"`
	Alias              string                           `json:"alias,omitempty"`
	NodeProcesses      []string                         `json:"nodeProcesses,omitempty"`
}

type ClusterPodGroupResource struct {
	CPURequest string `json:"cpu_request" required:"true"`
	CPULimit   string `json:"cpu_limit,omitempty"`
	RAMRequest string `json:"ram_request" required:"true"`
	RAMLimit   string `json:"ram_limit,omitempty"`
}

type ClusterPodGroupVolume struct {
	StorageClassName string `json:"storageClassName" required:"true"`
	Storage          string `json:"storage" required:"true"`
	Count            int    `json:"count" required:"true"`
}

type ClusterInfo struct {
	Services []ClusterInfoService `json:"services,omitempty"`
}

type ClusterInfoService struct {
	Type             string `json:"type" required:"true"`
	Exposed          bool   `json:"exposed,omitempty"`
	Description      string `json:"description,omitempty"`
	ConnectionString string `json:"connection_string" required:"true"`
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreate) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterUpdate) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterConfig) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterConfigFeature) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterConfigFeatureVolumeAutoresize) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterConfigFeatureVolumeAutoresizeObj) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterConfigSetting) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterConfigMaintenance) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterConfigMaintenanceBackup) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterConfigMaintenanceBackupObj) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterConfigMaintenanceCronTabs) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterConfigWarehouse) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterConfigWarehouseConnection) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterConfigWarehouseExtension) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterPodGroup) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterPodGroupResource) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterPodGroupVolume) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterInfo) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterInfoService) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Create performs request to create database cluster
func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	common.SetHeaders(client)
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(clustersURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	common.SetHeaders(client)
	resp, err := client.Get(clusterURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	common.SetHeaders(client)
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(clusterURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	common.SetHeaders(client)
	resp, err := client.Delete(clusterURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}
