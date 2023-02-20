package vkcs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func resourceMonitoringTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTriggerCreate,
		ReadContext:   resourceTriggerRead,
		UpdateContext: resourceTriggerUpdate,
		DeleteContext: resourceTriggerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Human-readable name for the trigger.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "on",
				Description: "Human-readable status for the trigger.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "mcs/vm",
				Description: "Namespace for metrics tenant. For vm mcs/vm",
			},
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Promql query for triggers like: cpu_usage_guest{vm_uuid=\"03e042d3-6b68-47eb-ac35-8b4f6c2f77e2\"} > 50",
			},
			"interval": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Run interval in seconds.",
			},
			"notification_title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Template for notification messages. Like: {{ $labels.host }}  {{$labels.__name__}} = {{ $value }}",
			},
			"notification_channels": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of channel uuids.",
			},
		},
		Description: "Manages monitoring triggers within VKCS.",
	}
}

func resourceTriggerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	MonitoringV1Client, err := config.MonitoringV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS monitoring client: %s", err)
	}

	var channels []string
	for _, ic := range d.Get("notification_channels").([]interface{}) {
		channels = append(channels, ic.(string))
	}

	tr := TriggerIn{
		Trigger: CreateTrigger{
			Name:                 d.Get("name").(string),
			Status:               d.Get("status").(string),
			Namespace:            d.Get("namespace").(string),
			Query:                d.Get("query").(string),
			Interval:             d.Get("interval").(int),
			NotificationTitle:    d.Get("notification_title").(string),
			NotificationChannels: channels,
		},
	}

	ch, err := triggerCreate(MonitoringV1Client, config.GetTenantID(), &tr).extract()
	if err != nil {
		return diag.Errorf("Error creating VKCS monitoring trigger: %s", err)
	}
	d.SetId(ch.Trigger.Id)

	return nil

}

func resourceTriggerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	MonitoringV1Client, err := config.MonitoringV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS monitoring client: %s", err)
	}
	ch, err := triggerGet(MonitoringV1Client, config.GetTenantID(), d.Id()).extract()
	if err != nil {
		return diag.Errorf("Error get VKCS monitoring trigger(%s): %s", d.Id(), err)
	}
	d.Set("name", ch.Trigger.Name)
	d.Set("status", ch.Trigger.Status)
	d.Set("query", ch.Trigger.Query)
	//d.Set("interval", ch.Trigger.Interval)
	d.Set("notification_title", ch.Trigger.NotificationTitle)
	d.Set("notification_channels", ch.Trigger.NotificationChannels)
	return nil
}

func resourceTriggerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	MonitoringV1Client, err := config.MonitoringV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS monitoring client: %s", err)
	}

	var channels []string
	for _, ic := range d.Get("notification_channels").([]interface{}) {
		channels = append(channels, ic.(string))
	}

	tr := TriggerIn{
		Trigger: CreateTrigger{
			Name:                 d.Get("name").(string),
			Status:               d.Get("status").(string),
			Namespace:            d.Get("namespace").(string),
			Query:                d.Get("query").(string),
			Interval:             d.Get("interval").(int),
			NotificationTitle:    d.Get("notification_title").(string),
			NotificationChannels: channels,
		},
	}
	_, err = triggerUpdate(MonitoringV1Client, config.GetTenantID(), d.Id(), &tr).extract()
	if err != nil {
		return diag.Errorf("Error update VKCS monitoring trigger: %s", err)
	}

	return nil
}

func resourceTriggerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	MonitoringV1Client, err := config.MonitoringV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS monitoring client: %s", err)
	}

	err = triggerDelete(MonitoringV1Client, config.GetTenantID(), d.Id()).extractErr()
	if err != nil {
		return diag.Errorf("Error de VKCS monitoring trigger: %s", err)
	}
	return nil
}
