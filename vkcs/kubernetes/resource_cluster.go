package kubernetes

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/valid"
)

const (
	operationCreate          = 60
	operationUpdate          = 60
	operationDelete          = 30
	createUpdateDelay        = 1
	createUpdatePollInterval = 20
	deleteDelay              = 30
	nodeGroupDeleteDelay     = 10
	deletePollInterval       = 10
)

type clusterStatus string

var (
	clusterStatusDeleting     clusterStatus = "DELETING"
	clusterStatusDeleted      clusterStatus = "DELETED"
	clusterStatusNotFound     clusterStatus = "NOT_FOUND"
	clusterStatusReconciling  clusterStatus = "RECONCILING"
	clusterStatusProvisioning clusterStatus = "PROVISIONING"
	clusterStatusRunning      clusterStatus = "RUNNING"
	clusterStatusError        clusterStatus = "ERROR"
	clusterStatusShutoff      clusterStatus = "SHUTOFF"
)

var stateStatusMap = map[clusterStatus]string{
	clusterStatusRunning: "turn_on_cluster",
	clusterStatusShutoff: "turn_off_cluster",
}

func ResourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesClusterCreate,
		ReadContext:   resourceKubernetesClusterRead,
		UpdateContext: resourceKubernetesClusterUpdate,
		DeleteContext: resourceKubernetesClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(operationCreate * time.Minute),
			Update: schema.DefaultTimeout(operationUpdate * time.Minute),
			Delete: schema.DefaultTimeout(operationDelete * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Region to use for the cluster. Default is a region configured for provider.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					name := val.(string)
					if err := valid.ClusterName(name); err != nil {
						errs = append(errs, err)
					}
					return
				},
				Description: "The name of the cluster. Changing this creates a new cluster. Should match the pattern `^[a-zA-Z][a-zA-Z0-9_.-]*$`.",
			},
			"project_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Description: "The project of the cluster.",
			},
			"user_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Description: "The user of the cluster.",
			},
			"created_at": {
				Type:        schema.TypeString,
				ForceNew:    false,
				Computed:    true,
				Description: "The time at which cluster was created.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				ForceNew:    false,
				Computed:    true,
				Description: "The time at which cluster was created.",
			},
			"api_address": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Description: "COE API address.",
			},
			"cluster_template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "The UUID of the Kubernetes cluster template. It can be obtained using the cluster_template data source.",
			},
			"master_flavor": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				Description: "The UUID of a flavor for the master nodes. If master_flavor is not present, value from cluster_template will be used.",
			},
			"keypair": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "The name of the Compute service SSH keypair. Changing this creates a new cluster.",
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Description: "The list of optional key value pairs representing additional properties of the cluster." +
					" _note_ Updating this attribute will not immediately apply the changes; these options will be used when recreating or deleting cluster nodes, for example, during an upgrade operation.\n\n" +
					"  * `calico_ipv4pool` to set subnet where pods will be created. Default 10.100.0.0/16. _note_ Updating this value while the cluster is running is dangerous because it can lead to loss of connectivity of the cluster nodes.\n" +
					"  * `clean_volumes` to remove pvc volumes when deleting a cluster. Default False. _note_ Changes to this value will be applied immediately.\n" +
					"  * `cloud_monitoring` to enable cloud monitoring feature. Default False.\n" +
					"  * `etcd_volume_size` to set etcd volume size in GB. Default 10.\n" +
					"  * `kube_log_level` to set log level for kubelet in range 0 to 8. Default 0.\n" +
					"  * `master_volume_size` to set master vm volume size in GB. Default 50.\n" +
					"  * `cluster_node_volume_type` to set master vm volume type. Default ceph-hdd.\n",
			},
			"all_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The read-only map of all cluster labels.",
			},
			"master_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The number of master nodes for the cluster. Changing this creates a new cluster.",
			},
			"master_addresses": {
				Type:        schema.TypeList,
				ForceNew:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "IP addresses of the master node of the cluster.",
			},
			"stack_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Description: "UUID of the Orchestration service stack.",
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "UUID of the cluster's network.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "UUID of the cluster's subnet.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				Description: "Current state of a cluster. Changing this to `RUNNING` or `SHUTOFF` will turn cluster on/off.",
			},
			"pods_network_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Network cidr of k8s virtual network",
			},
			"floating_ip_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Floating ip is enabled.",
			},
			"api_lb_vip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "API LoadBalancer vip. IP address field.",
			},
			"api_lb_fip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "API LoadBalancer fip. IP address field.",
			},
			"ingress_floating_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Floating IP created for ingress service.",
			},
			"registry_auth_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Docker registry access password.",
			},
			"loadbalancer_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The UUID of the load balancer's subnet. Changing this creates new cluster.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Availability zone of the cluster.",
			},
			"insecure_registries": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Addresses of registries from which you can download images without checking certificates. Changing this creates a new cluster.",
			},
			"dns_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Custom DNS cluster domain. Changing this creates a new cluster.",
			},
			"sync_security_policy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enables syncing of security policies of cluster. Default value is false.",
			},
		},
		Description: "Provides a kubernetes cluster resource. This can be used to create, modify and delete kubernetes clusters.",
	}
}

func resourceKubernetesClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	containerInfraClient, err := config.ContainerInfraV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra client: %s", err)
	}

	// Get and check labels map.
	rawLabels := d.Get("labels").(map[string]interface{})
	labels, err := extractKubernetesLabelsMap(rawLabels)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := clusters.CreateOpts{
		ClusterTemplateID:    d.Get("cluster_template_id").(string),
		MasterFlavorID:       d.Get("master_flavor").(string),
		Keypair:              d.Get("keypair").(string),
		Labels:               labels,
		Name:                 d.Get("name").(string),
		NetworkID:            d.Get("network_id").(string),
		SubnetID:             d.Get("subnet_id").(string),
		PodsNetworkCidr:      d.Get("pods_network_cidr").(string),
		FloatingIPEnabled:    d.Get("floating_ip_enabled").(bool),
		APILBVIP:             d.Get("api_lb_vip").(string),
		APILBFIP:             d.Get("api_lb_fip").(string),
		LoadbalancerSubnetID: d.Get("loadbalancer_subnet_id").(string),
		RegistryAuthPassword: d.Get("registry_auth_password").(string),
		AvailabilityZone:     d.Get("availability_zone").(string),
		DNSDomain:            d.Get("dns_domain").(string),
	}

	if masterCount, ok := d.GetOk("master_count"); ok {
		mCount := masterCount.(int)
		if mCount < 1 {
			return diag.Errorf("master_count if set must be greater or equal 1: %s", err)
		}
		createOpts.MasterCount = mCount
	}

	if registriesRaw, ok := d.GetOk("insecure_registries"); ok {
		registries := registriesRaw.([]interface{})
		insecureRegistries := make([]string, 0, len(registries))
		for _, val := range registries {
			insecureRegistries = append(insecureRegistries, val.(string))
		}
		createOpts.InsecureRegistries = insecureRegistries
	}

	if syncSecurityPolicyRaw, ok := d.GetOk("sync_security_policy"); ok {
		syncSecurityPolicy := syncSecurityPolicyRaw.(bool)
		createOpts.SecurityPolicySyncEnabled = &syncSecurityPolicy
	}

	s, err := clusters.Create(containerInfraClient, &createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating vkcs_kubernetes_cluster: %s", err)
	}

	// Store the cluster ID.
	d.SetId(s)

	stateConf := &retry.StateChangeConf{
		Pending:      []string{string(clusterStatusProvisioning)},
		Target:       []string{string(clusterStatusRunning)},
		Refresh:      kubernetesStateRefreshFunc(containerInfraClient, s),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        createUpdateDelay * time.Minute,
		PollInterval: createUpdatePollInterval * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for vkcs_kubernetes_cluster %s to become ready: %s", s, err)
	}

	log.Printf("[DEBUG] Created vkcs_kubernetes_cluster %s", s)
	return resourceKubernetesClusterRead(ctx, d, meta)
}

func resourceKubernetesClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	containerInfraClient, err := config.ContainerInfraV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra client: %s", err)
	}

	cluster, err := clusters.Get(containerInfraClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "error retrieving vkcs_kubernetes_cluster"))
	}

	log.Printf("[DEBUG] retrieved vkcs_kubernetes_cluster %s: %#v", d.Id(), cluster)

	d.Set("name", cluster.Name)
	d.Set("api_address", cluster.APIAddress)
	d.Set("cluster_template_id", cluster.ClusterTemplateID)
	d.Set("master_flavor", cluster.MasterFlavorID)
	d.Set("keypair", cluster.KeyPair)
	d.Set("master_count", cluster.MasterCount)
	d.Set("master_addresses", cluster.MasterAddresses)
	d.Set("stack_id", cluster.StackID)
	d.Set("status", cluster.NewStatus)
	d.Set("pods_network_cidr", cluster.PodsNetworkCidr)
	d.Set("floating_ip_enabled", cluster.FloatingIPEnabled)
	d.Set("api_lb_vip", cluster.APILBVIP)
	d.Set("api_lb_fip", cluster.APILBFIP)
	d.Set("ingress_floating_ip", cluster.IngressFloatingIP)
	d.Set("loadbalancer_subnet_id", cluster.LoadbalancerSubnetID)
	d.Set("registry_auth_password", cluster.RegistryAuthPassword)
	d.Set("availability_zone", cluster.AvailabilityZone)
	d.Set("region", util.GetRegion(d, config))
	d.Set("insecure_registries", cluster.InsecureRegistries)
	d.Set("dns_domain", cluster.DNSDomain)
	d.Set("sync_security_policy", cluster.SecurityPolicySyncEnabled)

	// Allow to read old api clusters
	if cluster.NetworkID != "" {
		d.Set("network_id", cluster.NetworkID)
	} else {
		d.Set("network_id", cluster.Labels["fixed_network"])
	}
	if cluster.SubnetID != "" {
		d.Set("subnet_id", cluster.SubnetID)
	} else {
		d.Set("subnet_id", cluster.Labels["fixed_subnet"])
	}

	if err := d.Set("created_at", util.GetTimestamp(&cluster.CreatedAt)); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_kubernetes_cluster created_at: %s", err)
	}
	if err := d.Set("updated_at", util.GetTimestamp(&cluster.UpdatedAt)); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_kubernetes_cluster updated_at: %s", err)
	}

	// Get and check labels map.
	rawLabels := d.Get("labels").(map[string]interface{})
	labels, err := extractKubernetesLabelsMap(rawLabels)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range cluster.Labels {
		if _, ok := labels[k]; ok {
			labels[k] = v
		}
	}

	d.Set("labels", labels)
	d.Set("all_labels", cluster.Labels)

	return nil
}

func resourceKubernetesClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	containerInfraClient, err := config.ContainerInfraV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra client: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Refresh:      kubernetesStateRefreshFunc(containerInfraClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        createUpdateDelay * time.Minute,
		PollInterval: createUpdatePollInterval * time.Second,
		Pending:      []string{string(clusterStatusReconciling)},
		Target:       []string{string(clusterStatusRunning)},
	}

	cluster, err := clusters.Get(containerInfraClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("error retrieving cluster: %s", err)
	}

	switch cluster.NewStatus {
	case string(clusterStatusShutoff):
		changed, err := checkForStatus(ctx, d, containerInfraClient, cluster)
		if err != nil {
			return diag.FromErr(err)
		}
		if changed {
			err := checkForClusterTemplateID(ctx, d, containerInfraClient, stateConf)
			if err != nil {
				return diag.FromErr(err)
			}
			err = checkForMasterFlavor(ctx, d, containerInfraClient, stateConf)
			if err != nil {
				return diag.FromErr(err)
			}
			err = checkForUpdate(ctx, d, containerInfraClient, stateConf)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			return diag.Errorf("changing cluster attributes is prohibited when cluster has SHUTOFF status")
		}
	case string(clusterStatusRunning):
		err := checkForClusterTemplateID(ctx, d, containerInfraClient, stateConf)
		if err != nil {
			return diag.FromErr(err)
		}
		err = checkForMasterFlavor(ctx, d, containerInfraClient, stateConf)
		if err != nil {
			return diag.FromErr(err)
		}
		err = checkForUpdate(ctx, d, containerInfraClient, stateConf)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = checkForStatus(ctx, d, containerInfraClient, cluster)
		if err != nil {
			return diag.FromErr(err)
		}
	default:
		return diag.Errorf("changes in cluster are prohibited when status is not RUNNING/SHUTOFF; current status: %s", cluster.NewStatus)
	}

	return resourceKubernetesClusterRead(ctx, d, meta)
}

func checkForClusterTemplateID(ctx context.Context, d *schema.ResourceData, containerInfraClient *gophercloud.ServiceClient, stateConf *retry.StateChangeConf) error {
	if d.HasChange("cluster_template_id") {
		upgradeOpts := clusters.UpgradeOpts{
			ClusterTemplateID: d.Get("cluster_template_id").(string),
			RollingEnabled:    true,
		}

		_, err := clusters.Upgrade(containerInfraClient, d.Id(), &upgradeOpts).Extract()
		if err != nil {
			return fmt.Errorf("error upgrade cluster : %s", err)
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmt.Errorf(
				"error waiting for vkcs_kubernetes_cluster %s to become upgraded: %s", d.Id(), err)
		}
	}
	return nil
}

func checkForMasterFlavor(ctx context.Context, d *schema.ResourceData, containerInfraClient *gophercloud.ServiceClient, stateConf *retry.StateChangeConf) error {
	if d.HasChange("master_flavor") {
		upgradeOpts := clusters.ActionsBaseOpts{
			Action: "resize_masters",
			Payload: map[string]string{
				"flavor": d.Get("master_flavor").(string),
			},
		}

		_, err := clusters.UpdateMasters(containerInfraClient, d.Id(), &upgradeOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating cluster's falvor : %s", err)
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmt.Errorf(
				"error waiting for vkcs_kubernetes_cluster %s to become updated: %s", d.Id(), err)
		}
	}
	return nil
}

