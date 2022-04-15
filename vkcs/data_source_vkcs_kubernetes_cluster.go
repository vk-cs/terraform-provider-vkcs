package vkcs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesClusterRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"discovery_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_flavor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"keypair": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"master_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"master_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"stack_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pods_network_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"floating_ip_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"api_lb_vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_lb_fip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ingress_floating_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registry_auth_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"k8s_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"loadbalancer_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"insecure_registries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceKubernetesClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	containerInfraClient, err := config.ContainerInfraV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra client: %s", err)
	}
	clusterIdentifierName, err := ensureOnlyOnePresented(d, "cluster_id", "name")
	if err != nil {
		return diag.FromErr(err)
	}
	clusterIdentifier := d.Get(clusterIdentifierName).(string)
	c, err := clusterGet(containerInfraClient, clusterIdentifier).Extract()
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
	d.Set("availability_zone", c.AvailabilityZone)
	d.Set("loadbalancer_subnet_id", c.LoadbalancerSubnetID)
	d.Set("insecure_registries", c.InsecureRegistries)

	k8sConfig, err := k8sConfigGet(containerInfraClient, c.UUID)
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

	d.Set("region", getRegion(d, config))

	return nil
}
