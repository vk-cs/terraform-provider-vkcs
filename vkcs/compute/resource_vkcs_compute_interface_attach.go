package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
)

func ResourceComputeInterfaceAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeInterfaceAttachCreate,
		ReadContext:   resourceComputeInterfaceAttachRead,
		DeleteContext: resourceComputeInterfaceAttachDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to create the interface attachment. If omitted, the `region` argument of the provider is used. Changing this creates a new attachment.",
			},

			"port_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"network_id"},
				Description: "The ID of the Port to attach to an Instance.\n" +
					"_NOTE_: This option and `network_id` are mutually exclusive.",
			},

			"network_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"port_id"},
				Description: "The ID of the Network to attach to an Instance. A port will be created automatically.\n" +
					"_NOTE_: This option and `port_id` are mutually exclusive.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the Instance to attach the Port or Network to.",
			},

			"fixed_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"port_id"},
				Description: "An IP address to assosciate with the port.\n" +
					"_NOTE_: This option cannot be used with port_id. You must specify a network_id. The IP address must lie in a range on the supplied network.",
			},
		},
		Description: "Attaches a Network Interface (a Port) to an Instance using the VKCS Compute API.",
	}
}

func resourceComputeInterfaceAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	var portID string
	if v, ok := d.GetOk("port_id"); ok {
		portID = v.(string)
	}

	var networkID string
	if v, ok := d.GetOk("network_id"); ok {
		networkID = v.(string)
	}

	if networkID == "" && portID == "" {
		return diag.Errorf("Must set one of network_id and port_id")
	}

	// For some odd reason the API takes an array of IPs, but you can only have one element in the array.
	var fixedIPs []attachinterfaces.FixedIP
	if v, ok := d.GetOk("fixed_ip"); ok {
		fixedIPs = append(fixedIPs, attachinterfaces.FixedIP{IPAddress: v.(string)})
	}

	attachOpts := attachinterfaces.CreateOpts{
		PortID:    portID,
		NetworkID: networkID,
		FixedIPs:  fixedIPs,
	}

	log.Printf("[DEBUG] vkcs_compute_interface_attach attach options: %#v", attachOpts)

	attachment, err := attachinterfaces.Create(computeClient, instanceID, attachOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ATTACHING"},
		Target:     []string{"ATTACHED"},
		Refresh:    computeInterfaceAttachAttachFunc(computeClient, instanceID, attachment.PortID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error creating vkcs_compute_interface_attach %s: %s", instanceID, err)
	}

	// Use the instance ID and attachment ID as the resource ID.
	id := fmt.Sprintf("%s/%s", instanceID, attachment.PortID)

	log.Printf("[DEBUG] Created vkcs_compute_interface_attach %s: %#v", id, attachment)

	d.SetId(id)

	return resourceComputeInterfaceAttachRead(ctx, d, meta)
}

func resourceComputeInterfaceAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	instanceID, attachmentID, err := ComputeInterfaceAttachParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	attachment, err := attachinterfaces.Get(computeClient, instanceID, attachmentID).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error retrieving vkcs_compute_interface_attach"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_compute_interface_attach %s: %#v", d.Id(), attachment)

	if len(attachment.FixedIPs) > 0 {
		d.Set("fixed_ip", attachment.FixedIPs[0].IPAddress)
	}

	d.Set("instance_id", instanceID)
	d.Set("port_id", attachment.PortID)
	d.Set("network_id", attachment.NetID)
	d.Set("region", util.GetRegion(d, config))

	return nil
}

func resourceComputeInterfaceAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	instanceID, attachmentID, err := ComputeInterfaceAttachParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{"DETACHED"},
		Refresh:    computeInterfaceAttachDetachFunc(computeClient, instanceID, attachmentID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error detaching vkcs_compute_interface_attach %s: %s", d.Id(), err)
	}

	return nil
}
