package clusters

import (
	"github.com/gophercloud/gophercloud"
)

// ClusterResp represents database cluster response
type Cluster struct {
	ID                string            `json:"id"`
	ClusterTemplateID string            `json:"cluster_template_id"`
	Configs           *ClusterConfig    `json:"configs"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	CreatedAt         string            `json:"created_at"`
	NetworkID         string            `json:"network_id"`
	PodGroups         []ClusterPodGroup `json:"pod_groups"`
	ProductName       string            `json:"product_name"`
	ProductType       string            `json:"product_type"`
	ProductVersion    string            `json:"product_version"`
	StackID           string            `json:"stack_id"`
	Status            string            `json:"status"`
	SubnetID          string            `json:"subnet_id"`
	Upgrades          []string          `json:"upgrades"`
	AvailabilityZone  string            `json:"availability_zone"`
	MultiAZ           bool              `json:"multi_az"`
	FloatingIPPool    string            `json:"floating_ip_pool"`
}

type ClusterConfig struct {
	Settings    []ClusterConfigSetting    `json:"settings"`
	Maintenance *ClusterConfigMaintenance `json:"maintenance"`
	Warehouses  []ClusterConfigWarehouse  `json:"warehouses"`
}

type ClusterConfigSetting struct {
	Alias string `json:"alias" required:"true"`
	Value string `json:"value" required:"true"`
}

type ClusterConfigMaintenance struct {
	Start    string                             `json:"start"`
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
	ID       string                 `json:"id"`
	Required bool                   `json:"required"`
	Name     string                 `json:"name"`
	Start    string                 `json:"start"`
	Settings []ClusterConfigSetting `json:"settings,omitempty"`
}

type ClusterConfigWarehouse struct {
	ID          string                             `json:"id"`
	Name        string                             `json:"name"`
	Connections []ClusterConfigWarehouseConnection `json:"connections"`
	Extensions  []ClusterConfigWarehouseExtension  `json:"extensions,omitempty"`
}

type ClusterConfigWarehouseConnection struct {
	Name      string                 `json:"name"`
	Plug      string                 `json:"plug"`
	Settings  []ClusterConfigSetting `json:"settings"`
	ID        string                 `json:"id"`
	CreatedAt string                 `json:"created_at"`
}

type ClusterConfigWarehouseExtension struct {
	Type      string                 `json:"type"`
	Version   string                 `json:"version,omitempty"`
	Settings  []ClusterConfigSetting `json:"settings,omitempty"`
	ID        string                 `json:"id"`
	CreatedAt string                 `json:"created_at"`
}

type ClusterPodGroup struct {
	ID                 string                           `json:"id"`
	Name               string                           `json:"name"`
	Count              int                              `json:"count"`
	Resource           *ClusterPodGroupResource         `json:"resource"`
	PodGroupTemplateID string                           `json:"pod_group_template_id" required:"true"`
	Volumes            map[string]ClusterPodGroupVolume `json:"volumes,omitempty"`
	FloatingIPPool     string                           `json:"floating_ip_pool,omitempty"`
	AvailabilityZone   string                           `json:"availability_zone,omitempty"`
	Alias              string                           `json:"alias,omitempty"`
}

type ClusterPodGroupResource struct {
	CPURequest string `json:"cpu_request"`
	CPULimit   string `json:"cpu_limit"`
	RAMRequest string `json:"ram_request"`
	RAMLimit   string `json:"ram_limit"`
}

type ClusterPodGroupVolume struct {
	StorageClassName string `json:"storageClassName"`
	Storage          string `json:"storage"`
	Count            int    `json:"count"`
}

type ClusterShortResp struct {
	ID string `json:"id"`
}

type commonClusterResult struct {
	gophercloud.Result
}

type commonShortClusterResult struct {
	gophercloud.Result
}

// CreateResult represents result of dataplatform cluster create
type CreateResult struct {
	commonShortClusterResult
}

type UpdateResult struct {
	commonShortClusterResult
}

// GetResult represents result of dataplatform cluster get
type GetResult struct {
	commonClusterResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}

// Extract is used to extract result into response struct
func (r commonClusterResult) Extract() (*Cluster, error) {
	var c *Cluster
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c, nil
}

// Extract is used to extract result into short response struct
func (r commonShortClusterResult) Extract() (*ClusterShortResp, error) {
	var c *ClusterShortResp
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c, nil
}
