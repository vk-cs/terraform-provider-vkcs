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
	Name              string                  `json:"name" required:"true"`
	ClusterTemplateID string                  `json:"cluster_template_id" required:"true"`
	NetworkID         string                  `json:"network_id" required:"true"`
	SubnetID          string                  `json:"subnet_id" required:"true"`
	ProductName       string                  `json:"product_name" required:"true"`
	ProductVersion    string                  `json:"product_version" required:"true"`
	AvailabilityZone  string                  `json:"availability_zone" required:"true"`
	Configs           *ClusterCreateConfig    `json:"configs" required:"true"`
	PodGroups         []ClusterCreatePodGroup `json:"pod_groups" required:"true"`
	Description       string                  `json:"description,omitempty"`
	StackID           string                  `json:"stack_id,omitempty"`
}

type ClusterUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ClusterCreateConfig struct {
	Settings    []ClusterCreateConfigSetting    `json:"settings,omitempty"`
	Maintenance *ClusterCreateConfigMaintenance `json:"maintenance" required:"true"`
	Warehouses  []ClusterCreateConfigWarehouse  `json:"warehouses" required:"true"`
}

type ClusterCreateConfigSetting struct {
	Alias string `json:"alias" required:"true"`
	Value string `json:"value" required:"true"`
}

type ClusterCreateConfigMaintenance struct {
	Start    string                                   `json:"start" required:"true"`
	Backup   *ClusterCreateConfigMaintenanceBackup    `json:"backup,omitempty"`
	CronTabs []ClusterCreateConfigMaintenanceCronTabs `json:"cron_tabs,omitempty"`
}

type ClusterCreateConfigMaintenanceBackup struct {
	Full         *ClusterCreateConfigMaintenanceBackupObj `json:"full,omitempty"`
	Incremental  *ClusterCreateConfigMaintenanceBackupObj `json:"incremental,omitempty"`
	Differential *ClusterCreateConfigMaintenanceBackupObj `json:"differential,omitempty"`
}

type ClusterCreateConfigMaintenanceBackupObj struct {
	Start     string `json:"start" required:"true"`
	KeepCount int    `json:"keep_count,omitempty"`
	KeepTime  int    `json:"keep_time,omitempty"`
}

type ClusterCreateConfigMaintenanceCronTabs struct {
	Name     string                       `json:"name" required:"true"`
	Start    string                       `json:"start" required:"true"`
	Settings []ClusterCreateConfigSetting `json:"settings,omitempty"`
}

type ClusterCreateConfigWarehouse struct {
	Name        string                                   `json:"name,omitempty"`
	Connections []ClusterCreateConfigWarehouseConnection `json:"connections" required:"true"`
	Extensions  []ClusterCreateConfigWarehouseExtension  `json:"extensions,omitempty"`
}

type ClusterCreateConfigWarehouseConnection struct {
	Name     string                       `json:"name" required:"true"`
	Plug     string                       `json:"plug" required:"true"`
	Settings []ClusterCreateConfigSetting `json:"settings" required:"true"`
}

type ClusterCreateConfigWarehouseExtension struct {
	Type     string                       `json:"type" required:"true"`
	Version  string                       `json:"version,omitempty"`
	Settings []ClusterCreateConfigSetting `json:"settings,omitempty"`
}

type ClusterCreatePodGroup struct {
	Count              int                                    `json:"count" required:"true"`
	Resource           *ClusterCreatePodGroupResource         `json:"resource" required:"true"`
	PodGroupTemplateID string                                 `json:"pod_group_template_id" required:"true"`
	Volumes            map[string]ClusterCreatePodGroupVolume `json:"volumes,omitempty"`
	FloatingIPPool     string                                 `json:"floating_ip_pool,omitempty"`
}

type ClusterCreatePodGroupResource struct {
	CPURequest string `json:"cpu_request" required:"true"`
	RAMRequest string `json:"ram_request" required:"true"`
}

type ClusterCreatePodGroupVolume struct {
	StorageClassName string `json:"storageClassName" required:"true"`
	Storage          string `json:"storage" required:"true"`
	Count            int    `json:"count" required:"true"`
}

type ClusterUpdateSettings struct {
	Settings []ClusterUpdateSetting `json:"settings" required:"true"`
}

type ClusterUpdateSetting struct {
	Alias string `json:"alias" required:"true"`
	Value string `json:"value" required:"true"`
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
func (opts *ClusterCreateConfig) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigSetting) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigMaintenance) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigMaintenanceBackup) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigMaintenanceBackupObj) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigMaintenanceCronTabs) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigWarehouse) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigWarehouseConnection) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigWarehouseExtension) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreatePodGroup) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreatePodGroupResource) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreatePodGroupVolume) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterUpdateSettings) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterUpdateSetting) Map() (map[string]interface{}, error) {
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

func UpdateSettings(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	common.SetHeaders(client)
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(clusterSettingsURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
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
