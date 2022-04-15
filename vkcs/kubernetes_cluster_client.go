package vkcs

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clusters"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"
)

const magnumAPIMicroVersion = "1.24"

var magnumAPIMicroVersionHeader = map[string]string{
	"MCS-API-Version": fmt.Sprintf("container-infra %s", magnumAPIMicroVersion),
}

func addMagnumMicroVersionHeader(reqOpts *gophercloud.RequestOpts) {
	reqOpts.MoreHeaders = magnumAPIMicroVersionHeader
}

type node struct {
	Name        string     `json:"name"`
	UUID        string     `json:"uuid"`
	NodeGroupID string     `json:"node_group_id"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type nodesFlatSchema []map[string]interface{}

func flattenNodes(nodes []*node) nodesFlatSchema {
	flatSchema := nodesFlatSchema{}
	for _, node := range nodes {
		flatSchema = append(flatSchema, map[string]interface{}{
			"name":          node.Name,
			"uuid":          node.UUID,
			"node_group_id": node.NodeGroupID,
			"created_at":    getTimestamp(node.CreatedAt),
			"updated_at":    getTimestamp(node.UpdatedAt),
		})
	}
	return flatSchema
}

type nodeGroupClusterPatchOpts []nodeGroupPatchParams

type nodeGroupPatchParams struct {
	Path  string      `json:"path,omitempty"`
	Value interface{} `json:"value,omitempty"`
	Op    string      `json:"op,omitempty"`
}

type nodeGroupBatchAddParams struct {
	Action  string      `json:"action,omitempty"`
	Payload []nodeGroup `json:"payload,omitempty"`
}

type nodeGroup struct {
	Name              string    `json:"name,omitempty"`
	NodeCount         int       `json:"node_count,omitempty"`
	MaxNodes          int       `json:"max_nodes,omitempty"`
	MinNodes          int       `json:"min_nodes,omitempty"`
	VolumeSize        int       `json:"volume_size,omitempty"`
	VolumeType        string    `json:"volume_type,omitempty"`
	FlavorID          string    `json:"flavor_id,omitempty"`
	ImageID           string    `json:"image_id,omitempty"`
	Autoscaling       bool      `json:"autoscaling_enabled,omitempty"`
	ClusterID         string    `json:"cluster_id,omitempty"`
	UUID              string    `json:"uuid,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
	Nodes             []*node   `json:"nodes,omitempty"`
	State             string    `json:"state,omitempty"`
	AvailabilityZones []string  `json:"availability_zones"`
}

type nodeGroupLabel struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

type nodeGroupTaint struct {
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
	Effect string `json:"effect,omitempty"`
}

// nodeGroupCreateOpts contains options to create node group.
type nodeGroupCreateOpts struct {
	ClusterID         string           `json:"cluster_id" required:"true"`
	Name              string           `json:"name"`
	Labels            []nodeGroupLabel `json:"labels,omitempty"`
	Taints            []nodeGroupTaint `json:"taints,omitempty"`
	NodeCount         int              `json:"node_count,omitempty"`
	MaxNodes          int              `json:"max_nodes,omitempty"`
	MinNodes          int              `json:"min_nodes,omitempty"`
	VolumeSize        int              `json:"volume_size,omitempty"`
	VolumeType        string           `json:"volume_type,omitempty"`
	FlavorID          string           `json:"flavor_id,omitempty"`
	Autoscaling       bool             `json:"autoscaling_enabled,omitempty"`
	AvailabilityZones []string         `json:"availability_zones,omitempty"`
}

// nodeGroupScaleOpts contains options to scale node group
type nodeGroupScaleOpts struct {
	Delta    int    `json:"delta" required:"true"`
	Rollback string `json:"rollback,omitempty"`
}

// clusterCreateOpts contains options to create cluster
type clusterCreateOpts struct {
	ClusterTemplateID    string            `json:"cluster_template_id" required:"true"`
	Keypair              string            `json:"keypair,omitempty"`
	Labels               map[string]string `json:"labels,omitempty"`
	MasterCount          int               `json:"master_count,omitempty"`
	MasterFlavorID       string            `json:"master_flavor_id,omitempty"`
	Name                 string            `json:"name"`
	NetworkID            string            `json:"network_id" required:"true"`
	SubnetID             string            `json:"subnet_id" required:"true"`
	PodsNetworkCidr      string            `json:"pods_network_cidr,omitempty"`
	FloatingIPEnabled    bool              `json:"floating_ip_enabled"`
	APILBVIP             string            `json:"api_lb_vip,omitempty"`
	APILBFIP             string            `json:"api_lb_fip,omitempty"`
	RegistryAuthPassword string            `json:"registry_auth_password,omitempty"`
	AvailabilityZone     string            `json:"availability_zone,omitempty"`
	LoadbalancerSubnetID string            `json:"loadbalancer_subnet_id,omitempty"`
	InsecureRegistries   []string          `json:"insecure_registries,omitempty"`
}

