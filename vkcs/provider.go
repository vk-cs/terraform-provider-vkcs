package vkcs

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/terraform/auth"
	"github.com/gophercloud/utils/terraform/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
)

const (
	maxRetriesCount         = 3
	defaultIdentityEndpoint = "https://infra.mail.ru/identity/v3/"
	defaultUsersDomainName  = "users"
	requestsMaxRetriesCount = 3
	requestsRetryDelay      = 30 * time.Millisecond
)

// configer is interface to work with gophercloud.Config calls
type configer interface {
	LoadAndValidate() error
	GetRegion() string
	ComputeV2Client(region string) (*gophercloud.ServiceClient, error)
	ImageV2Client(region string) (*gophercloud.ServiceClient, error)
	NetworkingV2Client(region string, sdn string) (*gophercloud.ServiceClient, error)
	BlockStorageV3Client(region string) (*gophercloud.ServiceClient, error)
	KeyManagerV1Client(region string) (*gophercloud.ServiceClient, error)
	ContainerInfraV1Client(region string) (ContainerClient, error)
	IdentityV3Client(region string) (ContainerClient, error)
	DatabaseV1Client(region string) (*gophercloud.ServiceClient, error)
	GetMutex() *mutexkv.MutexKV
}

// config uses openstackbase.Config as the base/foundation of this provider's
type config struct {
	auth.Config
}

var _ configer = &config{}

// GetRegion is implementation of getRegion method
func (c *config) GetRegion() string {
	return c.Region
}

func (c *config) ComputeV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.Config.ComputeV2Client(region)
}

func (c *config) ImageV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.Config.ImageV2Client(region)
}

func (c *config) NetworkingV2Client(region string, sdn string) (*gophercloud.ServiceClient, error) {
	client, err := c.Config.NetworkingV2Client(region)
	if err != nil {
		return client, err
	}
	client.MoreHeaders = map[string]string{
		"X-SDN": sdn,
	}
	return client, err
}

func (c *config) BlockStorageV3Client(region string) (*gophercloud.ServiceClient, error) {
	return c.Config.BlockStorageV3Client(region)
}

func (c *config) KeyManagerV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.Config.KeyManagerV1Client(region)
}

// DatabaseV1Client is implementation of DatabaseV1Client method
func (c *config) DatabaseV1Client(region string) (*gophercloud.ServiceClient, error) {
	client, clientErr := c.Config.DatabaseV1Client(region)
	client.ProviderClient.RetryFunc = func(context context.Context, method, url string, options *gophercloud.RequestOpts, err error, failCount uint) error {
		if failCount >= requestsMaxRetriesCount {
			return err
		}
		switch errType := err.(type) {
		case gophercloud.ErrDefault500, gophercloud.ErrDefault503:
			time.Sleep(requestsRetryDelay)
			return nil
		case gophercloud.ErrUnexpectedResponseCode:
			if errType.Actual == http.StatusGatewayTimeout {
				time.Sleep(requestsRetryDelay)
				return nil
			}
			return err
		default:
			return err
		}
	}
	return client, clientErr
}

// ContainerInfraV1Client is implementation of ContainerInfraV1Client method
func (c *config) ContainerInfraV1Client(region string) (ContainerClient, error) {
	return c.Config.ContainerInfraV1Client(region)
}

// IdentityV3Client is implementation of ContainerInfraV1Client method
func (c *config) IdentityV3Client(region string) (ContainerClient, error) {
	return c.Config.IdentityV3Client(region)
}

func (c *config) GetMutex() *mutexkv.MutexKV {
	return c.Config.MutexKV
}

