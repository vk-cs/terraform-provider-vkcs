package kubernetes

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

const (
	operationCreateNodeGroupV2 = 30
	operationUpdateNodeGroupV2 = 30
	operationDeleteNodeGroupV2 = 15
)

const (
	createNgDelayV2        = 3
	createNgPollIntervalV2 = 30
	updateNgDelayV2        = 3
	updateNgPollIntervalV2 = 30
)

func ResourceKubernetesNodeGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesNodeGroupV2Create,
		ReadContext:   resourceKubernetesNodeGroupV2Read,
		UpdateContext: resourceKubernetesNodeGroupV2Update,
		DeleteContext: resourceKubernetesNodeGroupV2Delete,
		CustomizeDiff: customizeNodeGroupV2Diff,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(operationCreateNodeGroupV2 * time.Minute),
			Update: schema.DefaultTimeout(operationUpdateNodeGroupV2 * time.Minute),
			Delete: schema.DefaultTimeout(operationDeleteNodeGroupV2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the target cluster. Changing this will force to create a new node group.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "The UUID of the node group. Generated automatically after node group creation.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					isNodeGroupNameV2,
				)),
				Description: "The name of node group to create. Changing this will force to create a new node group.",
			},
			"node_flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false, // при изменении будет запускаться вертикальное скалирование, а не пересоздание нод-группы
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					isUUID,
				)),
				Description: "Flavor ID of the nodes from node group.",
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// TODO: Добавить валидацию на указываемое значение с помощью ручки получения доступных AZs
				Description: "The availability zone of the node group.",
			},
			"scale_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
				ValidateFunc: validation.StringInSlice([]string{
					"fixed_scale",
					"auto_scale",
				}, false),
				Description: "Type of scaling for the node group. Must be either 'fixed_scale' or 'auto_scale'. If scale_type is 'auto_scale', then the condition 'auto_scale_min_size <= auto_scale_node_count <= auto_scale_max_size' must be met.",
			},
			"fixed_scale_node_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      false, // при изменении будет запускаться горизонтальное скалирование, а не пересоздание нод-группы
				ValidateFunc:  validation.IntAtLeast(0),
				ConflictsWith: []string{"auto_scale_min_size", "auto_scale_max_size", "auto_scale_node_count"},
				Description:   "The desired node count of the node group. Required if scale_type is 'fixed_scale'.",
			},
			"auto_scale_node_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      false, // при изменении будет запускаться горизонтальное скалирование, а не пересоздание нод-группы
				ValidateFunc:  validation.IntAtLeast(0),
				ConflictsWith: []string{"fixed_scale_node_count"},
				Description:   "When creating a cluster, this parameter allows to specify the desired initial number of nodes in the node group. During the cluster lifecycle, indicates the current number of nodes in the node group. Required if scale_type is 'auto_scale'.",
			},
			"auto_scale_min_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      false,
				ValidateFunc:  validation.IntAtLeast(0),
				ConflictsWith: []string{"fixed_scale_node_count"},
				Description:   "The minimum allowed nodes for this node group. Required if scale_type is 'auto_scale'.",
			},
			"auto_scale_max_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      false,
				ValidateFunc:  validation.IntAtLeast(0),
				ConflictsWith: []string{"fixed_scale_node_count"},
				Description:   "The maximum allowed nodes for this node group. Required if scale_type is 'auto_scale'.",
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: false,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// TODO: Добавить валидацию key, value
				Description: "The list of key-value pairs representing additional properties of the node group.",
			},
			// TODO: Добавить валидация для key, value
			"taints": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: false,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    false,
							Description: "Key of the taint.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    false,
							Description: "Value of the taint.",
						},
						"effect": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
							ValidateFunc: validation.StringInSlice([]string{
								"NoSchedule",
								"PreferNoSchedule",
								"NoExecute",
							}, false),
							Description: "Effect of the taint. Must be one of: 'NoSchedule', 'PreferNoSchedule', 'NoExecute'.",
						},
					},
				},
				Description: "The list of objects representing node group taints. Each object should have following attributes: key, value, effect.",
			},
			"parallel_upgrade_chunk": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    false,
				Description: "The maximum percent of nodes that can be unavailble during an upgrade.",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.IntBetween(1, 100),
				),
			},
			// TODO: Добавить ручку, которая возвращает список типов диска для запуска воркер-ноды
			"disk_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The root volume type to load nodes from.",
			},
			"disk_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				Description:  "The size in GB for volume to load nodes from.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Region to use for the node group. Default is a region configured for provider.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    false,
				Description: "The time at which node group was created.",
			},
		},
		Description: "Provides a cluster node group resource for V2 API. This can be used to create, modify and delete cluster's node group.",
	}
}

func resourceKubernetesNodeGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// creates client for managed-k8s
	config := meta.(clients.Config)
	containerInfraClientV2, err := config.ContainerInfraV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS container infra v2 client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)

	// build VM engine
	vmEngine := nodegroups.VMEngine{
		NovaEngine: nodegroups.NovaEngine{
			FlavorID: d.Get("node_flavor").(string),
		},
	}

	// build scale specification
	scaleType := d.Get("scale_type").(string)
	scaleSpec := nodegroups.ScaleSpec{}

	switch scaleType {
	case "fixed_scale":
		scaleSpec.FixedScale = &nodegroups.FixedScale{
			Size: d.Get("fixed_scale_node_count").(int),
		}
	case "auto_scale":
		scaleSpec.AutoScale = &nodegroups.AutoScale{
			MinSize: d.Get("auto_scale_min_size").(int),
			MaxSize: d.Get("auto_scale_max_size").(int),
			Size:    d.Get("auto_scale_node_count").(int),
		}
	}

	// build labels
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

	// build taints
	taints, err := getAsTaintsSlice(d, "taints")
	if err != nil {
		return diag.FromErr(err)
	}

	// build disk type
	diskTypeConfig := nodegroups.DiskType{
		CinderVolumeType: nodegroups.CinderVolumeType{
			Type: d.Get("disk_type").(string),
			Size: d.Get("disk_size").(int),
		},
	}

	// build AZs
	az := d.Get("availability_zone").(string)

	// build ng spec
	spec := nodegroups.NodeGroupSpec{
		Name:                 d.Get("name").(string),
		VMEngine:             vmEngine,
		Zones:                []string{az},
		ScaleSpec:            scaleSpec,
		Labels:               labels,
		Taints:               taints,
		ParallelUpgradeChunk: d.Get("parallel_upgrade_chunk").(int),
		DiskType:             diskTypeConfig,
	}

	// create node group
	createOpts := nodegroups.CreateOpts{
		ClusterID: clusterID,
		Spec:      spec,
	}

	nodeGroupID, err := nodegroups.Create(containerInfraClientV2, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating VKCS Kubernetes Node Group V2: %s", err)
	}

	d.SetId(nodeGroupID)

	// wait for node group to become active again after update
	stateConf := &retry.StateChangeConf{
		Pending:      []string{clusterStatusV2Provisioning, clusterStatusV2Reconciling},
		Target:       []string{clusterStatusV2Running},
		Refresh:      kubernetesStateRefreshFuncV2(containerInfraClientV2, d.Get("cluster_id").(string)),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        createNgDelayV2 * time.Minute,
		PollInterval: createNgPollIntervalV2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for vkcs_kubernetes_node_group_v2 %s to become ready after update: %s", d.Id(), err)
	}

	return resourceKubernetesNodeGroupV2Read(ctx, d, meta)
}

func resourceKubernetesNodeGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// creates client for managed-k8s
	config := meta.(clients.Config)
	containerInfraClientV2, err := config.ContainerInfraV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS container infra v2 client: %s", err)
	}

	nodeGroupID := d.Id()

	ng, err := nodegroups.Get(containerInfraClientV2, nodeGroupID).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error retrieving VKCS Kubernetes Node Group V2"))
	}

	d.Set("cluster_id", ng.ClusterID)
	d.Set("uuid", ng.UUID)
	d.Set("created_at", ng.CreatedAt)
	d.Set("name", ng.Name)

	if len(ng.Zones) < 1 {
		return diag.FromErr(fmt.Errorf("expected not empty list of availability zones for node_group %s", nodeGroupID))
	}
	d.Set("availability_zone", ng.Zones[0])

	d.Set("labels", ng.Labels)
	d.Set("parallel_upgrade_chunk", ng.ParallelUpgradeChunk)
	d.Set("node_flavor", ng.VMEngine.NovaEngine.FlavorID)

	// set scale spec
	if ng.ScaleSpec.FixedScale != nil {
		d.Set("scale_type", "fixed_scale")
		d.Set("fixed_scale_node_count", ng.ScaleSpec.FixedScale.Size)
	} else if ng.ScaleSpec.AutoScale != nil {
		d.Set("scale_type", "auto_scale")
		d.Set("auto_scale_min_size", ng.ScaleSpec.AutoScale.MinSize)
		d.Set("auto_scale_max_size", ng.ScaleSpec.AutoScale.MaxSize)
		d.Set("auto_scale_node_count", ng.ScaleSpec.AutoScale.Size)
	}

	// Set Taints
	if len(ng.Taints) > 0 {
		taints := make([]map[string]any, len(ng.Taints))
		for i, taint := range ng.Taints {
			taints[i] = map[string]any{
				"key":    taint.Key,
				"value":  taint.Value,
				"effect": taint.Effect,
			}
		}
		d.Set("taints", taints)
	} else {
		d.Set("taints", nil)
	}

	// set disk type
	d.Set("disk_type", ng.DiskType.CinderVolumeType.Type)
	d.Set("disk_size", ng.DiskType.CinderVolumeType.Size)

	return nil
}

func resourceKubernetesNodeGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// creates client for managed-k8s
	config := meta.(clients.Config)
	containerInfraClientV2, err := config.ContainerInfraV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra v2 client: %s", err)
	}

	// check which fields can be updated
	hasChanges := false
	spec := nodegroups.UpdateOpts{}

	// check for node_flavor update
	if d.HasChange("node_flavor") {
		hasChanges = true
		flavorID := d.Get("node_flavor").(string)
		spec.VMEngine = &nodegroups.VMEngine{
			NovaEngine: nodegroups.NovaEngine{
				FlavorID: flavorID,
			},
		}
	}

	// check for scale_type
	if d.HasChange("fixed_scale_node_count") || d.HasChange("auto_scale_node_count") || d.HasChange("auto_scale_min_size") || d.HasChange("auto_scale_max_size") {
		hasChanges = true
		scaleType := d.Get("scale_type").(string)
		scaleSpec := nodegroups.ScaleSpec{}

		switch scaleType {
		case "fixed_scale":
			fixedSize := d.Get("fixed_scale_node_count").(int)
			scaleSpec.FixedScale = &nodegroups.FixedScale{
				Size: fixedSize,
			}
		case "auto_scale":
			minSize := d.Get("auto_scale_min_size").(int)
			maxSize := d.Get("auto_scale_max_size").(int)
			size := d.Get("auto_scale_node_count").(int)
			scaleSpec.AutoScale = &nodegroups.AutoScale{
				MinSize: minSize,
				MaxSize: maxSize,
				Size:    size,
			}

		}
		spec.ScaleSpec = &scaleSpec
	}

	if d.HasChange("labels") {
		// build labels
		hasChanges = true
		rawLabels := d.Get("labels").(map[string]any)
		labels, err := extractKubernetesLabelsMap(rawLabels)
		if err != nil {
			return diag.FromErr(err)
		}
		spec.Labels = labels
	}

	if d.HasChange("taints") {
		// build taints
		hasChanges = true
		taints, err := getAsTaintsSlice(d, "taints")
		if err != nil {
			return diag.FromErr(err)
		}
		spec.Taints = taints
	}

	if d.HasChange("parallel_upgrade_chunk") {
		// build taints
		hasChanges = true
		parallelUpgradeChunk := d.Get("parallel_upgrade_chunk").(int)
		spec.ParallelUpgradeChunk = &parallelUpgradeChunk
	}

	if !hasChanges {
		log.Printf("[DEBUG] No changes detected for vkcs_kubernetes_node_group_v2 %s", d.Id())
		return resourceKubernetesNodeGroupV2Read(ctx, d, meta)
	}

	log.Printf("[DEBUG] Updating vkcs_kubernetes_node_group_v2: %#v", spec)

	err = nodegroups.Scale(containerInfraClientV2, d.Id(), spec)
	if err != nil {
		return diag.Errorf("error updating vkcs_kubernetes_node_group_v2: %s", err)
	}

	// wait for node group to become active again after update
	stateConf := &retry.StateChangeConf{
		Pending:      []string{clusterStatusV2Reconciling},
		Target:       []string{clusterStatusV2Running},
		Refresh:      kubernetesStateRefreshFuncV2(containerInfraClientV2, d.Get("cluster_id").(string)),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        updateNgDelayV2 * time.Minute,
		PollInterval: updateNgPollIntervalV2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for vkcs_kubernetes_node_group_v2 %s to become ready after update: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Updated vkcs_kubernetes_node_group_v2 %s", d.Id())
	return resourceKubernetesNodeGroupV2Read(ctx, d, meta)
}

func resourceKubernetesNodeGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// creates client for managed-k8s
	config := meta.(clients.Config)
	containerInfraClientV2, err := config.ContainerInfraV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra v2 client: %s", err)
	}

	// delete node group
	err = nodegroups.Delete(containerInfraClientV2, d.Id())
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "error deleting vkcs_kubernetes_node_group_v2"))
	}

	log.Printf("[DEBUG] Deleted vkcs_kubernetes_node_group_v2 %s", d.Id())
	return nil
}

func customizeNodeGroupV2Diff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	scaleType := d.Get("scale_type").(string)

	switch scaleType {
	case "fixed_scale":
		// Проверяем, что поля auto_scale не заданы
		if _, ok := d.GetOk("auto_scale_min_size"); ok {
			return fmt.Errorf("auto_scale_min_size cannot be set when scale_type is 'fixed_scale'")
		}
		if _, ok := d.GetOk("auto_scale_max_size"); ok {
			return fmt.Errorf("auto_scale_max_size cannot be set when scale_type is 'fixed_scale'")
		}
		if _, ok := d.GetOk("auto_scale_node_count"); ok {
			return fmt.Errorf("auto_scale_node_count cannot be set when scale_type is 'fixed_scale'")
		}

		// Проверяем, что fixed_scale_node_count задано
		if _, ok := d.GetOk("fixed_scale_node_count"); !ok {
			return fmt.Errorf("fixed_scale_node_count is required when scale_type is 'fixed_scale'")
		}
	case "auto_scale":
		// Проверяем, что fixed_scale_node_count не задано
		if _, ok := d.GetOk("fixed_scale_node_count"); ok {
			return fmt.Errorf("fixed_scale_node_count cannot be set when scale_type is 'auto_scale'")
		}

		// Проверяем, что все auto_scale поля заданы
		minSize, minOk := d.GetOk("auto_scale_min_size")
		maxSize, maxOk := d.GetOk("auto_scale_max_size")
		nodeCount, countOk := d.GetOk("auto_scale_node_count")

		if !minOk || !maxOk || !countOk {
			return fmt.Errorf("auto_scale_min_size, auto_scale_max_size, and auto_scale_node_count are required when scale_type is 'auto_scale'")
		}

		// Проверяем условие min <= node_count <= max
		if minSize.(int) > nodeCount.(int) || nodeCount.(int) > maxSize.(int) {
			return fmt.Errorf("for auto_scale, condition 'auto_scale_min_size <= auto_scale_node_count <= auto_scale_max_size' must be met")
		}
	}

	// TODO: добавить валидацию key, values для labels и taints
	return nil
}