type clusterActionsBaseOpts struct {
	Action  string      `json:"action" required:"true"`
	Payload interface{} `json:"payload,omitempty"`
}

type clusterUpgradeOpts struct {
	ClusterTemplateID string `json:"cluster_template_id" required:"true"`
	RollingEnabled    bool   `json:"rolling_enabled"`
}

type cluster struct {
	APIAddress           string             `json:"api_address"`
	ClusterTemplateID    string             `json:"cluster_template_id"`
	CreatedAt            time.Time          `json:"created_at"`
	DiscoveryURL         string             `json:"discovery_url"`
	KeyPair              string             `json:"keypair"`
	Labels               map[string]string  `json:"labels"`
	Links                []gophercloud.Link `json:"links"`
	MasterFlavorID       string             `json:"master_flavor_id"`
	MasterAddresses      []string           `json:"master_addresses"`
	MasterCount          int                `json:"master_count"`
	Name                 string             `json:"name"`
	ProjectID            string             `json:"project_id"`
	StackID              string             `json:"stack_id"`
	Status               clusterStatus      `json:"status"`
	NewStatus            clusterStatus      `json:"new_status"`
	StatusReason         string             `json:"status_reason"`
	UUID                 string             `json:"uuid"`
	UpdatedAt            time.Time          `json:"updated_at"`
	UserID               string             `json:"user_id"`
	NetworkID            string             `json:"network_id"`
	SubnetID             string             `json:"subnet_id"`
	PodsNetworkCidr      string             `json:"pods_network_cidr"`
	FloatingIPEnabled    bool               `json:"floating_ip_enabled"`
	APILBVIP             string             `json:"api_lb_vip"`
	APILBFIP             string             `json:"api_lb_fip"`
	IngressFloatingIP    string             `json:"ingress_floating_ip"`
	RegistryAuthPassword string             `json:"registry_auth_password"`
	AvailabilityZone     string             `json:"availability_zone"`
	LoadbalancerSubnetID string             `json:"loadbalancer_subnet_id"`
	InsecureRegistries   []string           `json:"insecure_registries,omitempty"`
}

type clusterTemplate struct {
	clustertemplates.ClusterTemplate
	DeprecatedAt time.Time `json:"deprecated_at"`
	Version      string    `json:"version"`
}

type clusterTemplates struct {
	Templates []clusterTemplate `json:"clustertemplates"`
}

// Map builds request params.
func (opts *clusterCreateOpts) Map() (map[string]interface{}, error) {
	cluster, err := gophercloud.BuildRequestBody(*opts, "")
	return cluster, err
}

// Map builds request params.
func (opts *clusterActionsBaseOpts) Map() (map[string]interface{}, error) {
	cluster, err := gophercloud.BuildRequestBody(*opts, "")
	return cluster, err
}

// Map builds request params.
func (opts *nodeGroupCreateOpts) Map() (map[string]interface{}, error) {
	cluster, err := gophercloud.BuildRequestBody(*opts, "")
	return cluster, err
}

