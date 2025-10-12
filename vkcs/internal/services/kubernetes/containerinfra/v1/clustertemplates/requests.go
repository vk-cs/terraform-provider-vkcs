package clustertemplates

import (
	"net/http"
	"time"

	"github.com/gophercloud/gophercloud"
)

// Represents a template for a Cluster Template
type ClusterTemplate struct {
	APIServerPort       int                `json:"apiserver_port"`
	COE                 string             `json:"coe"`
	ClusterDistro       string             `json:"cluster_distro"`
	CreatedAt           time.Time          `json:"created_at"`
	DNSNameServer       string             `json:"dns_nameserver"`
	DockerStorageDriver string             `json:"docker_storage_driver"`
	DockerVolumeSize    int                `json:"docker_volume_size"`
	ExternalNetworkID   string             `json:"external_network_id"`
	FixedNetwork        string             `json:"fixed_network"`
	FixedSubnet         string             `json:"fixed_subnet"`
	FlavorID            string             `json:"flavor_id"`
	FloatingIPEnabled   bool               `json:"floating_ip_enabled"`
	HTTPProxy           string             `json:"http_proxy"`
	HTTPSProxy          string             `json:"https_proxy"`
	ImageID             string             `json:"image_id"`
	InsecureRegistry    string             `json:"insecure_registry"`
	KeyPairID           string             `json:"keypair_id"`
	Labels              map[string]string  `json:"labels"`
	Links               []gophercloud.Link `json:"links"`
	MasterFlavorID      string             `json:"master_flavor_id"`
	MasterLBEnabled     bool               `json:"master_lb_enabled"`
	Name                string             `json:"name"`
	NetworkDriver       string             `json:"network_driver"`
	NoProxy             string             `json:"no_proxy"`
	ProjectID           string             `json:"project_id"`
	Public              bool               `json:"public"`
	RegistryEnabled     bool               `json:"registry_enabled"`
	ServerType          string             `json:"server_type"`
	TLSDisabled         bool               `json:"tls_disabled"`
	UUID                string             `json:"uuid"`
	UpdatedAt           time.Time          `json:"updated_at"`
	UserID              string             `json:"user_id"`
	VolumeDriver        string             `json:"volume_driver"`
	DeprecatedAt        time.Time          `json:"deprecated_at"`
	Version             string             `json:"version"`
}

type ClusterTemplates struct {
	Templates []ClusterTemplate `json:"clustertemplates"`
}

func Get(client *gophercloud.ServiceClient, id string) (r clusterTemplateResult) {
	var result *http.Response
	result, r.Err = client.Get(templateURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

func List(client *gophercloud.ServiceClient) (r clusterTemplatesResult) {
	var result *http.Response
	result, r.Err = client.Get(templatesURL(client), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