func checkForStatus(ctx context.Context, d *schema.ResourceData, containerInfraClient *gophercloud.ServiceClient, cluster *clusters.Cluster) (bool, error) {

	turnOffConf := &retry.StateChangeConf{
		Refresh:      kubernetesStateRefreshFunc(containerInfraClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        createUpdateDelay * time.Minute,
		PollInterval: createUpdatePollInterval * time.Second,
		Pending:      []string{string(clusterStatusRunning)},
		Target:       []string{string(clusterStatusShutoff)},
	}

	turnOnConf := &retry.StateChangeConf{
		Refresh:      kubernetesStateRefreshFunc(containerInfraClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        createUpdateDelay * time.Minute,
		PollInterval: createUpdatePollInterval * time.Second,
		Pending:      []string{string(clusterStatusShutoff)},
		Target:       []string{string(clusterStatusRunning)},
	}

	if d.HasChange("status") {
		currentStatus := clusterStatus(d.Get("status").(string))
		if cluster.NewStatus != string(clusterStatusRunning) && cluster.NewStatus != string(clusterStatusShutoff) {
			return false, fmt.Errorf("turning on/off is prohibited due to cluster's status %s", cluster.NewStatus)
		}
		switchStateOpts := clusters.ActionsBaseOpts{
			Action: stateStatusMap[currentStatus],
		}
		_, err := clusters.SwitchState(containerInfraClient, d.Id(), &switchStateOpts).Extract()
		if err != nil {
			return false, fmt.Errorf("error during switching state: %s", err)
		}

		var switchStateConf *retry.StateChangeConf
		switch currentStatus {
		case clusterStatusRunning:
			switchStateConf = turnOnConf
		case clusterStatusShutoff:
			switchStateConf = turnOffConf
		default:
			return false, fmt.Errorf("unknown status provided: %s", currentStatus)
		}

		_, err = switchStateConf.WaitForStateContext(ctx)
		if err != nil {
			return false, fmt.Errorf(
				"error waiting for vkcs_kubernetes_cluster %s to become updated: %s", d.Id(), err)
		}
		return true, nil

	}
	return false, nil
}

func checkForUpdate(ctx context.Context, d *schema.ResourceData, containerInfraClient *gophercloud.ServiceClient, stateConf *retry.StateChangeConf) error {
	updateOpts := []clusters.OptsBuilder{}

	if d.HasChange("labels") {
		rawLabels := d.Get("labels").(map[string]interface{})
		labels, err := extractKubernetesLabelsMap(rawLabels)
		if err != nil {
			return err
		}

		rawAllLabels := d.Get("all_labels").(map[string]interface{})
		allLabels, err := extractKubernetesLabelsMap(rawAllLabels)
		if err != nil {
			return err
		}

		for k, v := range labels {
			allLabels[k] = v
		}

		updateOpts = append(updateOpts, &clusters.UpdateOpts{
			Op:    clusters.ReplaceOp,
			Path:  "/labels",
			Value: allLabels,
		})
	}

	if d.HasChange("sync_security_policy") {
		syncSecurityPolicy := d.Get("sync_security_policy").(bool)
		updateOpts = append(updateOpts, &clusters.UpdateOpts{
			Op:    clusters.ReplaceOp,
			Path:  "/security_policy_sync_enabled",
			Value: syncSecurityPolicy,
		})
	}

	if len(updateOpts) > 0 {
		log.Printf("[DEBUG] Updating vkcs_kubernetes_cluster %s with options: %#v", d.Id(), updateOpts)

		_, err := clusters.Update(containerInfraClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating vkcs_kubernetes_cluster %s: %s", d.Id(), err)
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmt.Errorf("error waiting for vkcs_kubernetes_cluster %s to be updated: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceKubernetesClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	client, err := config.ContainerInfraV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("failed to get container infra client: %s", err)
	}

	if err := clusters.Delete(client, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "failed to delete vkcs_kubernetes_cluster"))
	}

	stateConf := &retry.StateChangeConf{
		Pending:      []string{string(clusterStatusReconciling), string(clusterStatusRunning), string(clusterStatusDeleting), string(clusterStatusDeleted)},
		Target:       []string{string(clusterStatusNotFound)},
		Refresh:      kubernetesStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        deleteDelay * time.Second,
		PollInterval: deletePollInterval * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf(
			"error waiting for vkcs_kubernetes_cluster %s to become deleted: %s", d.Id(), err)
	}

	return nil
}
