package vkcs

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

// dbCluster is used to send proper request for cluster creation
type dbCluster struct {
	Cluster *dbClusterCreateOpts `json:"cluster" required:"true"`
}

// dbClusterCreateOpts represents database cluster creation parameters
type dbClusterCreateOpts struct {
	Name                   string                        `json:"name" required:"true"`
	Datastore              *dataStoreShort               `json:"datastore" required:"true"`
	FloatingIPEnabled      bool                          `json:"allow_remote_access,omitempty"`
	AutoExpand             int                           `json:"volume_autoresize_enabled,omitempty"`
	MaxDiskSize            int                           `json:"volume_autoresize_max_size,omitempty"`
	WalAutoExpand          int                           `json:"wal_autoresize_enabled,omitempty"`
	WalMaxDiskSize         int                           `json:"wal_autoresize_max_size,omitempty"`
	Instances              []dbClusterInstanceCreateOpts `json:"instances"`
	Capabilities           []instanceCapabilityOpts      `json:"capabilities,omitempty"`
	RestorePoint           *restorePoint                 `json:"restorePoint,omitempty"`
	BackupSchedule         *backupSchedule               `json:"backup_schedule,omitempty"`
	CloudMonitoringEnabled bool                          `json:"cloud_monitoring_enabled,omitempty"`
}

// dbClusterInstanceCreateOpts represents database cluster instance creation parameters
type dbClusterInstanceCreateOpts struct {
	Keypair          string        `json:"key_name,omitempty"`
	AvailabilityZone string        `json:"availability_zone,omitempty"`
	FlavorRef        string        `json:"flavorRef,omitempty" mapstructure:"flavor_id"`
	Nics             []networkOpts `json:"nics" required:"true"`
	Volume           *volume       `json:"volume" required:"true"`
	Walvolume        *walVolume    `json:"wal_volume,omitempty"`
	ShardID          string        `json:"shard_id,omitempty"`
	SecurityGroups   []string      `json:"security_groups,omitempty"`
}

// dbClusterAttachConfigurationGroupOpts represents parameters of configuration group to be attached to database cluster
type dbClusterAttachConfigurationGroupOpts struct {
	ConfigurationAttach struct {
		ConfigurationID string `json:"configuration_id"`
	} `json:"configuration_attach"`
}

// dbClusterDetachConfigurationGroupOpts represents parameters of configuration group to be detached from database cluster
type dbClusterDetachConfigurationGroupOpts struct {
	ConfigurationDetach struct {
		ConfigurationID string `json:"configuration_id"`
	} `json:"configuration_detach"`
}

// dbClusterResizeVolumeOpts represents parameters of volume resize of database cluster
type dbClusterResizeVolumeOpts struct {
	Resize struct {
		Volume struct {
			Size int `json:"size"`
		} `json:"volume"`
	} `json:"resize"`
}

// dbClusterResizeWalVolumeOpts represents parameters of wal volume resize of database cluster
type dbClusterResizeWalVolumeOpts struct {
	Resize struct {
		Volume struct {
			Size int    `json:"size"`
			Kind string `json:"kind"`
		} `json:"volume"`
	} `json:"resize"`
}

// dbClusterResizeOpts represents database cluster resize parameters
type dbClusterResizeOpts struct {
	Resize struct {
		FlavorRef string `json:"flavorRef"`
	} `json:"resize"`
}

// dbClusterUpdateAutoExpandOpts represents autoresize parameters of volume of database cluster
type dbClusterUpdateAutoExpandOpts struct {
	Cluster struct {
		VolumeAutoresizeEnabled int `json:"volume_autoresize_enabled"`
		VolumeAutoresizeMaxSize int `json:"volume_autoresize_max_size"`
	} `json:"cluster"`
}

// dbClusterUpdateAutoExpandWalOpts represents autoresize parameters of wal volume of database cluster
type dbClusterUpdateAutoExpandWalOpts struct {
	Cluster struct {
		WalVolume struct {
			VolumeAutoresizeEnabled int `json:"autoresize_enabled"`
			VolumeAutoresizeMaxSize int `json:"autoresize_max_size"`
		} `json:"wal_volume"`
	} `json:"cluster"`
}

// dbClusterApplyCapabilityOpts represents parameters of capabilities to be applied to database cluster
type dbClusterApplyCapabilityOpts struct {
	ApplyCapability struct {
		Capabilities []instanceCapabilityOpts `json:"capabilities"`
	} `json:"apply_capability"`
}

// dbClusterGrowClusterOpts is used to send proper request to grow cluster
type dbClusterGrowClusterOpts struct {
	Grow []dbClusterGrowOpts `json:"grow"`
}

// dbClusterGrowOpts represents parameters of growing cluster
type dbClusterGrowOpts struct {
	Keypair          string     `json:"key_name"`
	AvailabilityZone string     `json:"availability_zone" required:"true"`
	FlavorRef        string     `json:"flavorRef" required:"true"`
	Volume           *volume    `json:"volume" required:"true"`
	Walvolume        *walVolume `json:"wal_volume,omitempty"`
}

// dbClusterShrinkClusterOpts is used to send proper request to shrink database cluster
type dbClusterShrinkClusterOpts struct {
	Shrink []dbClusterShrinkOpts `json:"shrink"`
}

// dbClusterShrinkOpts represents parameters of shrinking database cluster
type dbClusterShrinkOpts struct {
	ID string `json:"id" required:"true"`
}

