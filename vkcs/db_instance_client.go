package vkcs

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
	db "github.com/gophercloud/gophercloud/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/pagination"
)

// databaseClient performs request to dbaas api
type databaseClient interface {
	Get(url string, JSONResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	Post(url string, JSONBody interface{}, JSONResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	Patch(url string, JSONBody interface{}, JSONResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	Delete(url string, opts *gophercloud.RequestOpts) (*http.Response, error)
	Head(url string, opts *gophercloud.RequestOpts) (*http.Response, error)
	Put(url string, JSONBody interface{}, JSONResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	ServiceURL(parts ...string) string
}

// instanceResp represents result of database instance get
type instanceResp struct {
	СomputeInstanceID string                  `json:"compute_instance_id"`
	Configuration     *configuration          `json:"configuration"`
	ConfigurationID   string                  `json:"configuration_id"`
	ID                string                  `json:"id"`
	Created           dateTimeWithoutTZFormat `json:"created"`
	Updated           dateTimeWithoutTZFormat `json:"updated"`
	DataStore         *dataStore              `json:"datastore"`
	Flavor            *links                  `json:"flavor"`
	GaVersion         string                  `json:"ga_version"`
	HealthStatus      string                  `json:"health_status"`
	IP                *[]string               `json:"ip"`
	Links             *[]link                 `json:"links"`
	Name              string                  `json:"name"`
	Region            string                  `json:"region"`
	Status            string                  `json:"status"`
	Volume            *volume                 `json:"volume"`
	ReplicaOf         *links                  `json:"replica_of"`
	AutoExpand        int                     `json:"volume_autoresize_enabled"`
	MaxDiskSize       int                     `json:"volume_autoresize_max_size"`
	WalVolume         *walVolume              `json:"wal_volume"`
}

type instanceShortResp struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// volume represents database instance volume
type volume struct {
	Size       *int     `json:"size" required:"true"`
	Used       *float32 `json:"used,omitempty"`
	VolumeID   string   `json:"volume_id,,omitempty"`
	VolumeType string   `json:"type,,omitempty" required:"true"`
}

// walVolume represents database instance wal volume
type walVolume struct {
	Size        *int     `json:"size" required:"true"`
	Used        *float32 `json:"used,omitempty"`
	VolumeID    string   `json:"volume_id,,omitempty"`
	VolumeType  string   `json:"type,,omitempty" required:"true"`
	AutoExpand  int      `json:"autoresize_enabled,omitempty"`
	MaxDiskSize int      `json:"autoresize_max_size,omitempty"`
}

// walVolumeOpts represents parameters for creation of database instance wal volume
type walVolumeOpts struct {
	Size       int
	VolumeType string `mapstructure:"volume_type"`
}

// links represents database instance links
type links struct {
	ID    string  `json:"id"`
	Links *[]link `json:"links"`
}

// dataStore represents dbaas datastore
type dataStore struct {
	Type    string `json:"type" required:"true"`
	Version string `json:"version" required:"true"`
}

// link represents database instance link
type link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

// configuration represents database instance configuration
type configuration struct {
	ID    string  `json:"id"`
	Links *[]link `json:"links"`
	Name  string  `json:"name"`
}

// instanceAutoExpandOpts represents autoresize parameters of volume of database instance
type instanceAutoExpandOpts struct {
	AutoExpand  bool
	MaxDiskSize int `mapstructure:"max_disk_size"`
}

// instanceDetachReplicaOpts represents parameters of request to detach replica of database instance
type instanceDetachReplicaOpts struct {
	Instance struct {
		ReplicaOf string `json:"replica_of,omitempty"`
	} `json:"instance"`
}

// instanceAttachConfigurationGroupOpts represents parameters of configuration group to be attached to database instance
type instanceAttachConfigurationGroupOpts struct {
	Instance struct {
		Configuration string `json:"configuration"`
	} `json:"instance"`
}

// instanceDetachConfigurationGroupOpts represents parameters of configuration group to be detached from database instance
type instanceDetachConfigurationGroupOpts struct {
	Instance map[string]interface{} `json:"instance"`
}

// instanceUpdateAutoExpandOpts represents parameters of request to update autoresize properties of volume of database instance
type instanceUpdateAutoExpandOpts struct {
	Instance struct {
		VolumeAutoresizeEnabled int `json:"volume_autoresize_enabled"`
		VolumeAutoresizeMaxSize int `json:"volume_autoresize_max_size"`
	} `json:"instance"`
}

// instanceResizeVolumeOpts represents database instance volume resize parameters
type instanceResizeVolumeOpts struct {
	Resize struct {
		Volume struct {
			Size int `json:"size"`
		} `json:"volume"`
	} `json:"resize"`
}

// instanceUpdateAutoExpandWalOpts represents parameters of request to update autoresize properties of wal volume of database instance
type instanceUpdateAutoExpandWalOpts struct {
	Instance struct {
		WalVolume struct {
			VolumeAutoresizeEnabled int `json:"autoresize_enabled"`
			VolumeAutoresizeMaxSize int `json:"autoresize_max_size"`
		} `json:"wal_volume"`
	} `json:"instance"`
}

// instanceResizeWalVolumeOpts represents database instance wal volume resize parameters
type instanceResizeWalVolumeOpts struct {
	Resize struct {
		Volume struct {
			Size int    `json:"size"`
			Kind string `json:"kind"`
		} `json:"volume"`
	} `json:"resize"`
}

// instanceResizeOpts represents database instance resize parameters
type instanceResizeOpts struct {
	Resize struct {
		FlavorRef string `json:"flavorRef"`
	} `json:"resize"`
}

// instanceRootUserEnableOpts represents parameters of request to enable root user for database instance
type instanceRootUserEnableOpts struct {
	Password string `json:"password,omitempty"`
}

// instanceRespOpts is used to get instance response
type instanceRespOpts struct {
	Instance *instanceResp `json:"instance"`
}

type instanceShortRespOpts struct {
	Instance *instanceShortResp `json:"instance"`
}

// instanceCapabilityOpts represents parameters of database instance capabilities
type instanceCapabilityOpts struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params,omitempty" mapstructure:"settings"`
}

// databaseCapability represents capability info from dbaas
type databaseCapability struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params,omitempty"`
	Status string            `json:"status"`
}

// instanceApplyCapabilityOpts is used to send request to apply capability to database instance
type instanceApplyCapabilityOpts struct {
	ApplyCapability struct {
		Capabilities []instanceCapabilityOpts `json:"capabilities"`
	} `json:"apply_capability"`
}

type instanceGetCapabilityOpts struct {
	Capabilities []databaseCapability `json:"capabilities"`
}

// userBatchCreateOpts is used to send request to create database users
type userBatchCreateOpts struct {
	Users []userCreateOpts `json:"users"`
}

// userCreateOpts represents parameters of creation of database user
type userCreateOpts struct {
	Name      string             `json:"name" required:"true"`
	Password  string             `json:"password" required:"true"`
	Databases db.BatchCreateOpts `json:"databases,omitempty"`
	Host      string             `json:"host,omitempty"`
}

// userUpdateOpts represents parameters of update of database user
type userUpdateOpts struct {
	User struct {
		Name     string `json:"name,omitempty"`
		Password string `json:"password,omitempty"`
		Host     string `json:"host,omitempty"`
	} `json:"user"`
}

// userUpdateDatabasesOpts represents parameters of request to update users databases
type userUpdateDatabasesOpts struct {
	Databases []map[string]string `json:"databases"`
}

// databaseBatchCreateOpts is used to send request to create databases
type databaseBatchCreateOpts struct {
	Databases []databaseCreateOpts `json:"databases"`
}

// databaseCreateOpts represents parameters of creation of database
type databaseCreateOpts struct {
	Name    string `json:"name" required:"true"`
	CharSet string `json:"character_set,omitempty"`
	Collate string `json:"collate,omitempty"`
}

// createOptsBuilder is used to build create opts map
type createOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// Map converts opts to a map (for a request body)
func (opts dbInstance) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Map converts opts to a map (for a request body)
func (opts *instanceDetachReplicaOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *instanceDetachConfigurationGroupOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *instanceAttachConfigurationGroupOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *instanceResizeVolumeOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *instanceResizeWalVolumeOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *instanceResizeOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *instanceRootUserEnableOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *instanceUpdateAutoExpandOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *instanceUpdateAutoExpandWalOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *instanceApplyCapabilityOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *userBatchCreateOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *userUpdateOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *userUpdateDatabasesOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map converts opts to a map (for a request body)
func (opts *databaseBatchCreateOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// dbInstance is used to send request to create database instance
type dbInstance struct {
	Instance *dbInstanceCreateOpts `json:"instance" required:"true"`
}

// dbInstanceCreateOpts represents parameters of creation of database instance
type dbInstanceCreateOpts struct {
	FlavorRef         string                   `json:"flavorRef,omitempty"`
	Volume            *volume                  `json:"volume" required:"true"`
	Name              string                   `json:"name" required:"true"`
	Datastore         *dataStore               `json:"datastore" required:"true"`
	Nics              []networkOpts            `json:"nics" required:"true"`
	ReplicaOf         string                   `json:"replica_of,omitempty"`
	AvailabilityZone  string                   `json:"availability_zone,omitempty"`
	FloatingIPEnabled bool                     `json:"allow_remote_access,omitempty"`
	Keypair           string                   `json:"key_name,omitempty"`
	AutoExpand        *int                     `json:"volume_autoresize_enabled,omitempty"`
	MaxDiskSize       int                      `json:"volume_autoresize_max_size,omitempty"`
	Walvolume         *walVolume               `json:"wal_volume,omitempty"`
	Capabilities      []instanceCapabilityOpts `json:"capabilities,omitempty"`
}

// networkOpts represents network parameters of database instance
type networkOpts struct {
	UUID      string `json:"net-id,omitempty"`
	Port      string `json:"port-id,omitempty"`
	V4FixedIP string `json:"v4-fixed-ip,omitempty" mapstructure:"fixed_ip_v4"`
}

type commonInstanceResult struct {
	gophercloud.Result
}

type commonInstanceCapabilitiesResult struct {
	gophercloud.Result
}

type instanceShortResult struct {
	gophercloud.Result
}

// getInstanceResult represents result of database instance get
type getInstanceResult struct {
	commonInstanceResult
}

type getInstanceShortResult struct {
	instanceShortResult
}

type getInstanceCapabilitiesResult struct {
	commonInstanceCapabilitiesResult
}

// rootUserResp represents parameters of root user response
type rootUserResp struct {
	Password string `json:"password"`
	Name     string `json:"name"`
}

// rootUserRespOpts is used to get root user response
type rootUserRespOpts struct {
	User *rootUserResp `json:"user"`
}

type commonRootUserResult struct {
	gophercloud.Result
}

type commonRootUserErrResult struct {
	gophercloud.ErrResult
}

// createRootUserResult represents result of root user create
type createRootUserResult struct {
	commonRootUserResult
}

// deleteRootUserResult represents result of root user delete
type deleteRootUserResult struct {
	commonRootUserErrResult
}

// configurationResult represents result of configuration attach and detach
type configurationResult struct {
	gophercloud.ErrResult
}

// actionResult represents result of database instance action
type actionResult struct {
	gophercloud.ErrResult
}

// isRootUserEnabledResult represents result of getting root user status
type isRootUserEnabledResult struct {
	commonRootUserResult
}

// deleteResult represents result of database instance delete
type deleteResult struct {
	gophercloud.ErrResult
}

type commonUserResult struct {
	gophercloud.ErrResult
}

// updateUserResult represents result of database user update
type updateUserResult struct {
	commonUserResult
}

// updateUserDatabasesResult represents result of database user database update
type updateUserDatabasesResult struct {
	commonUserResult
}

// deleteUserDatabaseResult represents result of database user delete
type deleteUserDatabaseResult struct {
	commonUserResult
}

// userCreateResult represents result of database user create
type userCreateResult struct {
	commonUserResult
}

// userDeleteResult represents result of database user delete
type userDeleteResult struct {
	commonUserResult
}

type commonDatabaseResult struct {
	gophercloud.ErrResult
}

// databaseCreateResult represents result of database create
type databaseCreateResult struct {
	commonDatabaseResult
}

// databaseDeleteResult represents result of database delete
type databaseDeleteResult struct {
	commonDatabaseResult
}

// extract is used to extract result into response struct
func (r commonInstanceResult) extract() (*instanceResp, error) {
	var i *instanceRespOpts
	if err := r.ExtractInto(&i); err != nil {
		return nil, err
	}
	return i.Instance, nil
}

func (r instanceShortResult) extract() (*instanceShortResp, error) {
	var i *instanceShortRespOpts
	if err := r.ExtractInto(&i); err != nil {
		return nil, err
	}
	return i.Instance, nil
}

// extract is used to extract result into response struct
func (r commonRootUserResult) extract() (*rootUserResp, error) {
	var u *rootUserRespOpts
	err := r.ExtractInto(&u)
	return u.User, err
}

// extract is used to extract result into response struct
func (r isRootUserEnabledResult) extract() (bool, error) {
	return r.Body.(map[string]interface{})["rootEnabled"].(bool), r.Err
}

func (r commonInstanceCapabilitiesResult) extract() ([]databaseCapability, error) {
	var c *instanceGetCapabilityOpts
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c.Capabilities, nil
}

var instancesAPIPath = "instances"

// instanceCreate performs request to create database instance
func instanceCreate(client databaseClient, opts optsBuilder) (r getInstanceShortResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	reqOpts := getDBRequestOpts(200)
	result, r.Err = client.Post(baseURL(client, instancesAPIPath), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// instanceGet performs request to get database instance
func instanceGet(client databaseClient, id string) (r getInstanceResult) {
	reqOpts := getDBRequestOpts(200)
	var result *http.Response
	result, r.Err = client.Get(getURL(client, instancesAPIPath, id), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func instanceGetCapabilities(client databaseClient, id string) (r getInstanceCapabilitiesResult) {
	reqOpts := getDBRequestOpts(200)
	var result *http.Response
	result, r.Err = client.Get(instanceCapabilitiesURL(client, instancesAPIPath, id), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func clusterGetCapabilities(client databaseClient, id string) (r getInstanceCapabilitiesResult) {
	reqOpts := getDBRequestOpts(200)
	var result *http.Response
	result, r.Err = client.Get(instanceCapabilitiesURL(client, dbClustersAPIPath, id), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// instanceDetachReplica performs request to detach replica of database instance
func instanceDetachReplica(client databaseClient, id string, opts optsBuilder) (r instances.ActionResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getDBRequestOpts(202)
	var result *http.Response
	result, r.Err = client.Patch(getURL(client, instancesAPIPath, id), b, nil, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// instanceUpdateAutoExpand performs request to update database instance autoresize parameters
func instanceUpdateAutoExpand(client databaseClient, id string, opts optsBuilder) (r instances.ActionResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getDBRequestOpts(202)
	var result *http.Response
	result, r.Err = client.Patch(getURL(client, instancesAPIPath, id), &b, nil, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// instanceDetachConfigurationGroup performs request to detach configuration group from database instance
func instanceDetachConfigurationGroup(client databaseClient, id string) (r configurationResult) {
	opts := instanceDetachConfigurationGroupOpts{
		Instance: map[string]interface{}{},
	}
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getDBRequestOpts(202)
	var result *http.Response
	result, r.Err = client.Put(getURL(client, instancesAPIPath, id), b, nil, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// instanceAttachConfigurationGroup performs request to attach configuration group to database instance
func instanceAttachConfigurationGroup(client databaseClient, id string, opts optsBuilder) (r configurationResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getDBRequestOpts(202)
	var result *http.Response
	result, r.Err = client.Put(getURL(client, instancesAPIPath, id), b, nil, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// instanceAction performs request to perform an action on the database instance
func instanceAction(client databaseClient, id string, opts optsBuilder) (r actionResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getDBRequestOpts(202)
	var result *http.Response
	result, r.Err = client.Post(instanceActionURL(client, instancesAPIPath, id), b, nil, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// instanceRootUserEnable performs request to enable root user on database instance
func instanceRootUserEnable(client databaseClient, id string, opts optsBuilder) (r createRootUserResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getDBRequestOpts(200)
	var result *http.Response
	result, r.Err = client.Post(rootUserURL(client, instancesAPIPath, id), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// instanceRootUserGet performs request to get root user of database instance
func instanceRootUserGet(client databaseClient, id string) (r isRootUserEnabledResult) {
	var result *http.Response
	reqOpts := getDBRequestOpts(200)
	result, err := client.Get(rootUserURL(client, instancesAPIPath, id), &r.Body, reqOpts)
	if err == nil {
		r.Header = result.Header
	}
	return
}

// instanceRootUserDisable performs request to disable root user on database instance
func instanceRootUserDisable(client databaseClient, id string) (r deleteRootUserResult) {
	reqOpts := getDBRequestOpts(200)
	var result *http.Response
	result, r.Err = client.Delete(rootUserURL(client, instancesAPIPath, id), reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// instanceDelete performs request to delete database instance
func instanceDelete(client databaseClient, id string) (r deleteResult) {
	reqOpts := getDBRequestOpts()
	var result *http.Response
	result, r.Err = client.Delete(getURL(client, instancesAPIPath, id), reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// userCreate performs request to create database user
func userCreate(client databaseClient, id string, opts createOptsBuilder, dbmsType string) (r userCreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getDBRequestOpts(202)
	var result *http.Response
	if dbmsType == dbmsTypeInstance {
		result, r.Err = client.Post(instanceUsersURL(client, instancesAPIPath, id), b, nil, reqOpts)
	} else {
		result, r.Err = client.Post(instanceUsersURL(client, dbClustersAPIPath, id), b, nil, reqOpts)
	}
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// userList performs request to get list of database users
func userList(client databaseClient, id string, dbmsType string) pagination.Pager {
	var APIPath string
	if dbmsType == dbmsTypeInstance {
		APIPath = instancesAPIPath
	} else {
		APIPath = dbClustersAPIPath
	}
	return pagination.NewPager(client.(*gophercloud.ServiceClient), instanceUsersURL(client, APIPath, id), func(r pagination.PageResult) pagination.Page {
		return DBUserPage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// userDelete performs request to delete database user
func userDelete(client databaseClient, id string, userName string, dbmsType string) (r userDeleteResult) {
	reqOpts := getDBRequestOpts()
	var result *http.Response
	var APIPath string
	if dbmsType == dbmsTypeInstance {
		APIPath = instancesAPIPath
	} else {
		APIPath = dbClustersAPIPath
	}
	result, r.Err = client.Delete(instanceUserURL(client, APIPath, id, userName), reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// userUpdate performs request to update database user
func userUpdate(client databaseClient, id string, name string, opts optsBuilder, dbmsType string) (r updateUserResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getDBRequestOpts(202)
	var APIPath string
	if dbmsType == dbmsTypeInstance {
		APIPath = instancesAPIPath
	} else {
		APIPath = dbClustersAPIPath
	}
	var result *http.Response
	result, r.Err = client.Put(userURL(client, APIPath, id, name), b, nil, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// userUpdateDatabases performs request to update database user databases
func userUpdateDatabases(client databaseClient, id string, name string, opts optsBuilder, dbmsType string) (r updateUserDatabasesResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getDBRequestOpts(202)
	var result *http.Response
	var APIPath string
	if dbmsType == dbmsTypeInstance {
		APIPath = instancesAPIPath
	} else {
		APIPath = dbClustersAPIPath
	}
	result, r.Err = client.Put(userDatabasesURL(client, APIPath, id, name), b, nil, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// userDeleteDatabase performs request to delete database user
func userDeleteDatabase(client databaseClient, id string, userName string, dbName string, dbmsType string) (r deleteUserDatabaseResult) {
	reqOpts := getDBRequestOpts(202)
	var result *http.Response
	var APIPath string
	if dbmsType == dbmsTypeInstance {
		APIPath = instancesAPIPath
	} else {
		APIPath = dbClustersAPIPath
	}
	result, r.Err = client.Delete(userDatabaseURL(client, APIPath, id, userName, dbName), reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// databaseCreate performs request to create database
func databaseCreate(client databaseClient, id string, opts createOptsBuilder, dbmsType string) (r databaseCreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getDBRequestOpts(202)
	var result *http.Response
	var APIPath string
	if dbmsType == dbmsTypeInstance {
		APIPath = instancesAPIPath
	} else {
		APIPath = dbClustersAPIPath
	}
	result, r.Err = client.Post(instanceDatabasesURL(client, APIPath, id), b, nil, reqOpts)

	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// databaseList performs request to list databases
func databaseList(client databaseClient, id string, dbmsType string) pagination.Pager {
	var APIPath string
	if dbmsType == dbmsTypeInstance {
		APIPath = instancesAPIPath
	} else {
		APIPath = dbClustersAPIPath
	}
	return pagination.NewPager(client.(*gophercloud.ServiceClient), instanceDatabasesURL(client, APIPath, id), func(r pagination.PageResult) pagination.Page {
		return DBPage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// databaseDelete performs request to delete database
func databaseDelete(client databaseClient, id string, dbName string, dbmsType string) (r databaseDeleteResult) {
	reqOpts := getDBRequestOpts()
	var result *http.Response
	var APIPath string
	if dbmsType == dbmsTypeInstance {
		APIPath = instancesAPIPath
	} else {
		APIPath = dbClustersAPIPath
	}
	result, r.Err = client.Delete(instanceDatabaseURL(client, APIPath, id, dbName), reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