func newConfig(d *schema.ResourceData, terraformVersion string) (configer, diag.Diagnostics) {
	config := &config{
		auth.Config{
			CACertFile:       d.Get("cacert_file").(string),
			ClientCertFile:   d.Get("cert").(string),
			ClientKeyFile:    d.Get("key").(string),
			Password:         d.Get("password").(string),
			TenantID:         d.Get("project_id").(string),
			Region:           d.Get("region").(string),
			IdentityEndpoint: d.Get("auth_url").(string),
			AllowReauth:      true,
			MaxRetries:       maxRetriesCount,
			TerraformVersion: terraformVersion,
			SDKVersion:       meta.SDKVersionString(),
			MutexKV:          mutexkv.NewMutexKV(),
		},
	}

	if config.UserDomainID == "" {
		config.UserDomainID = os.Getenv("OS_USER_DOMAIN_ID")
	}

	v, ok := d.GetOk("insecure")
	if ok {
		insecure := v.(bool)
		config.Insecure = &insecure
	}

	if err := initWithUsername(d, config); err != nil {
		return nil, diag.FromErr(err)
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, diag.FromErr(err)
	}
	return config, nil
}

func initWithUsername(d *schema.ResourceData, config *config) error {
	config.UserDomainName = defaultUsersDomainName

	config.Username = os.Getenv("OS_USERNAME")
	if v, ok := d.GetOk("username"); ok {
		config.Username = v.(string)
	}
	if config.Username == "" {
		return fmt.Errorf("username must be specified")
	}
	return nil
}

