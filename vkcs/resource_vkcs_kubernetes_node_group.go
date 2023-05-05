package vkcs

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/nodegroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/randutil"
)

func resourceKubernetesNodeGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesNodeGroupCreate,
		ReadContext:   resourceKubernetesNodeGroupRead,
		UpdateContext: resourceKubernetesNodeGroupUpdate,
		DeleteContext: resourceKubernetesNodeGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(operationCreate * time.Minute),
			Update: schema.DefaultTimeout(operationUpdate * time.Minute),
			Delete: schema.DefaultTimeout(operationDelete * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The UUID of the existing cluster.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of node group to create. Changing this will force to create a new node group.",
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Description: "The list of objects representing representing additional properties of the node group. Each object should have attribute \"key\". Object may also have optional attribute \"value\".",
			},
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Description: "The list of objects representing node group taints. Each object should have following attributes: key, value, effect.",
			},
			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Suppress diff if node_count is managed by autoscaler when updating
					if d.Get("autoscaling_enabled").(bool) && old != "" {
						return true
					}
					return false
				},
				Description: "The node count for this node group. Should be greater than 0. If `autoscaling_enabled` parameter is set, this attribute will be ignored during update.",
			},
			"max_nodes": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				Description: "The maximum allowed nodes for this node group.",
			},
			"min_nodes": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				Description: "The minimum allowed nodes for this node group. Default to 0 if not set.",
			},
			"volume_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The size in GB for volume to load nodes from. Changing this will force to create a new node group.",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The volume type to load nodes from. Changing this will force to create a new node group.",
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The flavor UUID of this node group.",
			},
			"autoscaling_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Default:     false,
				Description: "Determines whether the autoscaling is enabled.",
			},
			"uuid": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Description: "The UUID of the cluster's node group.",
			},
			"state": {
				Type:        schema.TypeString,
				ForceNew:    false,
				Computed:    true,
				Description: "Determines current state of node group (RUNNING, SHUTOFF, ERROR).",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time at which node group was created.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time at which node group was created.",
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "The list of availability zones of the node group. Zones `MS1` and  `GZ1` are available. By default, node group is being created at cluster's zone.\n" +
					"**Important:** Receiving default AZ add it manually to your main.tf config to sync it with state to avoid node groups force recreation in the future.",
			},
			"max_node_unavailable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The maximum number of nodes that can fail during an upgrade. The default value is 25 percent.",
			},
		},
		Description: "Provides a cluster node group resource. This can be used to create, modify and delete cluster's node group.",
	}
}

func resourceKubernetesNodeGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	containerInfraClient, err := config.ContainerInfraV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra client: %s", err)
	}

	createOpts := nodegroups.CreateOpts{
		ClusterID:          d.Get("cluster_id").(string),
		FlavorID:           d.Get("flavor_id").(string),
		MaxNodes:           d.Get("max_nodes").(int),
		MinNodes:           d.Get("min_nodes").(int),
		VolumeSize:         d.Get("volume_size").(int),
		VolumeType:         d.Get("volume_type").(string),
		Autoscaling:        d.Get("autoscaling_enabled").(bool),
		MaxNodeUnavailable: d.Get("max_node_unavailable").(int),
	}

	if zonesRaw, ok := d.GetOk("availability_zones"); ok {
		zones := zonesRaw.([]interface{})
		az := make([]string, 0, len(zones))
		for _, zone := range zones {
			z := zone.(string)
			az = append(az, z)
		}
		createOpts.AvailabilityZones = az
	}

	if ngName, ok := d.GetOk("name"); ok {
		createOpts.Name = ngName.(string)
	} else {
		createOpts.Name = "ng-" + randutil.RandomName(5)
	}

	if lab, labOk := d.GetOk("labels"); labOk {
		rawLabels := lab.([]interface{})
		labels, err := extractNodeGroupLabelsList(rawLabels)
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.Labels = labels
	}

	if tnt, tntOk := d.GetOk("taints"); tntOk {
		rawTaints := tnt.([]interface{})
		taints, err := extractNodeGroupTaintsList(rawTaints)
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.Taints = taints
	}

	if nodeCount := d.Get("node_count").(int); nodeCount > 0 {
		createOpts.NodeCount = nodeCount
	} else {
		return diag.Errorf("node_count parameter must be > 0")
	}

	s, err := nodegroups.Create(containerInfraClient, &createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating vkcs_kubernetes_node_group: %s", err)
	}

	// Store the node Group ID.
	d.SetId(s.UUID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{string(clusterStatusReconciling)},
		Target:       []string{string(clusterStatusRunning)},
		Refresh:      kubernetesStateRefreshFunc(containerInfraClient, s.ClusterID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        createUpdateDelay * time.Minute,
		PollInterval: createUpdatePollInterval * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for vkcs_kubernetes_cluster %s to become ready: %s", s.ClusterID, err)
	}

	log.Printf("[DEBUG] Created vkcs_kubernetes_node_group %s", s.UUID)
	return resourceKubernetesNodeGroupRead(ctx, d, meta)
}

func resourceKubernetesNodeGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	containerInfraClient, err := config.ContainerInfraV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra client: %s", err)
	}

	s, err := nodegroups.Get(containerInfraClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "error retrieving vkcs_kubernetes_node_group"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_kubernetes_node_group %s: %#v", d.Id(), s)

	// Get and check labels list.
	rawLabels := d.Get("labels").([]interface{})
	labels, err := extractNodeGroupLabelsList(rawLabels)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("labels", flattenNodeGroupLabelsList(labels)); err != nil {
		return diag.Errorf("unable to set vkcs_kubernetes_node_group labels: %s", err)
	}

	// Get and check taints list.
	rawTaints := d.Get("taints").([]interface{})
	taints, err := extractNodeGroupTaintsList(rawTaints)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("taints", flattenNodeGroupTaintsList(taints)); err != nil {
		return diag.Errorf("unable to set vkcs_kubernetes_node_group taints: %s", err)
	}

	d.Set("name", s.Name)
	d.Set("node_count", s.NodeCount)
	d.Set("max_nodes", s.MaxNodes)
	d.Set("min_nodes", s.MinNodes)
	d.Set("volume_size", s.VolumeSize)
	d.Set("volume_type", s.VolumeType)
	d.Set("flavor_id", s.FlavorID)
	d.Set("autoscaling_enabled", s.Autoscaling)
	d.Set("cluster_id", s.ClusterID)
	d.Set("availability_zones", s.AvailabilityZones)
	d.Set("max_node_unavailable", s.MaxNodeUnavailable)

	if err := d.Set("created_at", getTimestamp(&s.CreatedAt)); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_kubernetes_node_group created_at: %s", err)
	}
	if err := d.Set("updated_at", getTimestamp(&s.UpdatedAt)); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_kubernetes_node_group updated_at: %s", err)
	}

	return nil
}

func resourceKubernetesNodeGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	containerInfraClient, err := config.ContainerInfraV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Refresh:      kubernetesStateRefreshFunc(containerInfraClient, d.Get("cluster_id").(string)),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        createUpdateDelay * time.Minute,
		PollInterval: createUpdatePollInterval * time.Second,
		Pending:      []string{string(clusterStatusReconciling)},
		Target:       []string{string(clusterStatusRunning)},
	}

	if d.HasChange("node_count") {
		s, err := nodegroups.Get(containerInfraClient, d.Id()).Extract()
		if err != nil {
			return diag.Errorf("error retrieving kubernetes_node_group : %s", err)
		}
		scaleOpts := nodegroups.ScaleOpts{
			Delta: d.Get("node_count").(int) - s.NodeCount,
		}

		_, err = nodegroups.Scale(containerInfraClient, d.Id(), &scaleOpts).Extract()
		if err != nil {
			return diag.Errorf("error scaling vkcs_kubernetes_node_group : %s", err)
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"error waiting for vkcs_kubernetes_node_group %s to become scaled: %s", d.Id(), err)
		}

	}

	var patchOpts nodegroups.PatchOpts

	if d.HasChange("max_nodes") {
		patchOpts = append(patchOpts, nodegroups.PatchParams{
			Path:  "/max_nodes",
			Value: d.Get("max_nodes").(int),
			Op:    "replace",
		})
	}

	if d.HasChange("min_nodes") {
		patchOpts = append(patchOpts, nodegroups.PatchParams{
			Path:  "/min_nodes",
			Value: d.Get("min_nodes").(int),
			Op:    "replace",
		})
	}

	if d.HasChange("autoscaling_enabled") {
		patchOpts = append(patchOpts, nodegroups.PatchParams{
			Path:  "/autoscaling_enabled",
			Value: strconv.FormatBool(d.Get("autoscaling_enabled").(bool)),
			Op:    "replace",
		})
	}

	if d.HasChange("labels") {
		rawLabels := d.Get("labels").([]interface{})
		labels, err := extractNodeGroupLabelsList(rawLabels)
		if err != nil {
			return diag.FromErr(err)
		}

		patchOpts = append(patchOpts, nodegroups.PatchParams{
			Path:  "/labels",
			Value: labels,
			Op:    "replace",
		})
	}

	if d.HasChange("taints") {
		rawTaints := d.Get("taints").([]interface{})
		taints, err := extractNodeGroupTaintsList(rawTaints)
		if err != nil {
			return diag.FromErr(err)
		}

		patchOpts = append(patchOpts, nodegroups.PatchParams{
			Path:  "/taints",
			Value: taints,
			Op:    "replace",
		})
	}

	if d.HasChange("max_node_unavailable") {
		patchOpts = append(patchOpts, nodegroups.PatchParams{
			Path:  "/max_node_unavailable",
			Value: d.Get("max_node_unavailable").(int),
			Op:    "replace",
		})
	}

	if len(patchOpts) > 0 {
		_, err := nodegroups.Patch(containerInfraClient, d.Id(), &patchOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating vkcs_kubernetes_node_group : %s", err)
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"error waiting for vkcs_kubernetes_node_group %s to become updated: %s", d.Id(), err)
		}
	}

	return resourceKubernetesNodeGroupRead(ctx, d, meta)
}

func resourceKubernetesNodeGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	containerInfraClient, err := config.ContainerInfraV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra client: %s", err)
	}

	if err := nodegroups.Delete(containerInfraClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(checkDeleted(d, err, "error deleting vkcs_kubernetes_node_group"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{string(clusterStatusReconciling)},
		Target:       []string{string(clusterStatusRunning)},
		Refresh:      kubernetesStateRefreshFunc(containerInfraClient, d.Get("cluster_id").(string)),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        nodeGroupDeleteDelay * time.Second,
		PollInterval: deletePollInterval * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for vkcs_kubernetes_node_group %s to become deleted: %s", d.Id(), err)
	}

	return nil
}
