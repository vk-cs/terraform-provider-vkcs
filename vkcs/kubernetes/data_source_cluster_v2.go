package kubernetes

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func DataSourceKubernetesClusterV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesClusterV2Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the Container Infra client. If omitted, the `region` argument of the provider is used.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the cluster.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the cluster.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kubernetes version of the cluster.",
			},
			"enable_public_ip": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable public IP for the cluster.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the cluster.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the cluster.",
			},
			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the cluster's network.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the cluster's subnet.",
			},
			"availability_zones": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Availability zones of the cluster.",
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the kubernetes cluster.",
			},
			"external_network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of external network.",
			},
			"insecure_registries": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Addresses of registries from which you can download images without checking certificates.",
			},
			"network_plugin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network plugin.",
			},
			"pods_ipv4_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IPv4 CIDR for pods network.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Key-value pairs of labels for the cluster.",
			},
			"loadbalancer_subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the load balancer's subnet.",
			},
			"loadbalancer_allowed_cidrs": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of CIDR blocks allowed to access the load balancer.",
			},
			"master_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of master nodes.",
			},
			"master_flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Flavor UUID for master nodes.",
			},
			"master_disks": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of master disk.",
						},
						"size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Size of master disk.",
						},
					},
				},
				Description: "List of master disks.",
			},
			"k8s_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Contents of the kubeconfig file.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time at which cluster was created.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster current status.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User project ID.",
			},
			"api_lb_fip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API LoadBalancer FIP (Floating IP).",
			},
			"api_lb_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API LoadBalancer VIP (Virtual IP).",
			},
			"api_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL address of cluster kubeapi-server.",
			},
			"node_groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of cluster node group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of cluster node group.",
						},
						"flavor": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Flavor of nodes in node group.",
						},
						"node_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of nodes in node group.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone of node group",
						},
					},
				},
				Description: "The list of cluster node groups.",
			},
		},
		Description: "Use this data source to get the ID of an available VKCS kubernetes cluster v2.",
	}
}

func dataSourceKubernetesClusterV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// creates client for managed-k8s
	config := meta.(clients.Config)
	containerInfraClientV2, err := config.ContainerInfraV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra v2 client: %s", err)
	}

	clusterID := d.Get("id").(string)
	if clusterID == "" {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "id must be specified",
			AttributePath: cty.Path{
				cty.GetAttrStep{Name: "id"},
			},
		}}
	}
	d.SetId(clusterID)

	cluster, err := clusters.Get(containerInfraClientV2, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "error retrieving vkcs_kubernetes_cluster_v2"))
	}

	log.Printf("[DEBUG] retrieved vkcs_kubernetes_cluster_v2 %s: %#v", d.Id(), cluster)

	// set basic fields
	d.Set("uuid", cluster.UUID)
	d.Set("name", cluster.Name)
	d.Set("version", cluster.Version)
	d.Set("status", cluster.Status)
	d.Set("labels", cluster.Labels)
	d.Set("description", cluster.Description)
	d.Set("insecure_registries", cluster.InsecureRegistries)
	d.Set("region", util.GetRegion(d, config))
	d.Set("created_at", cluster.CreatedAt)
	d.Set("project_id", cluster.ProjectID)
	d.Set("api_lb_fip", cluster.ExternalIP)
	d.Set("api_lb_vip", cluster.InternalIP)
	d.Set("api_address", cluster.ApiAddress)

	// set master spec fields
	d.Set("master_flavor", cluster.MasterSpec.Engine.NovaEngine.FlavorID)
	d.Set("master_count", cluster.MasterSpec.Replicas)
	masterDisks := make([]map[string]interface{}, 0, len(cluster.MasterSpec.Disks))
	for _, disk := range cluster.MasterSpec.Disks {
		masterDisks = append(masterDisks, map[string]interface{}{
			"type": disk.Type,
			"size": strconv.Itoa(disk.Size),
		})
	}
	d.Set("master_disks", masterDisks)

	// set deployment type fields
	if cluster.DeploymentType.ZonalDeployment != nil {
		d.Set("cluster_type", "standard")
		d.Set("availability_zones", []string{cluster.DeploymentType.ZonalDeployment.Zone})
	} else if cluster.DeploymentType.MultiZonalDeployment != nil {
		d.Set("cluster_type", "regional")
		// Преобразуем []string в []interface{}
		zones := make([]interface{}, len(cluster.DeploymentType.MultiZonalDeployment.Zones))
		for i, zone := range cluster.DeploymentType.MultiZonalDeployment.Zones {
			zones[i] = zone
		}
		d.Set("availability_zones", zones)
	}

	// set network config fields
	if cluster.NetworkConfig.Plugin.Calico != nil {
		d.Set("network_plugin", "calico")
		d.Set("pods_ipv4_cidr", cluster.NetworkConfig.Plugin.Calico.PodsIPv4CIDR)
	} else if cluster.NetworkConfig.Plugin.Cilium != nil {
		d.Set("network_plugin", "cilium")
		d.Set("pods_ipv4_cidr", cluster.NetworkConfig.Plugin.Cilium.PodsIPv4CIDR)
	}

	// set network engine fields
	d.Set("network_id", cluster.NetworkConfig.Engine.SprutEngine.NetworkID)
	d.Set("subnet_id", cluster.NetworkConfig.Engine.SprutEngine.SubnetID)
	d.Set("external_network_id", cluster.NetworkConfig.Engine.SprutEngine.ExternalNetworkID)

	// set load balancer config fields
	d.Set("enable_public_ip", cluster.LoadBalancerConfig.OctaviaEngine.EnablePublicIP)
	d.Set("loadbalancer_subnet_id", cluster.LoadBalancerConfig.OctaviaEngine.LoadbalancerSubnetID)
	d.Set("loadbalancer_allowed_cidrs", cluster.LoadBalancerConfig.OctaviaEngine.AllowedCIDRs)

	// set node_groups
	clusterNodeGroups := make([]map[string]any, len(cluster.NodeGroups))
	for i, ng := range cluster.NodeGroups {
		ngZone := ""
		if len(ng.Zones) > 0 {
			ngZone = ng.Zones[0]
		}
		clusterNodeGroups[i] = map[string]any{
			"id":                ng.ID,
			"name":              ng.Name,
			"flavor":            ng.VMEngine.NovaEngine.FlavorID,
			"node_count":        ng.GetActualSize(),
			"availability_zone": ngZone,
		}
	}
	d.Set("node_groups", clusterNodeGroups)

	// get cluster kubeconfig
	kubeconfig, err := clusters.GetKubeconfig(containerInfraClientV2, d.Id())
	if err != nil {
		return diag.Errorf("error retrieving vkcs_kubernetes_cluster_v2 %s kubeconfig: %s", d.Id(), err)
	}
	d.Set("k8s_config", kubeconfig)

	return nil
}
