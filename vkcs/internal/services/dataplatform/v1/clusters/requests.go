package clusters

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/common"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type OptsQueryBuilder interface {
	ToQuery() (string, error)
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
	FloatingIPPool    string                  `json:"floating_ip_pool,omitempty"`
}

type ClusterUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ClusterCreateConfig struct {
	Settings    []ClusterCreateConfigSetting    `json:"settings,omitempty"`
	Maintenance *ClusterCreateConfigMaintenance `json:"maintenance" required:"true"`
	Users       []ClusterCreateConfigUser       `json:"users,omitempty"`
	Warehouses  []ClusterCreateConfigWarehouse  `json:"warehouses" required:"true"`
}

type ClusterCreateConfigSetting struct {
	Alias string `json:"alias" required:"true"`
	Value string `json:"value" required:"true"`
}

type ClusterCreateConfigMaintenance struct {
	Start    string                                   `json:"start,omitempty"`
	Backup   *ClusterCreateConfigMaintenanceBackup    `json:"backup,omitempty"`
	CronTabs []ClusterCreateConfigMaintenanceCronTabs `json:"crontabs,omitempty"`
}

type ClusterCreateConfigMaintenanceBackup struct {
	Full         *ClusterCreateConfigMaintenanceBackupObj `json:"full,omitempty"`
	Incremental  *ClusterCreateConfigMaintenanceBackupObj `json:"incremental,omitempty"`
	Differential *ClusterCreateConfigMaintenanceBackupObj `json:"differential,omitempty"`
}

type ClusterCreateConfigMaintenanceBackupObj struct {
	Start     string `json:"start" required:"true"`
	KeepCount int    `json:"keepCount,omitempty"`
	KeepTime  int    `json:"keepTime,omitempty"`
}

type ClusterCreateConfigMaintenanceCronTabs struct {
	Name     string                       `json:"name" required:"true"`
	Start    string                       `json:"start" required:"true"`
	Settings []ClusterCreateConfigSetting `json:"settings,omitempty"`
}

type ClusterCreateConfigUser struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
	Role     string `json:"role,omitempty"`
}

type ClusterCreateConfigWarehouse struct {
	Name        string                                   `json:"name,omitempty"`
	Connections []ClusterCreateConfigWarehouseConnection `json:"connections,omitempty"`
	Extensions  []ClusterCreateConfigWarehouseExtension  `json:"extensions,omitempty"`
	Users       []string                                 `json:"users,omitempty"`
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
	Count              *int                                   `json:"count" required:"true"`
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

type ClusterUpdateUsers struct {
	Users []ClusterUpdateUser `json:"users" required:"true"`
}

type ClusterUpdateUser struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
	Role     string `json:"role,omitempty"`
}

type ClusterAddWarehouseConnections struct {
	WarehouseID string                 `json:"warehouse_id" required:"true"`
	Connections []ClusterAddConnection `json:"connections" required:"true"`
}

type ClusterAddConnection struct {
	Plug     string                        `json:"plug"`
	Name     string                        `json:"name"`
	Settings []ClusterAddConnectionSetting `json:"settings"`
}

type ClusterAddConnectionSetting struct {
	Alias string `json:"alias" required:"true"`
	Value string `json:"value" required:"true"`
}

type ClusterUpdatePodGroups struct {
	PodGroups []ClusterUpdatePodGroup `json:"pod_groups" required:"true"`
}

type ClusterUpdatePodGroup struct {
	ID       string                                 `json:"id" required:"true"`
	Resource ClusterUpdatePodGroupResource          `json:"resource,omitempty"`
	Count    *int                                   `json:"count" required:"true"`
	Volumes  map[string]ClusterUpdatePodGroupVolume `json:"volumes,omitempty"`
}

type ClusterUpdatePodGroupResource struct {
	CPURequest string `json:"cpu_request,omitempty"`
	RAMRequest string `json:"ram_request,omitempty"`
}

