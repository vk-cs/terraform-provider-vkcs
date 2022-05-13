package vkcs

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			},
			"max_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"min_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"autoscaling_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},
			"uuid": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				ForceNew: false,
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
			"availability_zones": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceKubernetesNodeGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	containerInfraClient, err := config.ContainerInfraV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra client: %s", err)
	}

	createOpts := nodeGroupCreateOpts{
		ClusterID:   d.Get("cluster_id").(string),
		FlavorID:    d.Get("flavor_id").(string),
		MaxNodes:    d.Get("max_nodes").(int),
		MinNodes:    d.Get("min_nodes").(int),
		VolumeSize:  d.Get("volume_size").(int),
		VolumeType:  d.Get("volume_type").(string),
		Autoscaling: d.Get("autoscaling_enabled").(bool),
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

	s, err := nodeGroupCreate(containerInfraClient, &createOpts).Extract()
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
	config := meta.(configer)
	containerInfraClient, err := config.ContainerInfraV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra client: %s", err)
	}

	s, err := nodeGroupGet(containerInfraClient, d.Id()).Extract()
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

	if err := d.Set("created_at", getTimestamp(&s.CreatedAt)); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_kubernetes_node_group created_at: %s", err)
	}
	if err := d.Set("updated_at", getTimestamp(&s.UpdatedAt)); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_kubernetes_node_group updated_at: %s", err)
	}

	return nil
}

func resourceKubernetesNodeGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
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
		s, err := nodeGroupGet(containerInfraClient, d.Id()).Extract()
		if err != nil {
			return diag.Errorf("error retrieving kubernetes_node_group : %s", err)
		}
		scaleOpts := nodeGroupScaleOpts{
			Delta: d.Get("node_count").(int) - s.NodeCount,
		}

		_, err = nodeGroupScale(containerInfraClient, d.Id(), &scaleOpts).Extract()
		if err != nil {
			return diag.Errorf("error scaling vkcs_kubernetes_node_group : %s", err)
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"error waiting for vkcs_kubernetes_node_group %s to become scaled: %s", d.Id(), err)
		}

	}

	var patchOpts nodeGroupClusterPatchOpts

	if d.HasChange("max_nodes") {
		patchOpts = append(patchOpts, nodeGroupPatchParams{
			Path:  "/max_nodes",
			Value: d.Get("max_nodes").(int),
			Op:    "replace",
		})
	}

	if d.HasChange("min_nodes") {
		patchOpts = append(patchOpts, nodeGroupPatchParams{
			Path:  "/min_nodes",
			Value: d.Get("min_nodes").(int),
			Op:    "replace",
		})
	}

	if d.HasChange("autoscaling_enabled") {
		patchOpts = append(patchOpts, nodeGroupPatchParams{
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

		patchOpts = append(patchOpts, nodeGroupPatchParams{
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

		patchOpts = append(patchOpts, nodeGroupPatchParams{
			Path:  "/taints",
			Value: taints,
			Op:    "replace",
		})
	}

	if len(patchOpts) > 0 {
		_, err := nodeGroupPatch(containerInfraClient, d.Id(), &patchOpts).Extract()
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
	config := meta.(configer)
	containerInfraClient, err := config.ContainerInfraV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra client: %s", err)
	}

	if err := nodeGroupDelete(containerInfraClient, d.Id()).ExtractErr(); err != nil {
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
