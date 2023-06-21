package provider

import (
	"context"

	sdkdiag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/blockstorage"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/compute"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/db"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/firewall"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/images"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/keymanager"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/lb"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/publicdns"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/regions"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/sharedfilesystem"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/vpnaas"

	wrapper "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/providerwrapper/sdk"
)

// SDKProvider returns a SDKv2 schema.Provider for VKCS.
func SDKProvider() *sdkschema.Provider {
	provider, err := wrapper.WrapProvider(SDKProviderBase())
	if err != nil {
		panic(err)
	}

	return provider
}

func SDKProviderBase() *sdkschema.Provider {
	provider := &sdkschema.Provider{
		Schema: map[string]*sdkschema.Schema{
			"auth_url": {
				Type:        sdkschema.TypeString,
				Optional:    true,
				Description: "The Identity authentication URL.",
			},
			"project_id": {
				Type:        sdkschema.TypeString,
				Optional:    true,
				Description: "The ID of Project to login with.",
			},
			"password": {
				Type:        sdkschema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password to login with.",
			},
			"username": {
				Type:        sdkschema.TypeString,
				Optional:    true,
				Description: "User name to login with.",
			},
			"user_domain_id": {
				Type:        sdkschema.TypeString,
				Optional:    true,
				Description: "The id of the domain where the user resides.",
			},
			"user_domain_name": {
				Type:        sdkschema.TypeString,
				Optional:    true,
				Description: "The name of the domain where the user resides.",
			},
			"region": {
				Type:        sdkschema.TypeString,
				Optional:    true,
				Description: "A region to use.",
			},
			"cloud_containers_api_version": {
				Type:        sdkschema.TypeString,
				Optional:    true,
				Description: "Cloud Containers API version to use. _note_ Only for custom VKCS deployments.",
			},
		},

		DataSourcesMap: map[string]*sdkschema.Resource{
			"vkcs_compute_keypair":               compute.DataSourceComputeKeypair(),
			"vkcs_compute_instance":              compute.DataSourceComputeInstance(),
			"vkcs_compute_availability_zones":    compute.DataSourceComputeAvailabilityZones(),
			"vkcs_compute_flavor":                compute.DataSourceComputeFlavor(),
			"vkcs_compute_quotaset":              compute.DataSourceComputeQuotaset(),
			"vkcs_images_image":                  images.DataSourceImagesImage(),
			"vkcs_networking_network":            networking.DataSourceNetworkingNetwork(),
			"vkcs_networking_subnet":             networking.DataSourceNetworkingSubnet(),
			"vkcs_networking_router":             networking.DataSourceNetworkingRouter(),
			"vkcs_networking_port":               networking.DataSourceNetworkingPort(),
			"vkcs_networking_secgroup":           firewall.DataSourceNetworkingSecGroup(),
			"vkcs_networking_floatingip":         networking.DataSourceNetworkingFloatingIP(),
			"vkcs_keymanager_container":          keymanager.DataSourceKeyManagerContainer(),
			"vkcs_blockstorage_volume":           blockstorage.DataSourceBlockStorageVolume(),
			"vkcs_blockstorage_snapshot":         blockstorage.DataSourceBlockStorageSnapshot(),
			"vkcs_lb_loadbalancer":               lb.DataSourceLoadBalancer(),
			"vkcs_sharedfilesystem_sharenetwork": sharedfilesystem.DataSourceSharedFilesystemShareNetwork(),
			"vkcs_sharedfilesystem_share":        sharedfilesystem.DataSourceSharedFilesystemShare(),
			"vkcs_db_database":                   db.DataSourceDatabaseDatabase(),
			"vkcs_db_instance":                   db.DataSourceDatabaseInstance(),
			"vkcs_db_user":                       db.DataSourceDatabaseUser(),
			"vkcs_kubernetes_clustertemplate":    kubernetes.DataSourceKubernetesClusterTemplate(),
			"vkcs_kubernetes_cluster":            kubernetes.DataSourceKubernetesCluster(),
			"vkcs_kubernetes_node_group":         kubernetes.DataSourceKubernetesNodeGroup(),
			"vkcs_region":                        regions.DataSourceVkcsRegion(),
			"vkcs_regions":                       regions.DataSourceVkcsRegions(),
			"vkcs_publicdns_zone":                publicdns.DataSourcePublicDNSZone(),
		},

		ResourcesMap: map[string]*sdkschema.Resource{
			"vkcs_compute_instance":                   compute.ResourceComputeInstance(),
			"vkcs_compute_interface_attach":           compute.ResourceComputeInterfaceAttach(),
			"vkcs_compute_keypair":                    compute.ResourceComputeKeypair(),
			"vkcs_compute_volume_attach":              compute.ResourceComputeVolumeAttach(),
			"vkcs_compute_floatingip_associate":       compute.ResourceComputeFloatingIPAssociate(),
			"vkcs_compute_servergroup":                compute.ResourceComputeServerGroup(),
			"vkcs_images_image":                       images.ResourceImagesImage(),
			"vkcs_networking_network":                 networking.ResourceNetworkingNetwork(),
			"vkcs_networking_subnet":                  networking.ResourceNetworkingSubnet(),
			"vkcs_networking_subnet_route":            networking.ResourceNetworkingSubnetRoute(),
			"vkcs_networking_router":                  networking.ResourceNetworkingRouter(),
			"vkcs_networking_router_interface":        networking.ResourceNetworkingRouterInterface(),
			"vkcs_networking_router_route":            networking.ResourceNetworkingRouterRoute(),
			"vkcs_networking_port":                    networking.ResourceNetworkingPort(),
			"vkcs_networking_port_secgroup_associate": networking.ResourceNetworkingPortSecGroupAssociate(),
			"vkcs_networking_secgroup":                firewall.ResourceNetworkingSecGroup(),
			"vkcs_networking_secgroup_rule":           firewall.ResourceNetworkingSecGroupRule(),
			"vkcs_networking_floatingip":              networking.ResourceNetworkingFloating(),
			"vkcs_networking_floatingip_associate":    networking.ResourceNetworkingFloatingIPAssociate(),
			"vkcs_keymanager_secret":                  keymanager.ResourceKeyManagerSecret(),
			"vkcs_keymanager_container":               keymanager.ResourceKeyManagerContainer(),
			"vkcs_blockstorage_volume":                blockstorage.ResourceBlockStorageVolume(),
			"vkcs_blockstorage_snapshot":              blockstorage.ResourceBlockStorageSnapshot(),
			"vkcs_lb_l7policy":                        lb.ResourceL7Policy(),
			"vkcs_lb_l7rule":                          lb.ResourceL7Rule(),
			"vkcs_lb_listener":                        lb.ResourceListener(),
			"vkcs_lb_loadbalancer":                    lb.ResourceLoadBalancer(),
			"vkcs_lb_member":                          lb.ResourceMember(),
			"vkcs_lb_members":                         lb.ResourceMembers(),
			"vkcs_lb_monitor":                         lb.ResourceMonitor(),
			"vkcs_lb_pool":                            lb.ResourcePool(),
			"vkcs_vpnaas_endpoint_group":              vpnaas.ResourceEndpointGroup(),
			"vkcs_vpnaas_ike_policy":                  vpnaas.ResourceIKEPolicy(),
			"vkcs_vpnaas_ipsec_policy":                vpnaas.ResourceIPSecPolicy(),
			"vkcs_vpnaas_service":                     vpnaas.ResourceService(),
			"vkcs_vpnaas_site_connection":             vpnaas.ResourceSiteConnection(),
			"vkcs_sharedfilesystem_securityservice":   sharedfilesystem.ResourceSharedFilesystemSecurityService(),
			"vkcs_sharedfilesystem_sharenetwork":      sharedfilesystem.ResourceSharedFilesystemShareNetwork(),
			"vkcs_sharedfilesystem_share":             sharedfilesystem.ResourceSharedFilesystemShare(),
			"vkcs_sharedfilesystem_share_access":      sharedfilesystem.ResourceSharedFilesystemShareAccess(),
			"vkcs_db_backup":                          db.ResourceDatabaseBackup(),
			"vkcs_db_instance":                        db.ResourceDatabaseInstance(),
			"vkcs_db_database":                        db.ResourceDatabaseDatabase(),
			"vkcs_db_user":                            db.ResourceDatabaseUser(),
			"vkcs_db_cluster":                         db.ResourceDatabaseCluster(),
			"vkcs_db_cluster_with_shards":             db.ResourceDatabaseClusterWithShards(),
			"vkcs_db_config_group":                    db.ResourceDatabaseConfigGroup(),
			"vkcs_kubernetes_cluster":                 kubernetes.ResourceKubernetesCluster(),
			"vkcs_kubernetes_node_group":              kubernetes.ResourceKubernetesNodeGroup(),
			"vkcs_publicdns_zone":                     publicdns.ResourcePublicDNSZone(),
			"vkcs_publicdns_record":                   publicdns.ResourcePublicDNSRecord(),
		},
	}

	provider.ConfigureContextFunc = func(_ context.Context, d *sdkschema.ResourceData) (interface{}, sdkdiag.Diagnostics) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return clients.ConfigureSdkProvider(d, terraformVersion)
	}

	return provider
}