// dbClusterResp represents database cluster response
type dbClusterResp struct {
	ConfigurationID string                  `json:"configuration_id"`
	Created         dateTimeWithoutTZFormat `json:"created"`
	DataStore       *dataStoreShort         `json:"datastore"`
	HealthStatus    string                  `json:"health_status"`
	ID              string                  `json:"id"`
	Instances       []dbClusterInstanceResp `json:"instances"`
	Links           *[]link                 `json:"links"`
	LoadbalancerID  string                  `json:"loadbalancer_id"`
	Name            string                  `json:"name"`
	Task            dbClusterTask           `json:"task"`
	Updated         dateTimeWithoutTZFormat `json:"updated"`
	AutoExpand      int                     `json:"volume_autoresize_enabled"`
	MaxDiskSize     int                     `json:"volume_autoresize_max_size"`
	WalAutoExpand   int                     `json:"wal_autoresize_enabled"`
	WalMaxDiskSize  int                     `json:"wal_autoresize_max_size"`
}

// dbClusterInstanceResp represents database cluster instance response
type dbClusterInstanceResp struct {
	СomputeInstanceID string     `json:"compute_instance_id"`
	Flavor            *links     `json:"flavor"`
	GaVersion         string     `json:"ga_version"`
	ID                string     `json:"id"`
	IP                *[]string  `json:"ip"`
	Links             *[]link    `json:"links"`
	Name              string     `json:"name"`
	Role              string     `json:"role"`
	Status            string     `json:"status"`
	Type              string     `json:"type"`
	Volume            *volume    `json:"volume"`
	WalVolume         *walVolume `json:"wal_volume"`
	ShardID           string     `json:"shard_id"`
}

type dbClusterShortResp struct {
	ID string `json:"id"`
}

// dbClusterRespOpts is used to properly extract database cluster response
type dbClusterRespOpts struct {
	Cluster *dbClusterResp `json:"cluster"`
}

type dbClusterShortRespOpts struct {
	Cluster *dbClusterShortResp `json:"cluster"`
}

// dbClusterTask represents database cluster task
type dbClusterTask struct {
	Description string `json:"description"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
}

// Map converts opts to a map (for a request body)
func (opts dbCluster) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Map converts opts to a map (for a request body)
func (opts *dbClusterDetachConfigurationGroupOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *dbClusterAttachConfigurationGroupOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *dbClusterResizeVolumeOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *dbClusterResizeWalVolumeOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *dbClusterResizeOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *dbClusterUpdateAutoExpandOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *dbClusterUpdateAutoExpandWalOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *dbClusterApplyCapabilityOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *dbClusterGrowClusterOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *dbClusterShrinkClusterOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

type commonClusterResult struct {
	gophercloud.Result
}

type shortClusterResult struct {
	gophercloud.Result
}

// createClusterResult represents result of database cluster create
type createClusterResult struct {
	shortClusterResult
}

// getClusterResult represents result of database cluster get
type getClusterResult struct {
	commonClusterResult
}

type deleteClusterResult struct {
	gophercloud.ErrResult
}

// clusterActionResult represents result of database cluster action
type clusterActionResult struct {
	gophercloud.ErrResult
}

// extract is used to extract result into response struct
func (r commonClusterResult) extract() (*dbClusterResp, error) {
	var c *dbClusterRespOpts
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c.Cluster, nil
}

// extract is used to extract result into short response struct
func (r shortClusterResult) extract() (*dbClusterShortResp, error) {
	var c *dbClusterShortRespOpts
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c.Cluster, nil
}

var dbClustersAPIPath = "clusters"

// dbClusterCreate performs request to create database cluster
func dbClusterCreate(client databaseClient, opts optsBuilder) (r createClusterResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	reqOpts := getRequestOpts(200)
	result, r.Err = client.Post(baseURL(client, dbClustersAPIPath), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// dbClusterGet performs request to get database cluster
func dbClusterGet(client databaseClient, id string) (r getClusterResult) {
	reqOpts := getRequestOpts(200)
	var result *http.Response
	result, r.Err = client.Get(getURL(client, dbClustersAPIPath, id), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func dbClusterDelete(client databaseClient, id string) (r deleteClusterResult) {
	reqOpts := getRequestOpts()
	var result *http.Response
	result, r.Err = client.Delete(getURL(client, dbClustersAPIPath, id), reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// dbClusterAction performs request to perform an action on the database cluster
func dbClusterAction(client databaseClient, id string, opts optsBuilder) (r clusterActionResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getRequestOpts(202)
	var result *http.Response
	result, r.Err = client.Post(getURL(client, dbClustersAPIPath, id), b, nil, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// dbClusterUpdateAutoExpand performs request to update database cluster autoresize parameters
func dbClusterUpdateAutoExpand(client databaseClient, id string, opts optsBuilder) (r clusterActionResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getRequestOpts(202)
	var result *http.Response
	result, r.Err = client.Patch(getURL(client, dbClustersAPIPath, id), b, nil, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func dbClusterUpdateBackupSchedule(client databaseClient, id string, opts optsBuilder) (r updateBackupScheduleResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getRequestOpts(200)
	var result *http.Response
	result, r.Err = client.Put(backupScheduleURL(client, clustersAPIPath, id), b, nil, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func dbClusterGetBackupSchedule(client databaseClient, id string) (r getClusterBackupScheduleResult) {
	reqOpts := getRequestOpts(200)
	var result *http.Response
	result, r.Err = client.Get(backupScheduleURL(client, clustersAPIPath, id), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
