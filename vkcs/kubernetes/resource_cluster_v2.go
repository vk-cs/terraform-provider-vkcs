package kubernetes

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

const (
	operationCreateV2 = 60
	operationUpdateV2 = 60
	operationDeleteV2 = 30
)

const (
	createUpdateDelayV2        = 5
	createUpdatePollIntervalV2 = 30
	deleteDelayV2              = 5
	deletePollIntervalV2       = 30
)

// cluster statuses from new API
const (
	clusterStatusV2Provisioning = "PROVISIONING"
	clusterStatusV2Starting     = "STARTING"
	clusterStatusV2Running      = "RUNNING"
	clusterStatusV2Reconciling  = "RECONCILING"
	clusterStatusV2Deleting     = "DELETING"
	clusterStatusV2Failed       = "FAILED"
	clusterStatusV2Deleted      = "DELETED"
)

func ResourceKubernetesClusterV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesClusterV2Create,
		ReadContext:   resourceKubernetesClusterV2Read,
		UpdateContext: resourceKubernetesClusterV2Update,
		DeleteContext: resourceKubernetesClusterV2Delete,

		CustomizeDiff: customizeClusterV2Diff,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(operationCreateV2 * time.Minute),
			Update: schema.DefaultTimeout(operationUpdateV2 * time.Minute),
			Delete: schema.DefaultTimeout(operationDeleteV2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					isUUID,
				)),
				Description: "UUID of the cluster. If not provided, will be generated automatically.",
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false, // при изменении будет запускаться обновление, а не пересоздание кластера
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringIsNotEmpty,
				)),
				// TODO: Добавить валидацию, которая делает запрос в mk8s-api и получает список доступных версий Kubernetes
				Description: "Kubernetes version of the cluster. Changing this upgrades a cluster.",
			},
			"enable_public_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Enable public IP for the cluster. If true, a floating IP will be assigned to the cluster. Default false.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					isClusterNameV2,
				)),
				Description: "The name of the cluster. Should match the pattern `^[a-zA-Z][a-zA-Z0-9_.-]*$`.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    false,
				Description: "Description of the cluster.",
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					isUUID,
				)),
				Description: "ID of the cluster's network. Cluster can be created only in network with SDN=sprut.",
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					isUUID,
				)),
				Description: "ID of the cluster's subnet.",
			},
			"availability_zones": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				// TODO: Добавить валидацию, которая делает запрос в mk8s-api и получает список зон доступности
				Description: "Availability zones of the cluster. Use 1 zone for standard cluster, 3 zones for regional cluster.",
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"standard",
					"regional",
				}, false),
				Description: "Type of the kubernetes cluster. Must be either 'standard' or 'regional'.",
			},
			"external_network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					isUUID,
				)),
				Description: "ID of external network.",
			},
			"insecure_registries": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: false,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateInsecureRegistryURLV2,
				},
				Description: "Addresses of registries from which you can download images without checking certificates.",
			},
			"network_plugin": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"calico",
					// "cilium", // пока не работает
				}, false),
				Description: "Network plugin. Must be 'calico'.",
			},
			"pods_ipv4_cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "IPv4 CIDR for pods network.",
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: false,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// TODO: Добавить валидацию для labels keys values
				Description: "Key-value pairs of labels for the cluster.",
			},
			"loadbalancer_subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					isUUID,
				)),
				Description: "The UUID of the load balancer's subnet.",
			},
			"loadbalancer_allowed_cidrs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: false,
				ForceNew: false,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Description: "List of CIDR blocks allowed to access the load balancer.",
			},
			"master_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "Number of master nodes. Use 1 for standard cluster, 3 or 5 for regional cluster. Changing this creates a new cluster.",
			},
			"master_flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false, // при изменении будет запускаться вертикальное скалирование, а не пересоздание кластера
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					isUUID,
				)),
				Description: "Flavor UUID for master nodes. Changing this scales master nodes to the target flavor.",
			},
			"master_disks": {
				Type:     schema.TypeSet,
				Computed: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							ForceNew:    false,
							Description: "Type of master disk.",
						},
						"size": {
							Type:        schema.TypeString,
							Computed:    true,
							ForceNew:    false,
							Description: "Size of master disk.",
						},
					},
				},
				Description: "List of master disks.",
			},
			"k8s_config": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    false,
				Description: "Contents of the kubeconfig file. Use it to authenticate to Kubernetes cluster.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    false,
				Description: "The time at which cluster was created.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    false,
				Description: "Cluster current status.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    false,
				Description: "User project ID.",
			},
			"api_lb_fip": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    false,
				Description: "API LoadBalancer FIP (Floating IP). IP address field.",
			},
			"api_lb_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    false,
				Description: "API LoadBalancer VIP (Virtual IP). IP address field.",
			},
			"api_address": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    false,
				Description: "URL address of cluster kubeapi-server.",
			},
			"node_groups": {
				Type:     schema.TypeSet,
				Computed: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							ForceNew:    false,
							Description: "ID of cluster node group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							ForceNew:    false,
							Description: "Name of cluster node group.",
						},
						"flavor": {
							Type:        schema.TypeString,
							Computed:    true,
							ForceNew:    false,
							Description: "Flavor of nodes in node group.",
						},
						"node_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							ForceNew:    false,
							Description: "Number of nodes in node group.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							ForceNew:    false,
							Description: "Availability zone of node group",
						},
					},
				},
				Description: "The list of cluster node groups.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Region to use for the cluster. Default is a region configured for provider.",
			},
		},
		Description: "Provides a kubernetes cluster resource v2. This can be used to create, modify and delete kubernetes clusters v2.",
	}
}

