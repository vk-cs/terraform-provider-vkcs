package clusters

import (
	"github.com/gophercloud/gophercloud"
)

type (
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

	MasterSpec struct {
		Engine   MasterEngine     `json:"vm_engine"`
		Replicas int              `json:"replicas"`
		Disks    []MasterSpecDisk `json:"disks"`
	}

	MasterSpecDisk struct {
		Type string `json:"type"`
		Size int    `json:"size"`
	}

	MasterEngine struct {
		NovaEngine NovaEngine `json:"nova_engine"`
	}

	NovaEngine struct {
		FlavorID string `json:"flavor_id"`
	}

	DeploymentType struct {
		ZonalDeployment      *ZonalDeployment      `json:"zonal_deployment"`
		MultiZonalDeployment *MultiZonalDeployment `json:"multi_zonal_deployment"`
	}

	ZonalDeployment struct {
		Zone string `json:"zone"`
	}

	MultiZonalDeployment struct {
		Zones []string `json:"zones"`
	}

	NetworkConfig struct {
		Plugin NetworkPlugin `json:"plugin"`
		Engine NetworkEngine `json:"engine"`
	}

	NetworkPlugin struct {
		Calico *CalicoPlugin `json:"calico"`
		Cilium *CiliumPlugin `json:"cilium"`
	}

	CalicoPlugin struct {
		PodsIPv4CIDR string `json:"pods_ipv4_cidr"`
	}

	CiliumPlugin struct {
		PodsIPv4CIDR string `json:"pods_ipv4_cidr"`
	}

	NetworkEngine struct {
		SprutEngine SprutEngine `json:"sprut_engine"`
	}

	SprutEngine struct {
		NetworkID         string `json:"network_id"`
		SubnetID          string `json:"subnet_id"`
		ExternalNetworkID string `json:"external_network_id"`
	}

	LoadBalancerConfig struct {
		OctaviaEngine OctaviaEngine `json:"octavia_engine"`
	}

	OctaviaEngine struct {
		LoadbalancerSubnetID string   `json:"loadbalancer_subnet_id"`
		EnablePublicIP       bool     `json:"enable_public_ip"`
		AllowedCIDRs         []string `json:"allowed_cidrs"`
	}

	ScaleSpec struct {
		FixedScale *FixedScale `json:"fixed_scale"`
		AutoScale  *AutoScale  `json:"auto_scale"`
	}

	FixedScale struct {
		Size int `json:"size"`
	}

	AutoScale struct {
		MinSize int `json:"min_size"`
		MaxSize int `json:"max_size"`
		Size    int `json:"size"`
	}

	Taint struct {
		Key    string `json:"key"`
		Value  string `json:"value"`
		Effect string `json:"effect"`
	}

	VMEngine struct {
		NovaEngine NovaEngine `json:"nova_engine"`
	}

	DiskType struct {
		CinderVolumeType CinderVolumeType `json:"cinder_volume_type"`
	}

	CinderVolumeType struct {
		Type string `json:"type"`
		Size int    `json:"size"`
	}

	ClusterID struct {
		ID string `json:"id"`
	}

	ListClusterAZ struct {
		AZs []string `json:"azs"`
	}

	ListK8SVersion struct {
		Versions []K8SVersion `json:"k8s_versions"`
	}

	ListVolumeTypes struct {
		StorageClasses []StorageClass `json:"storage_classes"`
	}

	StorageClass struct {
		Name  string   `json:"name"`
		Zones []string `json:"zones"`
	}

	K8SVersion struct {
		Version string `json:"version"`
	}

	CreateResult struct {
		gophercloud.Result
	}

	GetResult struct {
		gophercloud.Result
	}

	GetAZsResult struct {
		gophercloud.Result
	}

	GetListK8SVersionResult struct {
		gophercloud.Result
	}

	GetVolumeTypesResult struct {
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

func (r CreateResult) Extract() (string, error) {
	var id ClusterID
	err := r.ExtractInto(&id)
	return id.ID, err
}

func (r GetResult) Extract() (*Cluster, error) {
	var cluster Cluster
	err := r.ExtractInto(&cluster)
	return &cluster, err
}

func (r GetAZsResult) Extract() (ListClusterAZ, error) {
	var azs ListClusterAZ
	err := r.ExtractInto(&azs)
	return azs, err
}

func (r GetListK8SVersionResult) Extract() (ListK8SVersion, error) {
	var k8sListVersions ListK8SVersion
	err := r.ExtractInto(&k8sListVersions)
	return k8sListVersions, err
}

func (r GetVolumeTypesResult) Extract() (ListVolumeTypes, error) {
	var azs ListVolumeTypes
	err := r.ExtractInto(&azs)
	return azs, err
}

// Less reports whether d is less than other.
// Ordering: Type (lexicographic), then Size.
func (d MasterSpecDisk) Less(other MasterSpecDisk) bool {
	if d.Type < other.Type {
		return true
	}
	if d.Type > other.Type {
		return false
	}
	return d.Size < other.Size
}

// Less reports whether ng is less than other.
// Ordering: UUID (lexicographic).
func (ng NodeGroup) Less(other NodeGroup) bool {
	return ng.UUID < other.UUID
}
