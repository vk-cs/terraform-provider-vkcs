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
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
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