func resourceKubernetesClusterV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// creates client for managed-k8s
	config := meta.(clients.Config)
	containerInfraClientV2, err := config.ContainerInfraV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra v2 client: %s", err)
	}

	// get and check labels map
	rawLabels := d.Get("labels").(map[string]any)
	labels, err := extractKubernetesLabelsMap(rawLabels)
	if err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
			AttributePath: cty.Path{
				cty.GetAttrStep{Name: "labels"},
			},
		}}
	}

	// build master spec
	masterSpec := clusters.MasterSpecOpts{
		Engine: clusters.MasterEngineOpts{
			NovaEngine: clusters.NovaEngineOpts{
				FlavorID: d.Get("master_flavor").(string),
			},
		},
		Replicas: d.Get("master_count").(int),
	}

	// build deployment type
	clusterType := d.Get("cluster_type").(string)
	var deploymentType clusters.DeploymentTypeOpts

	availabilityZones, err := getAsStringSlice(d, "availability_zones")
	if err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
			AttributePath: cty.Path{
				cty.GetAttrStep{Name: "availability_zones"},
			},
		}}
	}
	switch clusterType {
	case "standard":
		deploymentType = clusters.DeploymentTypeOpts{
			ZonalDeployment: &clusters.ZonalDeploymentOpts{
				Zone: availabilityZones[0],
			},
		}
	case "regional":
		deploymentType = clusters.DeploymentTypeOpts{
			MultiZonalDeployment: &clusters.MultiZonalDeploymentOpts{
				Zones: availabilityZones,
			},
		}
	}

	// build network plugin
	networkPlugin := d.Get("network_plugin").(string)
	var plugin clusters.NetworkPluginOpts

	switch networkPlugin {
	case "calico":
		plugin = clusters.NetworkPluginOpts{
			Calico: &clusters.CalicoPluginOpts{
				PodsIPv4CIDR: d.Get("pods_ipv4_cidr").(string),
			},
		}
	// TODO: разобраться, что с cillium - предположение, что без нод cillium не поднять
	// (пока что этот case никогда не выполнится, потому что есть валидация на поле network_plugin - только calico)
	case "cilium":
		plugin = clusters.NetworkPluginOpts{
			Cilium: &clusters.CiliumPluginOpts{
				PodsIPv4CIDR: d.Get("pods_ipv4_cidr").(string),
			},
		}
	}

	// build network engine
	sprutEngine := clusters.SprutEngineOpts{
		NetworkID: d.Get("network_id").(string),
		SubnetID:  d.Get("subnet_id").(string),
	}

	// add external network ID if specified
	if externalNetworkID, ok := d.GetOk("external_network_id"); ok {
		sprutEngine.ExternalNetworkID = externalNetworkID.(string)
	}

	networkEngine := clusters.NetworkEngineOpts{
		SprutEngine: sprutEngine,
	}

	// build network config
	networkConfig := clusters.NetworkConfigOpts{
		Plugin: plugin,
		Engine: networkEngine,
	}

	// build load balancer config if needed
	var loadBalancerConfig clusters.LoadBalancerConfigOpts
	octaviaEngine := clusters.OctaviaEngineOpts{
		EnablePublicIP:       d.Get("enable_public_ip").(bool),
		LoadbalancerSubnetID: d.Get("loadbalancer_subnet_id").(string),
	}

	if _, hasAllowedCIDRs := d.GetOk("loadbalancer_allowed_cidrs"); hasAllowedCIDRs {
		allowedCidrs, err := getAsStringSlice(d, "loadbalancer_allowed_cidrs")
		if err != nil {
			return diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
				AttributePath: cty.Path{
					cty.GetAttrStep{Name: "loadbalancer_allowed_cidrs"},
				},
			}}
		}
		octaviaEngine.AllowedCIDRs = allowedCidrs
	}

	loadBalancerConfig = clusters.LoadBalancerConfigOpts{
		OctaviaEngine: octaviaEngine,
	}

	// build insecure registries
	var insecureRegistries []string
	if _, ok := d.GetOk("insecure_registries"); ok {
		insecureRegistries, err = getAsStringSlice(d, "insecure_registries")
		if err != nil {
			return diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
				AttributePath: cty.Path{
					cty.GetAttrStep{Name: "insecure_registries"},
				},
			}}
		}
	}

	// build create options
	createOpts := clusters.CreateOpts{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Version:            d.Get("version").(string),
		Labels:             labels,
		MasterSpec:         masterSpec,
		DeploymentType:     deploymentType,
		NetworkConfig:      networkConfig,
		LoadBalancerConfig: loadBalancerConfig,
		InsecureRegistries: insecureRegistries,
	}

	// add UUID if specified
	if uuid, ok := d.GetOk("uuid"); ok {
		createOpts.UUID = uuid.(string)
	}

	// create cluster
	clusterID, err := clusters.Create(containerInfraClientV2, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating vkcs_kubernetes_cluster_v2: %s", err)
	}

	// store the cluster ID
	d.SetId(clusterID)

	// wait for cluster to become active
	stateConf := &retry.StateChangeConf{
		Pending:      []string{clusterStatusV2Provisioning, clusterStatusV2Starting, clusterStatusV2Reconciling},
		Target:       []string{clusterStatusV2Running},
		Refresh:      kubernetesStateRefreshFuncV2(containerInfraClientV2, clusterID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        createUpdateDelayV2 * time.Minute,
		PollInterval: createUpdatePollIntervalV2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for vkcs_kubernetes_cluster_v2 %s to become ready: %s", clusterID, err)
	}

	log.Printf("[DEBUG] Created vkcs_kubernetes_cluster_v2 %s", clusterID)
	if diagns := resourceKubernetesClusterV2Read(ctx, d, meta); diagns != nil {
		return diagns
	}

	return nil
}

