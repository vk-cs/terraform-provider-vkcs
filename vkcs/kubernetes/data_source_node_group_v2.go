package kubernetes

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func DataSourceKubernetesNodeGroupV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesNodeGroupV2Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Region to use for the node group. Default is a region configured for provider.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the node group.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the cluster.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the node group.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of node group.",
			},
			"node_flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Flavor ID of the nodes from node group.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The availability zone of the node group.",
			},
			"scale_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of scaling for the node group.",
			},
			"fixed_scale_node_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The node count of the node group, if scale_type is 'fixed_scale'.",
			},
			"auto_scale_node_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The node count of the node group, if scale_type is 'auto_scale'.",
			},
			"auto_scale_min_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum allowed nodes for this node group.",
			},
			"auto_scale_max_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum allowed nodes for this node group.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of key-value pairs representing additional properties of the node group.",
			},
			"taints": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key of the taint.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of the taint.",
						},
						"effect": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Effect of the taint.",
						},
					},
				},
				Description: "The list of objects representing node group taints.",
			},
			"parallel_upgrade_chunk": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum percent of nodes that can be unavailable during an upgrade.",
			},
			"disk_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The volume type to load nodes from.",
			},
			"disk_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size in GB for volume to load nodes from.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time at which node group was created.",
			},
		},
		Description: "Use this data source to get information on VKCS Kubernetes cluster's node group (V2 API).",
	}
}

func dataSourceKubernetesNodeGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// creates client for managed-k8s
	config := meta.(clients.Config)
	containerInfraClientV2, err := config.ContainerInfraV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS container infra v2 client: %s", err)
	}

	ngID := d.Get("id").(string)
	if ngID == "" {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "id must be specified",
			AttributePath: cty.Path{
				cty.GetAttrStep{Name: "id"},
			},
		}}
	}
	d.SetId(ngID)

	ng, err := nodegroups.Get(containerInfraClientV2, ngID).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error retrieving VKCS Kubernetes Node Group V2"))
	}

	d.Set("cluster_id", ng.ClusterID)
	d.Set("uuid", ng.UUID)
	d.Set("created_at", ng.CreatedAt)
	d.Set("name", ng.Name)

	if len(ng.Zones) < 1 {
		return diag.FromErr(fmt.Errorf("expected not empty list of availability zones for node_group %s", ngID))
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
	d.Set("region", util.GetRegion(d, config))

	return nil
}
