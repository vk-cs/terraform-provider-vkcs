package clusters

import (
	"github.com/gophercloud/gophercloud"
)

type (
	// Cluster represents a Kubernetes cluster in V2 API
	Cluster struct {
		ID                 string             `json:"id"`
		UUID               string             `json:"uuid"`
		Name               string             `json:"name"`
		Version            string             `json:"version"`
		CreatedAt          string             `json:"created_at"`
		Status             string             `json:"status"`
		NetworkConfig      NetworkConfig      `json:"network_config"`
		LoadBalancerConfig LoadBalancerConfig `json:"load_balancer_config"`
		Labels             map[string]string  `json:"labels"`
		NodeGroups         []NodeGroup        `json:"node_groups"`
		MasterSpec         MasterSpec         `json:"master"`
		DeploymentType     DeploymentType     `json:"deployment_type"`
		ProjectID          string             `json:"project_id"`
		ExternalIP         string             `json:"external_ip"`
		InternalIP         string             `json:"internal_ip"`
		ApiAddress         string             `json:"api_address"`
		InsecureRegistries []string           `json:"insecure_registries"`
		Description        string             `json:"description"`
	}

	// Cluster represents a cluster node group in V2 API
	NodeGroup struct {
		ID                   string            `json:"id"`
		Name                 string            `json:"name"`
		CreatedAt            string            `json:"created_at"`
		Zones                []string          `json:"zones"`
		ScaleSpec            ScaleSpec         `json:"scale_spec"`
		Labels               map[string]string `json:"labels"`
		Taints               []Taint           `json:"taints"`
		VMEngine             VMEngine          `json:"vm_engine"`
		ParallelUpgradeChunk int               `json:"parallel_upgrade_chunk"`
		ClusterID            string            `json:"cluster_id"`
		DiskType             DiskType          `json:"disk_type"`
		UUID                 string            `json:"uuid"`
	}

	// MasterSpec represents master node specification
	MasterSpec struct {
		Engine   MasterEngine     `json:"vm_engine"`
		Replicas int              `json:"replicas"`
		Disks    []MasterSpecDisk `json:"disks"`
	}

	// MasterSpec represents master node specification
	MasterSpecDisk struct {
		Type string `json:"type"`
		Size int    `json:"size"`
	}

	// MasterEngine represents master node engine configuration
	MasterEngine struct {
		NovaEngine NovaEngine `json:"nova_engine"`
	}

	// NovaEngine represents Nova-based engine configuration
	NovaEngine struct {
		FlavorID string `json:"flavor_id"`
	}

	// DeploymentType represents cluster deployment type
	DeploymentType struct {
		ZonalDeployment      *ZonalDeployment      `json:"zonal_deployment"`
		MultiZonalDeployment *MultiZonalDeployment `json:"multi_zonal_deployment"`
	}

	// ZonalDeployment represents single zone deployment
	ZonalDeployment struct {
		Zone string `json:"zone"`
	}

	// MultiZonalDeployment represents multi-zone deployment
	MultiZonalDeployment struct {
		Zones []string `json:"zones"`
	}

	// NetworkConfig represents network configuration
	NetworkConfig struct {
		Plugin NetworkPlugin `json:"plugin"`
		Engine NetworkEngine `json:"engine"`
	}

	// NetworkPlugin represents network plugin configuration
	NetworkPlugin struct {
		Calico *CalicoPlugin `json:"calico"`
		Cilium *CiliumPlugin `json:"cilium"`
	}

	// CalicoPlugin represents Calico network plugin configuration
	CalicoPlugin struct {
		PodsIPv4CIDR string `json:"pods_ipv4_cidr"`
	}

	// CiliumPlugin represents Cilium network plugin configuration
	CiliumPlugin struct {
		PodsIPv4CIDR string `json:"pods_ipv4_cidr"`
	}

	// NetworkEngine represents network engine configuration
	NetworkEngine struct {
		SprutEngine SprutEngine `json:"sprut_engine"`
	}

	// SprutEngine represents Sprut network engine configuration
	SprutEngine struct {
		NetworkID         string `json:"network_id"`
		SubnetID          string `json:"subnet_id"`
		ExternalNetworkID string `json:"external_network_id"`
	}

	// LoadBalancerConfig represents load balancer configuration
	LoadBalancerConfig struct {
		OctaviaEngine OctaviaEngine `json:"octavia_engine"`
	}

	// OctaviaEngine represents Octavia load balancer engine configuration
	OctaviaEngine struct {
		LoadbalancerSubnetID string   `json:"loadbalancer_subnet_id"`
		EnablePublicIP       bool     `json:"enable_public_ip"`
		AllowedCIDRs         []string `json:"allowed_cidrs"`
	}

	// ScaleSpec represents scaling specification
	ScaleSpec struct {
		FixedScale *FixedScale `json:"fixed_scale"`
		AutoScale  *AutoScale  `json:"auto_scale"`
	}

	// FixedScale represents fixed scaling configuration
	FixedScale struct {
		Size int `json:"size"`
	}

	// AutoScale represents auto scaling configuration
	AutoScale struct {
		MinSize int `json:"min_size"`
		MaxSize int `json:"max_size"`
		Size    int `json:"size"`
	}

	// Taint represents node taint configuration
	Taint struct {
		Key    string `json:"key"`
		Value  string `json:"value"`
		Effect string `json:"effect"`
	}

	// VMEngine represents virtual machine engine configuration
	VMEngine struct {
		NovaEngine NovaEngine `json:"nova_engine"`
	}

	// DiskType represents disk type configuration
	DiskType struct {
		CinderVolumeType CinderVolumeType `json:"cinder_volume_type"`
	}

	// CinderVolumeType represents Cinder volume type configuration
	CinderVolumeType struct {
		Type string `json:"type"`
		Size int    `json:"size"`
	}

	// ClusterID represents the result ID of a Create operation
	ClusterID struct {
		ID string `json:"id"`
	}

	// CreateResult represents the result of a Create operation
	CreateResult struct {
		gophercloud.Result
	}

	// GetResult represents the result of a Get operation
	GetResult struct {
		gophercloud.Result
	}
)

func (ng *NodeGroup) GetActualSize() int {
	if ng.ScaleSpec.AutoScale != nil {
		return ng.ScaleSpec.AutoScale.Size
	}
	if ng.ScaleSpec.FixedScale != nil {
		return ng.ScaleSpec.FixedScale.Size
	}
	return 0
}

// ExtractCreate extracts a cluster ID from a CreateResult
func (r CreateResult) Extract() (string, error) {
	var id ClusterID
	err := r.ExtractInto(&id)
	return id.ID, err
}

// ExtractGet extracts a cluster from a GetResult
func (r GetResult) Extract() (*Cluster, error) {
	var cluster Cluster
	err := r.ExtractInto(&cluster)
	return &cluster, err
}