func resourceKubernetesClusterV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// creates client for managed-k8s
	config := meta.(clients.Config)
	containerInfraClientV2, err := config.ContainerInfraV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra v2 client: %s", err)
	}

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

func resourceKubernetesClusterV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// creates client for managed-k8s
	config := meta.(clients.Config)
	containerInfraClientV2, err := config.ContainerInfraV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra v2 client: %s", err)
	}

	// build options
	upgradeOpts := clusters.UpgradeOpts{}
	scaleOpts := clusters.ScaleOpts{}

	// check for version update
	if d.HasChange("version") {
		upgradeOpts.Version = d.Get("version").(string)
	}

	// check for master_flavor update
	if d.HasChange("master_flavor") {
		masterSpec := clusters.MasterSpecOpts{
			Engine: clusters.MasterEngineOpts{
				NovaEngine: clusters.NovaEngineOpts{
					FlavorID: d.Get("master_flavor").(string),
				},
			},
			Replicas: d.Get("master_count").(int),
		}
		scaleOpts.MasterSpec = masterSpec
	}

	// if there are no changes, return early
	if upgradeOpts.Version == "" && scaleOpts.MasterSpec.Engine.NovaEngine.FlavorID == "" {
		log.Printf("[DEBUG] No changes detected for vkcs_kubernetes_cluster_v2 %s", d.Id())
		return resourceKubernetesClusterV2Read(ctx, d, meta)
	}

	if upgradeOpts.Version != "" {
		// ugrade cluster
		err = clusters.Upgrade(containerInfraClientV2, d.Id(), upgradeOpts)
		if err != nil {
			return diag.Errorf("error upgrading vkcs_kubernetes_cluster_v2: %s", err)
		}
	}

	if scaleOpts.MasterSpec.Engine.NovaEngine.FlavorID != "" {
		// scale cluster
		err = clusters.Scale(containerInfraClientV2, d.Id(), scaleOpts)
		if err != nil {
			return diag.Errorf("error scaling vkcs_kubernetes_cluster_v2: %s", err)
		}
	}

	// wait for cluster to become active again after upgrade
	stateConf := &retry.StateChangeConf{
		Pending:      []string{clusterStatusV2Reconciling},
		Target:       []string{clusterStatusV2Running},
		Refresh:      kubernetesStateRefreshFuncV2(containerInfraClientV2, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        createUpdateDelayV2 * time.Minute,
		PollInterval: createUpdatePollIntervalV2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for vkcs_kubernetes_cluster_v2 %s to become ready after update: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Updated vkcs_kubernetes_cluster_v2 %s", d.Id())

	return resourceKubernetesClusterV2Read(ctx, d, meta)
}

func resourceKubernetesClusterV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// creates client for managed-k8s
	config := meta.(clients.Config)
	containerInfraClientV2, err := config.ContainerInfraV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra v2 client: %s", err)
	}

	// delete cluster
	err = clusters.Delete(containerInfraClientV2, d.Id())
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "error deleting vkcs_kubernetes_cluster_v2"))
	}

	// wait for cluster to be deleted
	stateConf := &retry.StateChangeConf{
		Pending:      []string{clusterStatusV2Running, clusterStatusV2Deleting, clusterStatusV2Reconciling},
		Target:       []string{clusterStatusV2Deleted},
		Refresh:      kubernetesStateRefreshFuncV2(containerInfraClientV2, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        deleteDelayV2 * time.Second,
		PollInterval: deletePollIntervalV2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for vkcs_kubernetes_cluster_v2 %s to be deleted: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Deleted vkcs_kubernetes_cluster_v2 %s", d.Id())
	return nil
}

func customizeClusterV2Diff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if d.HasChange("cluster_type") || d.HasChange("master_count") {
		clusterType := d.Get("cluster_type").(string)
		masterCount := d.Get("master_count").(int)

		switch clusterType {
		case "standard":
			if masterCount != 1 {
				return fmt.Errorf("standard cluster requires exactly 1 master node, got %d", masterCount)
			}
		case "regional":
			if masterCount != 3 && masterCount != 5 {
				return fmt.Errorf("regional cluster requires 3 or 5 master nodes, got %d", masterCount)
			}
		}
	}

	// validate availability_zones for regional cluster
	if d.HasChange("cluster_type") || d.HasChange("availability_zones") {
		clusterType := d.Get("cluster_type").(string)
		availabilityZones, err := getAsStringSlice(d, "availability_zones")
		if err != nil {
			return err
		}

		if clusterType == "regional" && len(availabilityZones) != 3 {
			// regional cluster must have 3 a zones
			return fmt.Errorf("regional cluster requires 3 availability zones, got %d", len(availabilityZones))
		}

		if clusterType == "standard" && len(availabilityZones) != 1 {
			return fmt.Errorf("standard cluster requires exactly 1 availability zone, got %d", len(availabilityZones))
		}
	}

	// !!! Note: node_groups validation removed as it's not part of the current schema
	// node groups will be handled in a separate resource vkcs_kubernetes_node_group_v2

	// validate external network is specified when public IP is enabled
	if d.HasChange("enable_public_ip") {
		enablePublicIP := d.Get("enable_public_ip").(bool)
		externalNetworkID := d.Get("external_network_id").(string)

		if enablePublicIP && externalNetworkID == "" {
			return fmt.Errorf("external_network_id must be specified when enable_public_ip is true")
		}
	}

	oldVersion, newVersion := d.GetChange("version")
	oldMasterFlavor, newMasterFlavor := d.GetChange("master_flavor")
	versionAndFlavorWereNotEmpty := oldVersion.(string) != "" && newVersion.(string) != "" &&
		oldMasterFlavor.(string) != "" && newMasterFlavor.(string) != ""

	versionAndFlavorChanged := d.HasChange("version") && d.HasChange("master_flavor")

	if versionAndFlavorChanged && versionAndFlavorWereNotEmpty {
		return fmt.Errorf("simultaneous scaling and cluster upgrade is not available")
	}

	return nil
}
