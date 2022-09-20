package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKubernetesNodeGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesNodeGroupRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of cluster that node group belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the node group.",
			},
			"node_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    false,
				Description: "The count of nodes in node group.",
			},
			"max_nodes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    false,
				Description: "The maximum amount of nodes in node group.",
			},
			"min_nodes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    false,
				Description: "The minimum amount of nodes in node group.",
			},
			"volume_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    false,
				Description: "The amount of memory of volume in Gb.",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				Description: "The type of volume.",
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				Description: "The id of flavor.",
			},
			"autoscaling_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				Description: "Determines whether the autoscaling is enabled.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the cluster's node group.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Determines current state of node group (RUNNING, SHUTOFF, ERROR).",
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UUID of node.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the node.",
						},
						"node_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The node group id",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time when node was created.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time when node was updated.",
						},
					},
				},
				Description: "The list of node group's node objects.",
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of availability zones of the node group.",
			},
			"max_node_unavailable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specified as a percentage. The maximum number of nodes that can fail during an upgrade.",
			},
		},
		Description: "Use this data source to get the ID of an available VKCS kubernetes clusters node group.",
	}
}

func dataSourceKubernetesNodeGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	containerInfraClient, err := config.ContainerInfraV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating container infra client: %s", err)
	}

	nodeGroup, err := nodeGroupGet(containerInfraClient, d.Get("uuid").(string)).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "error retrieving vkcs_kubernetes_node_group"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_kubernetes_node_group %s: %#v", d.Id(), nodeGroup)

	d.SetId(nodeGroup.UUID)
	d.Set("cluster_id", nodeGroup.ClusterID)
	d.Set("name", nodeGroup.Name)
	d.Set("node_count", nodeGroup.NodeCount)
	d.Set("max_nodes", nodeGroup.MaxNodes)
	d.Set("min_nodes", nodeGroup.MinNodes)
	d.Set("volume_size", nodeGroup.VolumeSize)
	d.Set("volume_type", nodeGroup.VolumeType)
	d.Set("flavor_id", nodeGroup.FlavorID)
	d.Set("autoscaling_enabled", nodeGroup.Autoscaling)
	d.Set("nodes", flattenNodes(nodeGroup.Nodes))
	d.Set("state", nodeGroup.State)
	d.Set("availability_zones", nodeGroup.AvailabilityZones)
	d.Set("max_node_unavailable", nodeGroup.MaxNodeUnavailable)

	if err := d.Set("created_at", getTimestamp(&nodeGroup.CreatedAt)); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_kubernetes_node_group created_at: %s", err)
	}
	if err := d.Set("updated_at", getTimestamp(&nodeGroup.UpdatedAt)); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_kubernetes_node_group updated_at: %s", err)
	}

	return nil
}