type ClusterUpdatePodGroupVolume struct {
	StorageClassName string `json:"storageClassName" required:"true"`
	Storage          string `json:"storage" required:"true"`
	Count            int    `json:"count" required:"true"`
}

type ClusterDeleteUsers struct {
	ClusterUsersIDs []string `q:"cluster_users_ids"`
}

type ClusterUpdateConfigsMaintenance struct {
	Start    *string                                  `json:"start,omitempty"`
	Backup   *ClusterCreateConfigMaintenanceBackup    `json:"backup,omitempty"`
	Crontabs *ClusterUpdateConfigsMaintenanceCrontabs `json:"crontabs,omitempty"`
}

type ClusterUpdateConfigsMaintenanceCrontabs struct {
	Create []ClusterUpdateConfigsMaintenanceCrontabsCreate `json:"create,omitempty"`
	Update []ClusterUpdateConfigsMaintenanceCrontabsUpdate `json:"update,omitempty"`
	Delete []ClusterUpdateConfigsMaintenanceCrontabsDelete `json:"delete,omitempty"`
}

type ClusterUpdateConfigsMaintenanceCrontabsCreate struct {
	Name     string                       `json:"name"`
	Start    string                       `json:"start"`
	Settings []ClusterCreateConfigSetting `json:"settings,omitempty"`
}

type ClusterUpdateConfigsMaintenanceCrontabsUpdate struct {
	ID       string                       `json:"id"`
	Start    string                       `json:"start"`
	Settings []ClusterCreateConfigSetting `json:"settings,omitempty"`
}

type ClusterUpdateConfigsMaintenanceCrontabsDelete struct {
	ID string `json:"id"`
}

type ClusterDeleteWarehouseConnections struct {
	ClusterConnections []string `q:"cluster_connections"`
	WarehouseID        string   `q:"warehouse_id"`
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreate) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterUpdate) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfig) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigSetting) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigMaintenance) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigMaintenanceBackup) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigMaintenanceBackupObj) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigMaintenanceCronTabs) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigUser) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigWarehouse) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigWarehouseConnection) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreateConfigWarehouseExtension) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreatePodGroup) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreatePodGroupResource) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterCreatePodGroupVolume) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterUpdateSettings) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterUpdateSetting) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterUpdateUsers) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterUpdateUser) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterAddWarehouseConnections) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterAddConnection) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterAddConnectionSetting) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterUpdatePodGroups) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterUpdatePodGroup) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

func (opts *ClusterDeleteUsers) ToQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(*opts)
	return q.String(), err
}

func (opts *ClusterDeleteWarehouseConnections) ToQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(*opts)
	return q.String(), err
}

// Map converts opts to a map (for a request body)
func (opts *ClusterUpdateConfigsMaintenance) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
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

func UpdateMaintenance(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	common.SetHeaders(client)
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(clusterMaintenanceURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
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

func AddClusterUsers(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	common.SetHeaders(client)
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(clusterUsersURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func DeleteClusterUsers(client *gophercloud.ServiceClient, id string, opts OptsQueryBuilder) (r DeleteResult) {
	common.SetHeaders(client)
	url := clusterUsersURL(client, id)
	query, err := opts.ToQuery()
	if err != nil {
		r.Err = err
		return
	}
	url += query

	resp, err := client.Delete(url, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func AddClusterConfigsWarehouseConnections(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	common.SetHeaders(client)
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(clusterConnectionsURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func DeleteClusterConnections(client *gophercloud.ServiceClient, clusterID string, opts OptsQueryBuilder) (r DeleteResult) {
	common.SetHeaders(client)
	url := clusterConnectionsURL(client, clusterID)
	query, err := opts.ToQuery()
	if err != nil {
		r.Err = err
		return
	}
	url += query

	resp, err := client.Delete(url, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func UpdateClusterPodGroup(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	common.SetHeaders(client)
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(clusterPodGroupsURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
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