// Provider returns a schema.Provider for VKCS.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", defaultIdentityEndpoint),
				Description: "The Identity authentication URL.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PROJECT_ID", ""),
				Description: "The ID of Project to login with.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PASSWORD", ""),
				Description: "Password to login with.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_NAME", ""),
				Description: "User name to login with.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_REGION", "RegionOne"),
				Description: "A region to use.",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INSECURE", nil),
				Description: "Trust self-signed certificates.",
			},
			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CACERT", ""),
				Description: "A Custom CA certificate.",
			},
			"cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CERT", ""),
				Description: "A client certificate to authenticate with.",
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KEY", ""),
				Description: "A client private key to authenticate with.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vkcs_compute_keypair":               dataSourceComputeKeypair(),
			"vkcs_compute_instance":              dataSourceComputeInstance(),
			"vkcs_compute_availability_zones":    dataSourceComputeAvailabilityZones(),
			"vkcs_compute_flavor":                dataSourceComputeFlavor(),
			"vkcs_compute_quotaset":              dataSourceComputeQuotaset(),
			"vkcs_images_image":                  dataSourceImagesImage(),
			"vkcs_networking_network":            dataSourceNetworkingNetwork(),
			"vkcs_networking_subnet":             dataSourceNetworkingSubnet(),
			"vkcs_networking_router":             dataSourceNetworkingRouter(),
			"vkcs_networking_port":               dataSourceNetworkingPort(),
			"vkcs_networking_secgroup":           dataSourceNetworkingSecGroup(),
			"vkcs_networking_floatingip":         dataSourceNetworkingFloatingIP(),
			"vkcs_keymanager_secret":             dataSourceKeyManagerSecret(),
			"vkcs_keymanager_container":          dataSourceKeyManagerContainer(),
			"vkcs_blockstorage_volume":           dataSourceBlockStorageVolume(),
			"vkcs_blockstorage_snapshot":         dataSourceBlockStorageSnapshot(),
			"vkcs_sharedfilesystem_sharenetwork": dataSourceSharedFilesystemShareNetwork(),
			"vkcs_sharedfilesystem_share":        dataSourceSharedFilesystemShare(),
			"vkcs_db_database":                   dataSourceDatabaseDatabase(),
			"vkcs_db_instance":                   dataSourceDatabaseInstance(),
			"vkcs_db_user":                       dataSourceDatabaseUser(),
			"vkcs_kubernetes_clustertemplate":    dataSourceKubernetesClusterTemplate(),
			"vkcs_kubernetes_clustertemplates":   dataSourceKubernetesClusterTemplates(),
			"vkcs_kubernetes_cluster":            dataSourceKubernetesCluster(),
			"vkcs_kubernetes_node_group":         dataSourceKubernetesNodeGroup(),
			"vkcs_region":                        dataSourceVkcsRegion(),
			"vkcs_regions":                       dataSourceVkcsRegions(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"vkcs_compute_instance":                   resourceComputeInstance(),
			"vkcs_compute_interface_attach":           resourceComputeInterfaceAttach(),
			"vkcs_compute_keypair":                    resourceComputeKeypair(),
			"vkcs_compute_volume_attach":              resourceComputeVolumeAttach(),
			"vkcs_compute_floatingip_associate":       resourceComputeFloatingIPAssociate(),
			"vkcs_compute_servergroup":                resourceComputeServerGroup(),
			"vkcs_images_image":                       resourceImagesImage(),
			"vkcs_networking_network":                 resourceNetworkingNetwork(),
			"vkcs_networking_subnet":                  resourceNetworkingSubnet(),
			"vkcs_networking_subnet_route":            resourceNetworkingSubnetRoute(),
			"vkcs_networking_router":                  resourceNetworkingRouter(),
			"vkcs_networking_router_interface":        resourceNetworkingRouterInterface(),
			"vkcs_networking_router_route":            resourceNetworkingRouterRoute(),
			"vkcs_networking_port":                    resourceNetworkingPort(),
			"vkcs_networking_port_secgroup_associate": resourceNetworkingPortSecGroupAssociate(),
			"vkcs_networking_secgroup":                resourceNetworkingSecGroup(),
			"vkcs_networking_secgroup_rule":           resourceNetworkingSecGroupRule(),
			"vkcs_networking_floatingip":              resourceNetworkingFloating(),
			"vkcs_networking_floatingip_associate":    resourceNetworkingFloatingIPAssociate(),
			"vkcs_keymanager_secret":                  resourceKeyManagerSecret(),
			"vkcs_keymanager_container":               resourceKeyManagerContainer(),
			"vkcs_blockstorage_volume":                resourceBlockStorageVolume(),
			"vkcs_blockstorage_snapshot":              resourceBlockStorageSnapshot(),
			"vkcs_lb_l7policy":                        resourceL7Policy(),
			"vkcs_lb_l7rule":                          resourceL7Rule(),
			"vkcs_lb_listener":                        resourceListener(),
			"vkcs_lb_loadbalancer":                    resourceLoadBalancer(),
			"vkcs_lb_member":                          resourceMember(),
			"vkcs_lb_members":                         resourceMembers(),
			"vkcs_lb_monitor":                         resourceMonitor(),
			"vkcs_lb_pool":                            resourcePool(),
			"vkcs_vpnaas_endpoint_group":              resourceEndpointGroup(),
			"vkcs_vpnaas_ike_policy":                  resourceIKEPolicy(),
			"vkcs_vpnaas_ipsec_policy":                resourceIPSecPolicy(),
			"vkcs_vpnaas_service":                     resourceService(),
			"vkcs_vpnaas_site_connection":             resourceSiteConnection(),
			"vkcs_sharedfilesystem_securityservice":   resourceSharedFilesystemSecurityService(),
			"vkcs_sharedfilesystem_sharenetwork":      resourceSharedFilesystemShareNetwork(),
			"vkcs_sharedfilesystem_share":             resourceSharedFilesystemShare(),
			"vkcs_sharedfilesystem_share_access":      resourceSharedFilesystemShareAccess(),
			"vkcs_db_instance":                        resourceDatabaseInstance(),
			"vkcs_db_database":                        resourceDatabaseDatabase(),
			"vkcs_db_user":                            resourceDatabaseUser(),
			"vkcs_db_cluster":                         resourceDatabaseCluster(),
			"vkcs_db_cluster_with_shards":             resourceDatabaseClusterWithShards(),
			"vkcs_kubernetes_cluster":                 resourceKubernetesCluster(),
			"vkcs_kubernetes_node_group":              resourceKubernetesNodeGroup(),
		},
	}

	provider.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return newConfig(d, terraformVersion)
	}

	return provider
}
