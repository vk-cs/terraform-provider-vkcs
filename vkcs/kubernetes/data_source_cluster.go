package kubernetes

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func DataSourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesClusterRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the Container Infra client. If omitted, the `region` argument of the provider is used.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the cluster. _note_ Only one of `name` or `cluster_id` must be specified.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The UUID of the Kubernetes cluster template. _note_ Only one of `name` or `cluster_id` must be specified.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The project of the cluster.",
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user of the cluster.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time at which cluster was created.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time at which cluster was created.",
			},
			"api_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "COE API address.",
			},
			"cluster_template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the V1 Container Infra cluster template.",
			},
			"discovery_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL used for cluster node discovery.",
			},
			"master_flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the flavor for the master nodes.",
			},
			"keypair": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Compute service SSH keypair.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The list of key value pairs representing additional properties of the cluster.",
			},
			"master_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of master nodes for the cluster.",
			},
			"master_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "IP addresses of the master node of the cluster.",
			},
			"stack_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the Orchestration service stack.",
			},
			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the cluster's network.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the cluster's subnet.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current state of a cluster.",
			},
			"pods_network_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network cidr of k8s virtual network.",
			},
			"floating_ip_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether floating ip is enabled for cluster.",
			},
			"api_lb_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API LoadBalancer vip.",
			},
			"api_lb_fip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API LoadBalancer fip.",
			},
			"ingress_floating_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "This argument is deprecated as Ingress controller is not currently installed by default.",
				Description: "Floating IP created for ingress service.",
			},
			"registry_auth_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Docker registry access password.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Availability zone of the cluster.",
			},
			"availability_zones": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Availability zones of the regional cluster",
			},
			"k8s_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kubeconfig for cluster",
			},
			"loadbalancer_subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of load balancer's subnet.",
			},
			"insecure_registries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Addresses of registries from which you can download images without checking certificates.",
			},
			"dns_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Custom DNS cluster domain.",
			},
			"sync_security_policy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enables syncing of security policies of cluster.",
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the kubernetes cluster, may be `standard` or `regional`",
			},
		},
		Description: "Use this data source to get the ID of an available VKCS kubernetes cluster.",
	}
}

func dataSourceKubernetesClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	containerInfraClient, err := config.ContainerInfraV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra client: %s", err)
	}
	clusterIdentifierName, err := util.EnsureOnlyOnePresented(d, "cluster_id", "name")
	if err != nil {
		return diag.FromErr(err)
	}
	clusterIdentifier := d.Get(clusterIdentifierName).(string)
	c, err := clusters.Get(containerInfraClient, clusterIdentifier).Extract()
	if err != nil {
		return diag.Errorf("error getting vkcs_kubernetes_cluster %s: %s", clusterIdentifier, err)
	}

	d.SetId(c.UUID)
	d.Set("name", c.Name)
	d.Set("project_id", c.ProjectID)
	d.Set("user_id", c.UserID)
	d.Set("api_address", c.APIAddress)
	d.Set("cluster_template_id", c.ClusterTemplateID)
	d.Set("discovery_url", c.DiscoveryURL)
	d.Set("master_flavor", c.MasterFlavorID)
	d.Set("keypair", c.KeyPair)
	d.Set("master_count", c.MasterCount)
	d.Set("master_addresses", c.MasterAddresses)
	d.Set("stack_id", c.StackID)
	d.Set("network_id", c.NetworkID)
	d.Set("subnet_id", c.SubnetID)
	d.Set("status", c.NewStatus)
	d.Set("pods_network_cidr", c.PodsNetworkCidr)
	d.Set("floating_ip_enabled", c.FloatingIPEnabled)
	d.Set("api_lb_vip", c.APILBVIP)
	d.Set("api_lb_fip", c.APILBFIP)
	d.Set("ingress_floating_ip", c.IngressFloatingIP)
	d.Set("registry_auth_password", c.RegistryAuthPassword)
	d.Set("loadbalancer_subnet_id", c.LoadbalancerSubnetID)
	d.Set("insecure_registries", c.InsecureRegistries)
	d.Set("dns_domain", c.DNSDomain)
	d.Set("sync_security_policy", c.SecurityPolicySyncEnabled)
	d.Set("cluster_type", c.ClusterType)
	d.Set("availability_zone", c.AvailabilityZone)
	d.Set("availability_zones", c.AvailabilityZones)

	k8sConfig, err := clusters.KubeConfigGet(containerInfraClient, c.UUID)
	if err != nil {
		log.Printf("[DEBUG] error getting k8s config for cluster %s: %s", c.UUID, err)
		d.Set("k8s_config", "error")
	} else {
		d.Set("k8s_config", k8sConfig)
	}

	if err := d.Set("labels", c.Labels); err != nil {
		log.Printf("[DEBUG] unable to set labels for vkcs_kubernetes_cluster %s: %s", c.UUID, err)
	}
	if err := d.Set("created_at", c.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] unable to set created_at for vkcs_kubernetes_cluster %s: %s", c.UUID, err)
	}
	if err := d.Set("updated_at", c.UpdatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] unable to set updated_at for vkcs_kubernetes_cluster %s: %s", c.UUID, err)
	}

	d.Set("region", util.GetRegion(d, config))

	return nil
}
