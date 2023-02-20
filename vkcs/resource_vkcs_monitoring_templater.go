package vkcs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func resourceMonitoringTemplater() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTemplaterCreate,
		ReadContext:   resourceTemplaterRead,
		UpdateContext: resourceTemplaterUpdate,
		DeleteContext: resourceTemplaterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Compute instance id.",
			},
			"script": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Script for vm. It's contains script for agent installation.",
			},
		},
		Description: "Manages monitoring template within (for compute instances) VKCS.",
	}
}

func resourceTemplaterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	MonitoringTemplaterV2Client, err := config.MonitoringTemplaterV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS monitoring client: %s", err)
	}

	tmp := TemplateIn{
		InstanceId:   d.Get("instance_id").(string),
		Capabilities: []string{"telegraf"},
	}

	t, err := templateCreate(MonitoringTemplaterV2Client, config.GetTenantID(), &tmp).extract()
	if err != nil {
		return diag.Errorf("Error creating VKCS monitoring template: %s", err)
	}
	d.SetId(t.LinkId)
	d.Set("script", t.Script)
	return nil
}

func resourceTemplaterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceTemplaterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceTemplaterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
