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
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: false,
			},
			"max_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: false,
			},
			"min_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: false,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: false,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: false,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: false,
			},
			"autoscaling_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: false,
			},
			"uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_group_id": {
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
					},
				},
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
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

	if err := d.Set("created_at", getTimestamp(&nodeGroup.CreatedAt)); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_kubernetes_node_group created_at: %s", err)
	}
	if err := d.Set("updated_at", getTimestamp(&nodeGroup.UpdatedAt)); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_kubernetes_node_group updated_at: %s", err)
	}

	return nil
}