// Map builds request params.
func (opts *nodeGroup) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map builds request params.
func (opts *nodeGroupScaleOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map builds request params.
func (opts *clusterUpgradeOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Map builds request params.
func (opts *nodeGroupBatchAddParams) Map() (map[string]interface{}, error) {
	batch, err := gophercloud.BuildRequestBody(*opts, "")
	return batch, err
}

// Map builds request params.
func (opts *nodeGroupPatchParams) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// PatchMap collects all the params.
func (opts *nodeGroupClusterPatchOpts) PatchMap() ([]map[string]interface{}, error) {
	var lb []map[string]interface{}
	for _, opt := range *opts {
		b, err := opt.Map()
		if err != nil {
			return nil, err
		}
		lb = append(lb, b)
	}
	return lb, nil
}

const (
	clustersAPIPath        = "clusters"
	nodeGroupsAPIPath      = "nodegroups"
	clusterTemplateAPIPath = "clustertemplates"
)

type commonResult struct {
	gophercloud.Result
}

type clusterConfigResult struct {
	commonResult
}

type nodeGroupResult struct {
	commonResult
}

type nodeGroupDeleteResult struct {
	gophercloud.ErrResult
}

type clusterDeleteResult struct {
	gophercloud.ErrResult
}

type nodeGroupScaleResult struct {
	commonResult
}

type clusterTemplateResult struct {
	commonResult
}

// Extract returns uuid.
func (r nodeGroupScaleResult) Extract() (string, error) {
	var s struct {
		UUID string
	}
	err := r.ExtractInto(&s)
	return s.UUID, err
}

// Extract parses result into params for cluster.
func (r clusterConfigResult) Extract() (*cluster, error) {
	var s *cluster
	err := r.ExtractInto(&s)
	return s, err
}

// Extract parses result into params for cluster template.
func (r clusterTemplateResult) Extract() (*clusterTemplate, error) {
	var s *clusterTemplate
	err := r.ExtractInto(&s)
	return s, err
}

type clusterTemplatesResult struct {
	commonResult
}

// Extract parses result into params for cluster templates.
func (r clusterTemplatesResult) Extract() ([]clusterTemplate, error) {
	var s *clusterTemplates
	err := r.ExtractInto(&s)
	return s.Templates, err
}

// Extract parses result into params for node group.
func (r nodeGroupResult) Extract() (*nodeGroup, error) {
	var s *nodeGroup
	err := r.ExtractInto(&s)
	return s, err
}

func clusterTemplateGet(client ContainerClient, id string) (r clusterTemplateResult) {
	var result *http.Response
	reqOpts := getRequestOpts(200)
	result, r.Err = client.Get(getURL(client, clusterTemplateAPIPath, id), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func clusterTemplateList(client ContainerClient) (r clusterTemplatesResult) {
	var result *http.Response
	reqOpts := getRequestOpts(200)
	result, r.Err = client.Get(getURL(client, clusterTemplateAPIPath, ""), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func clusterCreate(client ContainerClient, opts optsBuilder) (r clusters.CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	reqOpts := getRequestOpts(202)
	result, r.Err = client.Post(baseURL(client, clustersAPIPath), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func clusterUpgrade(client ContainerClient, id string, opts optsBuilder) (r clusters.UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getRequestOpts(200, 202)
	var result *http.Response
	result, r.Err = client.Patch(upgradeURL(client, clustersAPIPath, id), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func clusterUpdateMasters(client ContainerClient, id string, opts optsBuilder) (r clusters.UpdateResult) {
	log.Printf("UPDATE masters for cluster %s", id)
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getRequestOpts(200, 202)
	var result *http.Response
	result, r.Err = client.Post(actionsURL(client, clustersAPIPath, id), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func clusterSwitchState(client ContainerClient, id string, opts optsBuilder) (r clusters.UpdateResult) {
	reqBody, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getRequestOpts(202)
	var result *http.Response
	result, r.Err = client.Post(actionsURL(client, clustersAPIPath, id), reqBody, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// clusterGet gets cluster data from vkcs.
func clusterGet(client ContainerClient, id string) (r clusterConfigResult) {
	log.Printf("GET cluster %s", id)
	var result *http.Response
	reqOpts := getRequestOpts(200)
	result, r.Err = client.Get(getURL(client, clustersAPIPath, id), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func clusterDelete(client ContainerClient, id string) (r clusterDeleteResult) {
	var result *http.Response
	reqOpts := getRequestOpts()
	result, r.Err = client.Delete(deleteURL(client, clustersAPIPath, id), reqOpts)
	r.Header = result.Header
	return
}

func nodeGroupGet(client ContainerClient, id string) (r nodeGroupResult) {
	var result *http.Response
	reqOpts := getRequestOpts(200)
	result, r.Err = client.Get(getURL(client, nodeGroupsAPIPath, id), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func nodeGroupScale(client ContainerClient, id string, opts optsBuilder) (r nodeGroupResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getRequestOpts(202)
	var result *http.Response
	result, r.Err = client.Patch(scaleURL(client, nodeGroupsAPIPath, id), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func nodeGroupCreate(client ContainerClient, opts optsBuilder) (r nodeGroupResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	reqOpts := getRequestOpts(202)
	result, r.Err = client.Post(baseURL(client, nodeGroupsAPIPath), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func nodeGroupDelete(client ContainerClient, id string) (r nodeGroupDeleteResult) {
	var result *http.Response
	reqOpts := getRequestOpts(204)
	result, r.Err = client.Delete(getURL(client, nodeGroupsAPIPath, id), reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func k8sConfigGet(client ContainerClient, id string) (string, error) {
	var result *http.Response
	reqOpts := getRequestOpts(200)
	reqOpts.KeepResponseBody = true
	result, err := client.Get(kubeConfigURL(client, clustersAPIPath, id), nil, reqOpts)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(make([]byte, 0, result.ContentLength))
	_, err = io.Copy(buf, result.Body)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func nodeGroupPatch(client ContainerClient, id string, opts patchOptsBuilder) (r nodeGroupScaleResult) {
	b, err := opts.PatchMap()
	if err != nil {
		r.Err = err
	}
	var result *http.Response
	reqOpts := getRequestOpts(200)
	result, r.Err = client.Patch(getURL(client, nodeGroupsAPIPath, id), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
