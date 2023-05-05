package vkcs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
)

// Provider returns a schema.Provider for VKCS.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", clients.DefaultIdentityEndpoint),
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
				DefaultFunc: schema.EnvDefaultFunc("OS_USERNAME", ""),
				Description: "User name to login with.",
			},
			"user_domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_DOMAIN_ID", ""),
				Description: "The id of the domain where the user resides.",
			},
			"user_domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_DOMAIN_NAME", clients.DefaultUserDomainName),
				Description: "The name of the domain where the user resides.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_REGION_NAME", clients.DefaultRegionName),
				Description: "A region to use.",
			},
			"cloud_containers_api_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  clients.ContainerInfraAPIVersion,
				Description: "Cloud Containers API version to use.\n" +
					"_NOTE_ Only for custom VKCS deployments.",
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
			"vkcs_lb_loadbalancer":               dataSourceLoadBalancer(),
			"vkcs_sharedfilesystem_sharenetwork": dataSourceSharedFilesystemShareNetwork(),
			"vkcs_sharedfilesystem_share":        dataSourceSharedFilesystemShare(),
			"vkcs_db_database":                   dataSourceDatabaseDatabase(),
			"vkcs_db_datastore":                  dataSourceDatabaseDatastore(),
			"vkcs_db_datastore_capabilities":     dataSourceDatabaseDatastoreCapabilities(),
			"vkcs_db_datastore_parameters":       dataSourceDatabaseDatastoreParameters(),
			"vkcs_db_datastores":                 dataSourceDatabaseDatastores(),
			"vkcs_db_instance":                   dataSourceDatabaseInstance(),
			"vkcs_db_user":                       dataSourceDatabaseUser(),
			"vkcs_db_backup":                     dataSourceDatabaseBackup(),
			"vkcs_db_config_group":               dataSourceDatabaseConfigGroup(),
			"vkcs_kubernetes_clustertemplate":    dataSourceKubernetesClusterTemplate(),
			"vkcs_kubernetes_clustertemplates":   dataSourceKubernetesClusterTemplates(),
			"vkcs_kubernetes_cluster":            dataSourceKubernetesCluster(),
			"vkcs_kubernetes_node_group":         dataSourceKubernetesNodeGroup(),
			"vkcs_region":                        dataSourceVkcsRegion(),
			"vkcs_regions":                       dataSourceVkcsRegions(),
			"vkcs_publicdns_zone":                dataSourcePublicDNSZone(),
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
			"vkcs_db_backup":                          resourceDatabaseBackup(),
			"vkcs_db_instance":                        resourceDatabaseInstance(),
			"vkcs_db_database":                        resourceDatabaseDatabase(),
			"vkcs_db_user":                            resourceDatabaseUser(),
			"vkcs_db_cluster":                         resourceDatabaseCluster(),
			"vkcs_db_cluster_with_shards":             resourceDatabaseClusterWithShards(),
			"vkcs_db_config_group":                    resourceDatabaseConfigGroup(),
			"vkcs_kubernetes_cluster":                 resourceKubernetesCluster(),
			"vkcs_kubernetes_node_group":              resourceKubernetesNodeGroup(),
			"vkcs_publicdns_zone":                     resourcePublicDNSZone(),
			"vkcs_publicdns_record":                   resourcePublicDNSRecord(),
		},
	}

	provider.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return clients.ConfigureProvider(d, terraformVersion)
	}

	return provider
}
