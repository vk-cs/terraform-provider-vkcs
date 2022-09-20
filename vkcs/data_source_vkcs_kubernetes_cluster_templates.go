package vkcs

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKubernetesClusterTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVkcsClusterTemplatesRead,
		Schema: map[string]*schema.Schema{
			"cluster_templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_template_uuid": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The UUID of the cluster template.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the cluster template.",
						},
						"version": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The version of the cluster template.",
						},
					},
				},
				// Currently, Computed composite fields are treated by terraform-schema as "Objects". Thus, subfields and their descriptions are not available to documentation generator. We have to put these descriptions into parent field description as workaround.
				Description: "A list of available kubernetes cluster templates.\n  - `cluster_template_uuid` **String** The UUID of the cluster template.\n\n  - `name` **String** The name of the cluster template.\n\n  - `version` **String** The version of the cluster template.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Random identifier of the data source.",
			},
		},
		Description: "`vkcs_kubernetes_cluster_templates` returns list of available VKCS Kubernetes Cluster Templates. To get details of each cluster template the data source can be combined with the `vkcs_kubernetes_clustertemplate` data source.",
	}
}

type clusterTemplateResponse struct {
	Version string
	UUID    string
	Name    string
}

type clusterTemplateFlatSchema []map[string]interface{}

func flattenClusterTemplates(templates []clusterTemplateResponse) clusterTemplateFlatSchema {
	flatSchema := clusterTemplateFlatSchema{}
	for _, template := range templates {
		flatSchema = append(flatSchema, map[string]interface{}{
			"name":                  template.Name,
			"cluster_template_uuid": template.UUID,
			"version":               template.Version,
		})
	}
	return flatSchema
}

func dataSourceVkcsClusterTemplatesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	client, err := config.ContainerInfraV1Client(config.GetRegion())
	if err != nil {
		return diag.Errorf("failed to init identity v3 client: %s", err)
	}

	templates, err := clusterTemplateList(client).Extract()
	if err != nil {
		return diag.Errorf("failed to list cluster templates: %s", err)
	}

	clusterTemplates := make([]clusterTemplateResponse, 0, len(templates))
	for _, t := range templates {
		clusterTemplates = append(clusterTemplates, clusterTemplateResponse{
			UUID:    t.UUID,
			Name:    t.Name,
			Version: t.Version,
		})
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	if err := d.Set("cluster_templates", flattenClusterTemplates(clusterTemplates)); err != nil {
		return diag.Errorf("failed to set cluster templates: %s", err)
	}

	return nil
}
