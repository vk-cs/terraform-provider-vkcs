package vkcs

import (
	"context"
	"fmt"
	"os"

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
)

// configer is interface to work with gophercloud.Config calls
type configer interface {
	LoadAndValidate() error
	GetRegion() string
	ComputeV2Client(region string) (*gophercloud.ServiceClient, error)
	ImageV2Client(region string) (*gophercloud.ServiceClient, error)
	NetworkingV2Client(region string) (*gophercloud.ServiceClient, error)
	BlockStorageV3Client(region string) (*gophercloud.ServiceClient, error)
	KeyManagerV1Client(region string) (*gophercloud.ServiceClient, error)
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

func (c *config) NetworkingV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.Config.NetworkingV2Client(region)
}

func (c *config) BlockStorageV3Client(region string) (*gophercloud.ServiceClient, error) {
	return c.Config.BlockStorageV3Client(region)
}

func (c *config) KeyManagerV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.Config.KeyManagerV1Client(region)
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
			AllowReauth:      true,
			MaxRetries:       maxRetriesCount,
			TerraformVersion: terraformVersion,
			SDKVersion:       meta.SDKVersionString(),
			MutexKV:          mutexkv.NewMutexKV(),
		},
	}

	if config.TenantID == "" {
		config.TenantID = os.Getenv("OS_PROJECT_ID")
	}
	if config.UserDomainID == "" {
		config.UserDomainID = os.Getenv("OS_USER_DOMAIN_ID")
	}
	if config.Password == "" {
		config.Password = os.Getenv("OS_PASSWORD")
	}
	if config.Username == "" {
		config.Username = os.Getenv("OS_USERNAME")
	}
	if config.Region == "" {
		config.Region = os.Getenv("OS_REGION")
	}

	v, ok := d.GetOk("insecure")
	if ok {
		insecure := v.(bool)
		config.Insecure = &insecure
	}
	v, ok = d.GetOk("auth_url")
	if ok {
		config.IdentityEndpoint = v.(string)
	} else {
		config.IdentityEndpoint = defaultIdentityEndpoint
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
				DefaultFunc: schema.EnvDefaultFunc("AUTH_URL", ""),
				Description: "The Identity authentication URL.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROJECT_ID", ""),
				Description: "The ID of Project to login with.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("PASSWORD", ""),
				Description: "Password to login with.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("USER_NAME", ""),
				Description: "User name to login with.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("REGION", "RegionOne"),
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
