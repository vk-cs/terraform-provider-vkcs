package clusters

import (
	"bytes"
	"io"

	"github.com/gophercloud/gophercloud"
)

// CreateOpts represents options for creating a cluster in V2 API
type (
	CreateOpts struct {
		UUID               string                 `json:"uuid"`
		Name               string                 `json:"name"`
		Description        string                 `json:"description"`
		Version            string                 `json:"version"`
		Labels             map[string]string      `json:"labels"`
		MasterSpec         MasterSpecOpts         `json:"master_spec"`
		DeploymentType     DeploymentTypeOpts     `json:"deployment_type"`
		NetworkConfig      NetworkConfigOpts      `json:"network_config"`
		LoadBalancerConfig LoadBalancerConfigOpts `json:"load_balancer_config"`
		InsecureRegistries []string               `json:"insecure_registries"`
	}

	// MasterSpecOpts represents master node specification
	MasterSpecOpts struct {
		Engine   MasterEngineOpts `json:"engine"`
		Replicas int              `json:"replicas"`
	}

	// MasterEngineOpts represents master node engine configuration
	MasterEngineOpts struct {
		NovaEngine NovaEngineOpts `json:"nova_engine"`
	}

	// NovaEngineOpts represents Nova-based engine configuration
	NovaEngineOpts struct {
		FlavorID string `json:"flavor_id"`
	}

	// DeploymentTypeOpts represents cluster deployment type
	DeploymentTypeOpts struct {
		ZonalDeployment      *ZonalDeploymentOpts      `json:"zonal_deployment"`
		MultiZonalDeployment *MultiZonalDeploymentOpts `json:"multi_zonal_deployment"`
	}

	// ZonalDeploymentOpts represents single zone deployment
	ZonalDeploymentOpts struct {
		Zone string `json:"zone"`
	}

	// MultiZonalDeploymentOpts represents multi-zone deployment
	MultiZonalDeploymentOpts struct {
		Zones []string `json:"zones"`
	}

	// NetworkConfigOpts represents network configuration
	NetworkConfigOpts struct {
		Plugin NetworkPluginOpts `json:"plugin"`
		Engine NetworkEngineOpts `json:"engine"`
	}

	// NetworkPluginOpts represents network plugin configuration
	NetworkPluginOpts struct {
		Calico *CalicoPluginOpts `json:"calico"`
		Cilium *CiliumPluginOpts `json:"cilium"`
	}

	// NetworkEngineOpts represents network engine configuration
	NetworkEngineOpts struct {
		SprutEngine SprutEngineOpts `json:"sprut_engine"`
	}

	// CalicoPluginOpts represents Calico network plugin configuration
	CalicoPluginOpts struct {
		PodsIPv4CIDR string `json:"pods_ipv4_cidr"`
	}

	// CiliumPluginOpts represents Cilium network plugin configuration
	CiliumPluginOpts struct {
		PodsIPv4CIDR string `json:"pods_ipv4_cidr"`
	}

	// SprutEngineOpts represents Sprut network engine configuration
	SprutEngineOpts struct {
		NetworkID         string `json:"network_id"`
		SubnetID          string `json:"subnet_id"`
		ExternalNetworkID string `json:"external_network_id"`
	}

	// LoadBalancerConfigOpts represents load balancer configuration
	LoadBalancerConfigOpts struct {
		OctaviaEngine OctaviaEngineOpts `json:"octavia_engine"`
	}

	// OctaviaEngineOpts represents Octavia load balancer engine configuration
	OctaviaEngineOpts struct {
		LoadbalancerSubnetID string   `json:"loadbalancer_subnet_id"`
		EnablePublicIP       bool     `json:"enable_public_ip"`
		AllowedCIDRs         []string `json:"allowed_cidrs"`
	}

	// ScaleOpts represents options for scaling a cluster
	ScaleOpts struct {
		MasterSpec MasterSpecOpts `json:"master_spec"`
	}

	// UpgradeOpts represents options for scaling a cluster
	UpgradeOpts struct {
		Version string `json:"target_version"`
	}
)

// ToClusterCreateMap builds a request body from CreateOpts
func (opts CreateOpts) ToClusterCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// ToClusterUpdateMap builds a request body from ScaleOpts
func (opts ScaleOpts) ToClusterUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// ToClusterUpdateMap builds a request body from UpgradeOpts
func (opts UpgradeOpts) ToClusterUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create creates a new cluster
func Create(c *gophercloud.ServiceClient, opts *CreateOpts) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToClusterCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = c.Post(rootURL(c), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

// Get retrieves a specific cluster based on its unique ID
func Get(c *gophercloud.ServiceClient, clusterID string) GetResult {
	var res GetResult
	_, res.Err = c.Get(resourceURL(c, clusterID), &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

// Scale scales an existing cluster
func Scale(c *gophercloud.ServiceClient, clusterID string, opts ScaleOpts) error {
	reqBody, err := opts.ToClusterUpdateMap()
	if err != nil {
		return err
	}

	_, err = c.Post(scaleURL(c, clusterID), reqBody, nil, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return err
}

// Upgrade upgrades an existing cluster
func Upgrade(c *gophercloud.ServiceClient, clusterID string, opts UpgradeOpts) error {
	reqBody, err := opts.ToClusterUpdateMap()
	if err != nil {
		return err
	}

	_, err = c.Patch(upgradeURL(c, clusterID), reqBody, nil, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return err
}

// Delete deletes an existing cluster
func Delete(c *gophercloud.ServiceClient, clusterID string) error {
	_, err := c.Delete(resourceURL(c, clusterID), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return err
}

// Get a raw cluster kubeconfig
func GetKubeconfig(c *gophercloud.ServiceClient, clusterID string) (string, error) {
	resp, err := c.Post(kubeconfigURL(c, clusterID), nil, nil, &gophercloud.RequestOpts{
		OkCodes:          []int{200},
		KeepResponseBody: true,
	})
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if resp.ContentLength > 0 {
		buf.Grow(int(resp.ContentLength))
	}
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Get available Kubernetes versions for the cluster
func GetListK8SVersion(c *gophercloud.ServiceClient) GetListK8SVersionResult {
	var res GetListK8SVersionResult
	_, res.Err = c.Get(listKubeVersionURL(c), &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

// Get available volume types for worker node root volume
func GetVolumeTypes(c *gophercloud.ServiceClient) GetVolumeTypesResult {
	var res GetVolumeTypesResult
	_, res.Err = c.Get(volumeTypesURL(c), &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}
