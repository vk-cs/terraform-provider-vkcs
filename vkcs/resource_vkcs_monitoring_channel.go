package vkcs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func resourceMonitoringChannel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChannelCreate,
		ReadContext:   resourceChannelRead,
		UpdateContext: resourceChannelUpdate,
		DeleteContext: resourceChannelDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Human-readable name for the channel.",
			},
			"channel_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of channel: email or sms.",
			},
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Address for channel email or phone.",
			},
		},
		Description: "Manages monitoring notification channels within VKCS.",
	}
}

func resourceChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config := meta.(configer)
	MonitoringV1Client, err := config.MonitoringV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS monitoring client: %s", err)
	}

	chn := ChannelIn{
		Name:        d.Get("name").(string),
		ChannelType: d.Get("channel_type").(string),
		Address:     d.Get("address").(string),
	}

	ch, err := channelCreate(MonitoringV1Client, config.GetTenantID(), &chn).extract()
	if err != nil {
		return diag.Errorf("Error creating VKCS monitoring channel: %s", err)
	}
	d.SetId(ch.Channel.ID)

	return nil
}

func resourceChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	MonitoringV1Client, err := config.MonitoringV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS monitoring client: %s", err)
	}
	ch, err := channelGet(MonitoringV1Client, config.GetTenantID(), d.Id()).extract()
	if err != nil {
		return diag.Errorf("Error get VKCS monitoring channel(%s): %s", d.Id(), err)
	}
	d.Set("name", ch.Channel.Name)
	d.Set("channel_type", ch.Channel.ChannelType)
	d.Set("address", ch.Channel.Address)
	return nil
}

func resourceChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	MonitoringV1Client, err := config.MonitoringV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS monitoring client: %s", err)
	}

	chn := ChannelIn{
		Name:        d.Get("name").(string),
		ChannelType: d.Get("channel_type").(string),
		Address:     d.Get("address").(string),
	}
	_, err = channelUpdate(MonitoringV1Client, config.GetTenantID(), d.Id(), &chn).extract()
	if err != nil {
		return diag.Errorf("Error update VKCS monitoring channel: %s", err)
	}

	return nil
}

func resourceChannelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	MonitoringV1Client, err := config.MonitoringV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS monitoring client: %s", err)
	}

	err = channelDelete(MonitoringV1Client, config.GetTenantID(), d.Id()).extractErr()
	if err != nil {
		return diag.Errorf("Error de VKCS monitoring channel: %s", err)
	}

	return nil
}
