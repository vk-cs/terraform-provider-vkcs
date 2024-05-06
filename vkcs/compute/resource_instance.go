package compute

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/availabilityzones"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/schedulerhints"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/shelveunshelve"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/tags"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	flavorsutils "github.com/gophercloud/utils/openstack/compute/v2/flavors"
	imagesutils "github.com/gophercloud/utils/openstack/imageservice/v2/images"
	"github.com/gophercloud/utils/terraform/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	ivolumes "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/blockstorage/v3/volumes"
	ibootfromvolume "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/bootfromvolume"
	iflavors "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/flavors"
	iimages "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/images"
	isecgroups "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/secgroups"
	iservers "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/servers"
	ishelveunshelve "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/shelveunshelve"
	istartstop "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/startstop"
	itags "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/tags"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func ResourceComputeInstance() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: resourceComputeInstanceCustomizeDiff,

		CreateContext: resourceComputeInstanceCreate,
		ReadContext:   resourceComputeInstanceRead,
		UpdateContext: resourceComputeInstanceUpdate,
		DeleteContext: resourceComputeInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceComputeInstanceImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to create the server instance. If omitted, the `region` argument of the provider is used. Changing this creates a new server.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "A unique name for the resource.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The image ID of the desired image for the server. Required if `image_name` is empty and not booting from a volume. Do not specify if booting from a volume. Changing this creates a new server.",
			},
			"image_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The name of the desired image for the server. Required if `image_id` is empty and not booting from a volume. Do not specify if booting from a volume. Changing this creates a new server.",
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				Description: "The flavor ID of the desired flavor for the server. Required if `flavor_name` is empty. Changing this resizes the existing server.",
			},
			"flavor_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				Description: "The name of the desired flavor for the server. Required if `flavor_id` is empty. Changing this resizes the existing server.",
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// just stash the hash for state & diff comparisons
				StateFunc: func(v interface{}) string {
					switch v := v.(type) {
					case string:
						hash := sha1.Sum([]byte(v))
						return hex.EncodeToString(hash[:])
					default:
						return ""
					}
				},
				Description: "The user data to provide when launching the instance.	Changing this creates a new server.",
			},
			"security_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "An array of one or more security group names to associate with the server. Changing this results in adding/removing security groups from the existing server. _note_ When attaching the instance to networks using Ports, place the security groups on the Port and not the instance. _note_ Names should be used and not ids, as ids trigger unnecessary updates.",
			},
			"availability_zone": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: suppressAvailabilityZoneDetailDiffs,
				Description:      "The availability zone in which to create the server. Conflicts with `availability_zone_hints`. Changing this creates a new server.",
			},
			"network_mode": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      false,
				ConflictsWith: []string{"network"},
				ValidateFunc: validation.StringInSlice([]string{
					"auto", "none",
				}, true),
				Description: "Special string for `network` option to create the server. `network_mode` can be `\"auto\"` or `\"none\"`. Please see the following [reference](https://docs.openstack.org/api-ref/compute/?expanded=create-server-detail#id11) for more information. Conflicts with `network`.",
			},
			"network": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "The network UUID to attach to the server. Optional if `port` or `name` is provided. Changing this creates a new server.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The human-readable name of the network. Optional if `uuid` or `port` is provided. Changing this creates a new server.",
							Computed:    true,
						},
						"port": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "The port UUID of a network to attach to the server. Optional if `uuid` or `name` is provided. Changing this creates a new server. _note_ If port is used, only its security groups will be applied instead of security_groups instance argument.",
						},
						"fixed_ip_v4": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "Specifies a fixed IPv4 address to be used on this network. Changing this creates a new server.",
						},
						"mac": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MAC address of the NIC on that network.",
						},
						"access_network": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Specifies if this network should be used for provisioning access. Accepts true or false. Defaults to false.",
						},
					},
				},
				Description: "An array of one or more networks to attach to the instance. The network object structure is documented below. Changing this creates a new server.",
			},
			"metadata": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    false,
				Description: "Metadata key/value pairs to make available from within the instance. Changing this updates the existing server metadata.",
			},
			"config_drive": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to use the config_drive feature to configure the instance. Changing this creates a new server.",
			},
			"admin_pass": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    false,
				Description: "The administrative password to assign to the server. Changing this changes the root password on the existing server.",
			},
			"access_ip_v4": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    false,
				Description: "The first detected Fixed IPv4 address.",
			},
			"key_pair": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of a key pair to put on the server. The key pair must already be created and associated with the tenant's account. Changing this creates a new server.",
			},
			"block_device": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The source type of the device. Must be one of \"blank\", \"image\", \"volume\", or \"snapshot\". Changing this creates a new server.",
						},
						"uuid": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The UUID of the image, volume, or snapshot. Optional if `source_type` is set to `\"blank\"`. Changing this creates a new server.",
						},
						"volume_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "The size of the volume to create (in gigabytes). Required in the following combinations: source=image and destination=volume, source=blank and destination=local, and source=blank and destination=volume. Changing this creates a new server.",
						},
						"destination_type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The type that gets created. Possible values are \"volume\" and \"local\". Changing this creates a new server.",
						},
						"boot_index": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "The boot index of the volume. It defaults to 0 if only one `block_device` is specified, and to -1 if more than one is configured. Changing this creates a new server. _note_ You must set the boot index to 0 for one of the block devices if more than one is defined.",
						},
						"delete_on_termination": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
							Description: "Delete the volume / block device upon termination of the instance. Defaults to false. Changing this creates a new server. __note__ It is important to enable `delete_on_termination` for volumes created with instance. If `delete_on_termination` is disabled for such volumes, then after instance deletion such volumes will stay orphaned and uncontrolled by terraform. __note__ It is important to disable `delete_on_termination` if volume is created as separate terraform resource and is attached to instance. Enabling `delete_on_termination` for such volumes will result in mismanagement between two terraform resources in case of instance deletion",
						},
						"guest_format": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Specifies the guest server disk file system format, such as `ext2`, `ext3`, `ext4`, `xfs` or `swap`. Swap block device mappings have the following restrictions: source_type must be blank and destination_type must be local and only one swap disk per server and the size of the swap disk must be less than or equal to the swap size of the flavor. Changing this creates a new server.",
						},
						"volume_type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The volume type that will be used. Changing this creates a new server.",
						},
						"device_type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The low-level device type that will be used. Most common thing is to leave this empty. Changing this creates a new server.",
						},
						"disk_bus": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The low-level disk bus that will be used. Most common thing is to leave this empty. Changing this creates a new server.",
						},
					},
				},
				Description: "Configuration of block devices. The block_device structure is documented below. Changing this creates a new server. You can specify multiple block devices which will create an instance with multiple disks. This configuration is very flexible, so please see the following [reference](https://docs.openstack.org/nova/latest/user/block-device-mapping.html) for more information.",
			},
			"scheduler_hints": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "A UUID of a Server Group. The instance will be placed into that group.",
						},
					},
				},
				Set:         resourceComputeSchedulerHintsHash,
				Description: "Provide the Nova scheduler with hints on how the instance should be launched. The available hints are described below.",
			},
			"personality": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The absolute path of the destination file.",
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The contents of the file. Limited to 255 bytes.",
						},
					},
				},
				Set:         resourceComputeInstancePersonalityHash,
				Description: "Customize the personality of an instance by defining one or more files and their contents. The personality structure is described below. _note_ 'config_drive' must be enabled.",
			},
			"stop_before_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to try stop instance gracefully before destroying it, thus giving chance for guest OS daemons to stop correctly. If instance doesn't stop within timeout, it will be destroyed anyway.",
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to force the compute instance to be forcefully deleted. This is useful for environments that have reclaim / soft deletion enabled.",
			},
			"all_metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Contains all instance metadata, even metadata not set by Terraform.",
			},
			"power_state": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{
					"active", "shutoff", "shelved_offloaded",
				}, true),
				DiffSuppressFunc: suppressPowerStateDiffs,
				Description:      "Provide the VM state. Only 'active' and 'shutoff' are supported values. _note_ If the initial power_state is the shutoff the VM will be stopped immediately after build and the provisioners like remote-exec or files are not supported.",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A set of string tags for the instance. Changing this updates the existing instance tags.",
			},
			"all_tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The collection of tags assigned on the instance, which have been explicitly and implicitly added.",
			},
			"vendor_options": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ignore_resize_confirmation": {
							Type:        schema.TypeBool,
							Default:     false,
							Optional:    true,
							Description: "Boolean to control whether to ignore manual confirmation of the instance resizing.",
						},
						"detach_ports_before_destroy": {
							Type:        schema.TypeBool,
							Default:     false,
							Optional:    true,
							Description: "Whether to try to detach all attached ports to the vm before destroying it to make sure the port state is correct after the vm destruction. This is helpful when the port is not deleted.",
						},
					},
				},
				Description: "Map of additional vendor-specific options. Supported options are described below.",
			},
		},
		Description: "Manages a compute VM instance resource._note_ All arguments including the instance admin password will be stored in the raw state as plain-text. [Read more about sensitive data in state](https://www.terraform.io/docs/language/state/sensitive-data.html).",
	}
}

func resourceComputeInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	imageClient, err := config.ImageV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS image client: %s", err)
	}

	computeClient.Microversion = computeAPIMicroVersion
	var createOpts servers.CreateOptsBuilder
	var availabilityZone string
	var networks interface{}
	var diags diag.Diagnostics

	// Determines the Image ID using the following rules:
	// If a bootable block_device was specified, ignore the image altogether.
	// If an image_id was specified, use it.
	// If an image_name was specified, look up the image ID, report if error.
	imageID, err := getImageIDFromConfig(imageClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Determines the Flavor ID using the following rules:
	// If a flavor_id was specified, use it.
	// If a flavor_name was specified, lookup the flavor ID, report if error.
	flavorID, err := getFlavorID(computeClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// determine if block_device configuration is correct
	// this includes valid combinations and required attributes
	diags = append(diags, checkBlockDeviceConfig(d)...)
	if diags.HasError() {
		return diags
	}

	if networkMode := d.Get("network_mode").(string); networkMode == "auto" || networkMode == "none" {
		// Use special string for network option
		networks = networkMode
		log.Printf("[DEBUG] Create with network options %s", networks)
	} else {
		log.Printf("[DEBUG] Create with specified network options")
		// Build a list of networks with the information given upon creation.
		// Error out if an invalid network configuration was used.
		allInstanceNetworks, err := getAllInstanceNetworks(d, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		// Build a []servers.Network to pass into the create options.
		networks = expandInstanceNetworks(allInstanceNetworks)
	}

	configDrive := d.Get("config_drive").(bool)

	// Retrieve tags and set microversion if they're provided.
	instanceTags := ComputeInstanceTags(d)

	if v, ok := d.GetOk("availability_zone"); ok {
		availabilityZone = v.(string)
	}

	createOpts = &servers.CreateOpts{
		Name:             d.Get("name").(string),
		ImageRef:         imageID,
		FlavorRef:        flavorID,
		SecurityGroups:   resourceInstanceSecGroupsV2(d),
		AvailabilityZone: availabilityZone,
		Networks:         networks,
		Metadata:         resourceInstanceMetadataV2(d),
		ConfigDrive:      &configDrive,
		AdminPass:        d.Get("admin_pass").(string),
		UserData:         []byte(d.Get("user_data").(string)),
		Personality:      resourceInstancePersonalityV2(d),
		Tags:             instanceTags,
	}

	if keyName, ok := d.Get("key_pair").(string); ok && keyName != "" {
		createOpts = &keypairs.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			KeyName:           keyName,
		}
	}

	if vL, ok := d.GetOk("block_device"); ok {
		blockDevices, err := ResourceInstanceBlockDevicesV2(d, vL.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}

		createOpts = &bootfromvolume.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			BlockDevice:       blockDevices,
		}
	}

	schedulerHintsRaw := d.Get("scheduler_hints").(*schema.Set).List()
	if len(schedulerHintsRaw) > 0 {
		log.Printf("[DEBUG] schedulerhints: %+v", schedulerHintsRaw)
		schedulerHints := resourceInstanceSchedulerHintsV2(d, schedulerHintsRaw[0].(map[string]interface{}))
		createOpts = &schedulerhints.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			SchedulerHints:    schedulerHints,
		}
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// If a block_device is used, use the bootfromvolume.Create function as it allows an empty ImageRef.
	// Otherwise, use the normal servers.Create function.
	var server *servers.Server
	if _, ok := d.GetOk("block_device"); ok {
		server, err = ibootfromvolume.Create(computeClient, createOpts).Extract()
	} else {
		server, err = iservers.Create(computeClient, createOpts).Extract()
	}

	if err != nil {
		return diag.Errorf("Error creating VKCS server: %s", err)
	}
	log.Printf("[INFO] Instance ID: %s", server.ID)

	// Store the ID now
	d.SetId(server.ID)

	// Wait for the instance to become running so we can get some attributes
	// that aren't available until later.
	log.Printf(
		"[DEBUG] Waiting for instance (%s) to become running",
		server.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    ServerStateRefreshFunc(computeClient, server.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	err = retry.RetryContext(ctx, stateConf.Timeout, func() *retry.RetryError {
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			log.Printf("[DEBUG] Retrying after error: %s", err)
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			server.ID, err)
	}

	vmState := d.Get("power_state").(string)
	if strings.ToLower(vmState) == "shutoff" {
		err = istartstop.Stop(computeClient, d.Id()).ExtractErr()
		if err != nil {
			return diag.Errorf("Error stopping VKCS instance: %s", err)
		}
		stopStateConf := &retry.StateChangeConf{
			// Pending:    []string{"ACTIVE"},
			Target:     []string{"SHUTOFF"},
			Refresh:    ServerStateRefreshFunc(computeClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		log.Printf("[DEBUG] Waiting for instance (%s) to stop", d.Id())
		_, err = stopStateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for instance (%s) to become inactive(shutoff): %s", d.Id(), err)
		}
	}

	diags = append(diags, resourceComputeInstanceRead(ctx, d, meta)...)
	return diags
}

func resourceComputeInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	server, err := iservers.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "server"))
	}

	log.Printf("[DEBUG] Retrieved Server %s: %+v", d.Id(), server)

	d.Set("name", server.Name)

	// Get the instance network and address information
	networks, err := flattenInstanceNetworks(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// Determine the best IPv4 addresses to access the instance with
	hostv4 := getInstanceAccessAddresses(d, networks)

	// AccessIPv4 isn't standard in OpenStack, but there have been reports
	// of them being used in some environments.
	if server.AccessIPv4 != "" && hostv4 == "" {
		hostv4 = server.AccessIPv4
	}

	d.Set("network", networks)
	d.Set("access_ip_v4", hostv4)

	// Determine the best IP address to use for SSH connectivity.
	// Prefer IPv4.
	var preferredSSHAddress string
	if hostv4 != "" {
		preferredSSHAddress = hostv4
	}

	if preferredSSHAddress != "" {
		// Initialize the connection info
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": preferredSSHAddress,
		})
	}

	d.Set("all_metadata", server.Metadata)

	secGrpNames := []string{}
	for _, sg := range server.SecurityGroups {
		secGrpNames = append(secGrpNames, sg["name"].(string))
	}
	d.Set("security_groups", secGrpNames)

	flavorID, ok := server.Flavor["id"].(string)
	if !ok {
		return diag.Errorf("Error setting VKCS server's flavor: %v", server.Flavor)
	}
	d.Set("flavor_id", flavorID)

	d.Set("key_pair", server.KeyName)
	flavor, err := iflavors.Get(computeClient, flavorID).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("flavor_name", flavor.Name)

	// Set the instance's image information appropriately
	if err := setImageInformation(computeClient, server, d); err != nil {
		return diag.FromErr(err)
	}

	// Build a custom struct for the availability zone extension
	var serverWithAZ struct {
		servers.Server
		availabilityzones.ServerAvailabilityZoneExt
	}

	// Do another Get so the above work is not disturbed.
	err = iservers.Get(computeClient, d.Id()).ExtractInto(&serverWithAZ)
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "server"))
	}
	// Set the availability zone
	d.Set("availability_zone", serverWithAZ.AvailabilityZone)

	// Set the region
	d.Set("region", util.GetRegion(d, config))

	// Set the current power_state
	currentStatus := strings.ToLower(server.Status)
	switch currentStatus {
	case "active", "shutoff", "error", "migrating", "shelved_offloaded", "shelved":
		d.Set("power_state", currentStatus)
	default:
		return diag.Errorf("Invalid power_state for instance %s: %s", d.Id(), server.Status)
	}

	// Populate tags.
	computeClient.Microversion = computeAPIMicroVersion
	instanceTags, err := itags.List(computeClient, server.ID).Extract()
	if err != nil {
		log.Printf("[DEBUG] Unable to get tags for vkcs_compute_instance: %s", err)
	} else {
		ComputeInstanceReadTags(d, instanceTags)
	}

	return nil
}

func resourceComputeInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	var updateOpts servers.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if updateOpts != (servers.UpdateOpts{}) {
		_, err := iservers.Update(computeClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating VKCS server: %s", err)
		}
	}

	if d.HasChange("power_state") {
		powerStateOldRaw, powerStateNewRaw := d.GetChange("power_state")
		powerStateOld := powerStateOldRaw.(string)
		powerStateNew := powerStateNewRaw.(string)
		if strings.ToLower(powerStateNew) == "shelved_offloaded" {
			err = ishelveunshelve.Shelve(computeClient, d.Id()).ExtractErr()
			if err != nil {
				return diag.Errorf("Error shelve VKCS instance: %s", err)
			}
			shelveStateConf := &retry.StateChangeConf{
				// Pending:    []string{"ACTIVE"},
				Target:     []string{"SHELVED_OFFLOADED"},
				Refresh:    ServerStateRefreshFunc(computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			log.Printf("[DEBUG] Waiting for instance (%s) to shelve", d.Id())
			_, err = shelveStateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to become shelve: %s", d.Id(), err)
			}
		}
		if strings.ToLower(powerStateNew) == "shutoff" {
			err = istartstop.Stop(computeClient, d.Id()).ExtractErr()
			if err != nil {
				return diag.Errorf("Error stopping VKCS instance: %s", err)
			}
			stopStateConf := &retry.StateChangeConf{
				// Pending:    []string{"ACTIVE"},
				Target:     []string{"SHUTOFF"},
				Refresh:    ServerStateRefreshFunc(computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			log.Printf("[DEBUG] Waiting for instance (%s) to stop", d.Id())
			_, err = stopStateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to become inactive(shutoff): %s", d.Id(), err)
			}
		}
		if strings.ToLower(powerStateNew) == "active" {
			if strings.ToLower(powerStateOld) == "shelved" || strings.ToLower(powerStateOld) == "shelved_offloaded" {
				unshelveOpt := &shelveunshelve.UnshelveOpts{
					AvailabilityZone: d.Get("availability_zone").(string),
				}
				err = ishelveunshelve.Unshelve(computeClient, d.Id(), unshelveOpt).ExtractErr()
				if err != nil {
					return diag.Errorf("Error unshelving VKCS instance: %s", err)
				}
			} else {
				err = istartstop.Start(computeClient, d.Id()).ExtractErr()
				if err != nil {
					return diag.Errorf("Error starting VKCS instance: %s", err)
				}
			}
			startStateConf := &retry.StateChangeConf{
				// Pending:    []string{"SHUTOFF"},
				Target:     []string{"ACTIVE"},
				Refresh:    ServerStateRefreshFunc(computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			log.Printf("[DEBUG] Waiting for instance (%s) to start/unshelve", d.Id())
			_, err = startStateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to become active: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("metadata") {
		oldMetadata, newMetadata := d.GetChange("metadata")
		var metadataToDelete []string

		// Determine if any metadata keys were removed from the configuration.
		// Then request those keys to be deleted.
		for oldKey := range oldMetadata.(map[string]interface{}) {
			var found bool
			for newKey := range newMetadata.(map[string]interface{}) {
				if oldKey == newKey {
					found = true
				}
			}

			if !found {
				metadataToDelete = append(metadataToDelete, oldKey)
			}
		}

		for _, key := range metadataToDelete {
			err := iservers.DeleteMetadatum(computeClient, d.Id(), key).ExtractErr()
			if err != nil {
				return diag.Errorf("Error deleting metadata (%s) from server (%s): %s", key, d.Id(), err)
			}
		}

		// Update existing metadata and add any new metadata.
		metadataOpts := make(servers.MetadataOpts)
		for k, v := range newMetadata.(map[string]interface{}) {
			metadataOpts[k] = v.(string)
		}

		_, err := iservers.UpdateMetadata(computeClient, d.Id(), metadataOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating VKCS server (%s) metadata: %s", d.Id(), err)
		}
	}

	if d.HasChange("security_groups") {
		oldSGRaw, newSGRaw := d.GetChange("security_groups")
		oldSGSet := oldSGRaw.(*schema.Set)
		newSGSet := newSGRaw.(*schema.Set)
		secgroupsToAdd := newSGSet.Difference(oldSGSet)
		secgroupsToRemove := oldSGSet.Difference(newSGSet)

		log.Printf("[DEBUG] Security groups to add: %v", secgroupsToAdd)

		log.Printf("[DEBUG] Security groups to remove: %v", secgroupsToRemove)

		for _, g := range secgroupsToRemove.List() {
			err := isecgroups.RemoveServer(computeClient, d.Id(), g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				if errutil.IsNotFound(err) {
					continue
				}

				return diag.Errorf("Error removing security group (%s) from VKCS server (%s): %s", g, d.Id(), err)
			}
			log.Printf("[DEBUG] Removed security group (%s) from instance (%s)", g, d.Id())
		}

		for _, g := range secgroupsToAdd.List() {
			err := isecgroups.AddServer(computeClient, d.Id(), g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				return diag.Errorf("Error adding security group (%s) to VKCS server (%s): %s", g, d.Id(), err)
			}
			log.Printf("[DEBUG] Added security group (%s) to instance (%s)", g, d.Id())
		}
	}

	if d.HasChange("admin_pass") {
		if newPwd, ok := d.Get("admin_pass").(string); ok {
			err := iservers.ChangeAdminPassword(computeClient, d.Id(), newPwd).ExtractErr()
			if err != nil {
				return diag.Errorf("Error changing admin password of VKCS server (%s): %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("flavor_id") || d.HasChange("flavor_name") {
		// Get vendor_options
		vendorOptionsRaw := d.Get("vendor_options").(*schema.Set)
		var ignoreResizeConfirmation bool
		if vendorOptionsRaw.Len() > 0 {
			vendorOptions := util.ExpandVendorOptions(vendorOptionsRaw.List())
			ignoreResizeConfirmation = vendorOptions["ignore_resize_confirmation"].(bool)
		}

		var newFlavorID string
		var err error
		if d.HasChange("flavor_id") {
			newFlavorID = d.Get("flavor_id").(string)
		} else {
			newFlavorName := d.Get("flavor_name").(string)
			newFlavorID, err = flavorsutils.IDFromName(computeClient, newFlavorName)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		resizeOpts := &servers.ResizeOpts{
			FlavorRef: newFlavorID,
		}
		log.Printf("[DEBUG] Resize configuration: %#v", resizeOpts)
		err = iservers.Resize(computeClient, d.Id(), resizeOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("Error resizing VKCS server: %s", err)
		}

		// Wait for the instance to finish resizing.
		log.Printf("[DEBUG] Waiting for instance (%s) to finish resizing", d.Id())

		// Resize instance without confirmation if specified by user.
		if ignoreResizeConfirmation {
			stateConf := &retry.StateChangeConf{
				Pending:    []string{"RESIZE", "VERIFY_RESIZE"},
				Target:     []string{"ACTIVE", "SHUTOFF"},
				Refresh:    ServerStateRefreshFunc(computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to resize: %s", d.Id(), err)
			}
		} else {
			stateConf := &retry.StateChangeConf{
				Pending:    []string{"RESIZE"},
				Target:     []string{"VERIFY_RESIZE"},
				Refresh:    ServerStateRefreshFunc(computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to resize: %s", d.Id(), err)
			}

			// Confirm resize.
			log.Printf("[DEBUG] Confirming resize")
			err = iservers.ConfirmResize(computeClient, d.Id()).ExtractErr()
			if err != nil {
				return diag.Errorf("Error confirming resize of VKCS server: %s", err)
			}

			stateConf = &retry.StateChangeConf{
				Pending:    []string{"VERIFY_RESIZE"},
				Target:     []string{"ACTIVE", "SHUTOFF"},
				Refresh:    ServerStateRefreshFunc(computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to confirm resize: %s", d.Id(), err)
			}
		}
	}

	// 	Perform any required updates to the tags.
	if d.HasChange("tags") {
		instanceTags := ComputeInstanceUpdateTags(d)
		instanceTagsOpts := tags.ReplaceAllOpts{Tags: instanceTags}
		computeClient.Microversion = computeAPIMicroVersion
		instanceTags, err := tags.ReplaceAll(computeClient, d.Id(), instanceTagsOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on vkcs_compute_instance %s: %s", d.Id(), err)
		}
		log.Printf("[DEBUG] Set tags %s on vkcs_compute_instance %s", instanceTags, d.Id())
	}

	return resourceComputeInstanceRead(ctx, d, meta)
}

func resourceComputeInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	if d.Get("stop_before_destroy").(bool) {
		err = istartstop.Stop(computeClient, d.Id()).ExtractErr()
		if err != nil {
			log.Printf("[WARN] Error stopping vkcs_compute_instance: %s", err)
		} else {
			stopStateConf := &retry.StateChangeConf{
				Pending:    []string{"ACTIVE"},
				Target:     []string{"SHUTOFF"},
				Refresh:    ServerStateRefreshFunc(computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutDelete),
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}
			log.Printf("[DEBUG] Waiting for instance (%s) to stop", d.Id())
			_, err = stopStateConf.WaitForStateContext(ctx)
			if err != nil {
				log.Printf("[WARN] Error waiting for instance (%s) to stop: %s, proceeding to delete", d.Id(), err)
			}
		}
	}
	vendorOptionsRaw := d.Get("vendor_options").(*schema.Set)
	var detachPortBeforeDestroy bool
	if vendorOptionsRaw.Len() > 0 {
		vendorOptions := util.ExpandVendorOptions(vendorOptionsRaw.List())
		detachPortBeforeDestroy = vendorOptions["detach_ports_before_destroy"].(bool)
	}
	if detachPortBeforeDestroy {
		allInstanceNetworks, err := getAllInstanceNetworks(d, meta)
		if err != nil {
			log.Printf("[WARN] Unable to get vkcs_compute_instance ports: %s", err)
		} else {
			for _, network := range allInstanceNetworks {
				if network.Port != "" {
					stateConf := &retry.StateChangeConf{
						Pending:    []string{""},
						Target:     []string{"DETACHED"},
						Refresh:    computeInterfaceAttachDetachFunc(computeClient, d.Id(), network.Port),
						Timeout:    d.Timeout(schema.TimeoutDelete),
						Delay:      5 * time.Second,
						MinTimeout: 5 * time.Second,
					}
					if _, err = stateConf.WaitForStateContext(ctx); err != nil {
						return diag.Errorf("Error detaching vkcs_compute_instance %s: %s", d.Id(), err)
					}
				}
			}
		}
	}
	if d.Get("force_delete").(bool) {
		log.Printf("[DEBUG] Force deleting VKCS Instance %s", d.Id())
		err = iservers.ForceDelete(computeClient, d.Id()).ExtractErr()
		if err != nil {
			return diag.FromErr(util.CheckDeleted(d, err, "Error force deleting vkcs_compute_instance"))
		}
	} else {
		log.Printf("[DEBUG] Deleting VKCS Instance %s", d.Id())
		err = iservers.Delete(computeClient, d.Id()).ExtractErr()
		if err != nil {
			return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_compute_instance"))
		}
	}

	log.Printf("[DEBUG] Deleting VKCS Instance %s", d.Id())
	err = iservers.Delete(computeClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_compute_instance"))
	}
	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE", "SHUTOFF"},
		Target:     []string{"DELETED", "SOFT_DELETED"},
		Refresh:    ServerStateRefreshFunc(computeClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for instance (%s) to Delete:  %s",
			d.Id(), err)
	}

	return nil
}

func resourceComputeInstanceCustomizeDiff(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	bdsRaw := diff.Get("block_device").([]interface{})
	bds := make([]map[string]interface{}, len(bdsRaw))
	for i, v := range bdsRaw {
		defaultBootIdx := -1
		if i == 0 && len(bds) == 1 {
			defaultBootIdx = 0
		}
		bd := v.(map[string]interface{})
		key := fmt.Sprintf("block_device.%d.boot_index", i)
		if _, ok := diff.GetOkExists(key); !ok {
			bd["boot_index"] = defaultBootIdx
		}
		bds[i] = bd
	}
	diff.SetNew("block_device", bds)
	return nil
}

func resourceComputeInstanceImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var serverWithAttachments struct {
		VolumesAttached []map[string]interface{} `json:"os-extended-volumes:volumes_attached"`
	}

	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("error creating VKCS compute client: %s", err)
	}

	results := make([]*schema.ResourceData, 1)
	diagErr := resourceComputeInstanceRead(ctx, d, meta)
	if diagErr != nil {
		return nil, fmt.Errorf("error reading vkcs_compute_instance %s: %v", d.Id(), diagErr)
	}

	raw := iservers.Get(computeClient, d.Id())
	if raw.Err != nil {
		return nil, util.CheckDeleted(d, raw.Err, "vkcs_compute_instance")
	}

	if err := raw.ExtractInto(&serverWithAttachments); err != nil {
		log.Printf("[DEBUG] unable to unmarshal raw struct to serverWithAttachments: %s", err)
	}

	log.Printf("[DEBUG] Retrieved vkcs_compute_instance %s volume attachments: %#v",
		d.Id(), serverWithAttachments)

	bds := []map[string]interface{}{}
	if len(serverWithAttachments.VolumesAttached) > 0 {
		blockStorageClient, err := config.BlockStorageV3Client(util.GetRegion(d, config))
		if err != nil {
			return nil, fmt.Errorf("error creating VKCS volume client: %s", err)
		}
		var volMetaData = struct {
			VolumeImageMetadata map[string]interface{} `json:"volume_image_metadata"`
			ID                  string                 `json:"id"`
			Size                int                    `json:"size"`
			Bootable            string                 `json:"bootable"`
		}{}
		for _, b := range serverWithAttachments.VolumesAttached {
			rawVolume := ivolumes.Get(blockStorageClient, b["id"].(string))
			if err := rawVolume.ExtractInto(&volMetaData); err != nil {
				log.Printf("[DEBUG] unable to unmarshal raw struct to volume metadata: %s", err)
			}

			log.Printf("[DEBUG] retrieved volume%+v", volMetaData)
			v := map[string]interface{}{
				"delete_on_termination": true,
				"uuid":                  volMetaData.VolumeImageMetadata["image_id"],
				"boot_index":            -1,
				"destination_type":      "volume",
				"source_type":           "image",
				"volume_size":           volMetaData.Size,
				"disk_bus":              "",
				"volume_type":           "",
				"device_type":           "",
			}

			if volMetaData.Bootable == "true" {
				bds = append(bds, v)
			}
		}

		bds[0]["boot_index"] = 0

		d.Set("block_device", bds)
	}

	metadata, err := iservers.Metadata(computeClient, d.Id()).Extract()
	if err != nil {
		return nil, fmt.Errorf("unable to read metadata for vkcs_compute_instance %s: %s", d.Id(), err)
	}

	d.Set("metadata", metadata)

	results[0] = d

	return results, nil
}

// ServerStateRefreshFunc returns a retry.StateRefreshFunc that is used to watch
// an VKCS instance.
func ServerStateRefreshFunc(client *gophercloud.ServiceClient, instanceID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := iservers.Get(client, instanceID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return s, "DELETED", nil
			}
			return nil, "", err
		}

		return s, s.Status, nil
	}
}

func resourceInstanceSecGroupsV2(d *schema.ResourceData) []string {
	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	res := make([]string, len(rawSecGroups))
	for i, raw := range rawSecGroups {
		res[i] = raw.(string)
	}
	return res
}

func resourceInstanceMetadataV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func ResourceInstanceBlockDevicesV2(_ *schema.ResourceData, bds []interface{}) ([]bootfromvolume.BlockDevice, error) {
	blockDeviceOpts := make([]bootfromvolume.BlockDevice, len(bds))
	for i, bd := range bds {
		bdM := bd.(map[string]interface{})
		blockDeviceOpts[i] = bootfromvolume.BlockDevice{
			UUID:                bdM["uuid"].(string),
			VolumeSize:          bdM["volume_size"].(int),
			BootIndex:           bdM["boot_index"].(int),
			DeleteOnTermination: bdM["delete_on_termination"].(bool),
			GuestFormat:         bdM["guest_format"].(string),
			VolumeType:          bdM["volume_type"].(string),
			DeviceType:          bdM["device_type"].(string),
			DiskBus:             bdM["disk_bus"].(string),
		}
		sourceType := bdM["source_type"].(string)
		switch sourceType {
		case "blank":
			blockDeviceOpts[i].SourceType = bootfromvolume.SourceBlank
		case "image":
			blockDeviceOpts[i].SourceType = bootfromvolume.SourceImage
		case "snapshot":
			blockDeviceOpts[i].SourceType = bootfromvolume.SourceSnapshot
		case "volume":
			blockDeviceOpts[i].SourceType = bootfromvolume.SourceVolume
		default:
			return blockDeviceOpts, fmt.Errorf("unknown block device source type %s", sourceType)
		}

		destinationType := bdM["destination_type"].(string)
		switch destinationType {
		case "local":
			blockDeviceOpts[i].DestinationType = bootfromvolume.DestinationLocal
		case "volume":
			blockDeviceOpts[i].DestinationType = bootfromvolume.DestinationVolume
		default:
			return blockDeviceOpts, fmt.Errorf("unknown block device destination type %s", destinationType)
		}
	}

	if len(blockDeviceOpts) == 1 {
		blockDeviceOpts[0].BootIndex = 0
	}

	log.Printf("[DEBUG] Block Device Options: %+v", blockDeviceOpts)
	return blockDeviceOpts, nil
}

func resourceInstanceSchedulerHintsV2(_ *schema.ResourceData, schedulerHintsRaw map[string]interface{}) schedulerhints.SchedulerHints {
	schedulerHints := schedulerhints.SchedulerHints{
		Group: schedulerHintsRaw["group"].(string),
	}

	return schedulerHints
}

func getImageIDFromConfig(imageClient *gophercloud.ServiceClient, d *schema.ResourceData) (string, error) {
	// If block_device was used, an Image does not need to be specified, unless an image/local
	// combination was used. This emulates normal boot behavior. Otherwise, ignore the image altogether.
	if vL, ok := d.GetOk("block_device"); ok {
		needImage := false
		for _, v := range vL.([]interface{}) {
			vM := v.(map[string]interface{})
			if vM["source_type"] == "image" && vM["destination_type"] == "local" {
				needImage = true
			}
		}
		if !needImage {
			return "", nil
		}
	}

	if imageID := d.Get("image_id").(string); imageID != "" {
		return imageID, nil
	}
	// try the OS_IMAGE_ID environment variable
	if v := os.Getenv("OS_IMAGE_ID"); v != "" {
		return v, nil
	}

	imageName := d.Get("image_name").(string)
	if imageName == "" {
		// try the OS_IMAGE_NAME environment variable
		if v := os.Getenv("OS_IMAGE_NAME"); v != "" {
			imageName = v
		}
	}

	if imageName != "" {
		imageID, err := imagesutils.IDFromName(imageClient, imageName)
		if err != nil {
			return "", err
		}
		return imageID, nil
	}

	return "", fmt.Errorf("neither a boot device, image ID, or image name were able to be determined")
}

func setImageInformation(computeClient *gophercloud.ServiceClient, server *servers.Server, d *schema.ResourceData) error {
	// If block_device was used, an Image does not need to be specified, unless an image/local
	// combination was used. This emulates normal boot behavior. Otherwise, ignore the image altogether.
	if vL, ok := d.GetOk("block_device"); ok {
		needImage := false
		for _, v := range vL.([]interface{}) {
			vM := v.(map[string]interface{})
			if vM["source_type"] == "image" && vM["destination_type"] == "local" {
				needImage = true
			}
		}
		if !needImage {
			d.Set("image_id", "Attempt to boot from volume - no image supplied")
			return nil
		}
	}

	if server.Image["id"] != nil {
		imageID := server.Image["id"].(string)
		if imageID != "" {
			d.Set("image_id", imageID)
			image, err := iimages.Get(computeClient, imageID).Extract()
			if err != nil {
				if errutil.IsNotFound(err) {
					// If the image name can't be found, set the value to "Image not found".
					// The most likely scenario is that the image no longer exists in the Image Service
					// but the instance still has a record from when it existed.
					d.Set("image_name", "Image not found")
					return nil
				}
				return err
			}
			d.Set("image_name", image.Name)
		}
	}

	return nil
}

func getFlavorID(computeClient *gophercloud.ServiceClient, d *schema.ResourceData) (string, error) {
	if flavorID := d.Get("flavor_id").(string); flavorID != "" {
		return flavorID, nil
	}
	// Try the OS_FLAVOR_ID environment variable
	if v := os.Getenv("OS_FLAVOR_ID"); v != "" {
		return v, nil
	}

	flavorName := d.Get("flavor_name").(string)
	if flavorName == "" {
		// Try the OS_FLAVOR_NAME environment variable
		if v := os.Getenv("OS_FLAVOR_NAME"); v != "" {
			flavorName = v
		}
	}

	if flavorName != "" {
		flavorID, err := flavorsutils.IDFromName(computeClient, flavorName)
		if err != nil {
			return "", err
		}
		return flavorID, nil
	}

	return "", fmt.Errorf("neither a flavor_id or flavor_name could be determined")
}

func resourceComputeSchedulerHintsHash(v interface{}) int {
	var buf bytes.Buffer

	m, ok := v.(map[string]interface{})
	if !ok {
		return hashcode.String(buf.String())
	}
	if m == nil {
		return hashcode.String(buf.String())
	}

	if m["group"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["group"].(string)))
	}

	return hashcode.String(buf.String())
}

func checkBlockDeviceConfig(d *schema.ResourceData) diag.Diagnostics {
	vL, ok := d.GetOk("block_device")
	if !ok {
		return nil
	}

	vLs := vL.([]interface{})
	isBootIdxZeroSet := len(vLs) <= 1
	for _, v := range vLs {
		vM := v.(map[string]interface{})

		if vM["source_type"] != "blank" && vM["uuid"] == "" {
			return diag.Errorf("you must specify a uuid for %s block device types", vM["source_type"])
		}

		if vM["source_type"] == "image" && vM["destination_type"] == "volume" {
			if vM["volume_size"] == 0 {
				return diag.Errorf("you must specify a volume_size when creating a volume from an image")
			}
		}

		if vM["source_type"] == "blank" && vM["destination_type"] == "local" {
			if vM["volume_size"] == 0 {
				return diag.Errorf("you must specify a volume_size when creating a blank block device")
			}
		}

		if vM["boot_index"] == 0 {
			isBootIdxZeroSet = true
		}
	}

	if !isBootIdxZeroSet {
		return diag.Errorf("you must set boot_index to 0 for one of block_devices")
	}

	diags := diag.Diagnostics{}
	for _, v := range vLs {
		vM := v.(map[string]interface{})
		deleteOnTermination := vM["delete_on_termination"].(bool)
		sourceType := vM["source_type"].(string)
		switch {
		case sourceType == "blank" || sourceType == "image" || sourceType == "snapshot":
			if !deleteOnTermination {
				path := cty.Path{
					cty.GetAttrStep{Name: "block_device"},
					cty.IndexStep{Key: cty.StringVal(vM["uuid"].(string))},
					cty.GetAttrStep{Name: "delete_on_termination"},
				}
				diags = append(diags, diag.Diagnostic{
					Severity:      diag.Warning,
					Summary:       fmt.Sprintf("delete_on_termination should be true, when source_type is %s", sourceType),
					AttributePath: path,
				})
			}
		case sourceType == "volume":
			if deleteOnTermination {
				path := cty.Path{
					cty.GetAttrStep{Name: "block_device"},
					cty.IndexStep{Key: cty.StringVal(vM["uuid"].(string))},
					cty.GetAttrStep{Name: "delete_on_termination"},
				}
				diags = append(diags, diag.Diagnostic{
					Severity:      diag.Warning,
					Summary:       "delete_on_termination should be false, when source_type is volume",
					AttributePath: path,
				})
			}
		}
	}

	return diags
}

func resourceComputeInstancePersonalityHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["file"].(string)))

	return hashcode.String(buf.String())
}

func resourceInstancePersonalityV2(d *schema.ResourceData) servers.Personality {
	var personalities servers.Personality

	if v := d.Get("personality"); v != nil {
		personalityList := v.(*schema.Set).List()
		if len(personalityList) > 0 {
			for _, p := range personalityList {
				rawPersonality := p.(map[string]interface{})
				file := servers.File{
					Path:     rawPersonality["file"].(string),
					Contents: []byte(rawPersonality["content"].(string)),
				}

				log.Printf("[DEBUG] vkcs_compute_instance Personality: %+v", file)

				personalities = append(personalities, &file)
			}
		}
	}

	return personalities
}

// suppressAvailabilityZoneDetailDiffs will suppress diffs when a user specifies an
// availability zone in the format of `az:host:node` and Nova/Compute responds with
// only `az`.
func suppressAvailabilityZoneDetailDiffs(_, old, new string, _ *schema.ResourceData) bool {
	if strings.Contains(new, ":") {
		parts := strings.Split(new, ":")
		az := parts[0]

		if az == old {
			return true
		}
	}

	return false
}

// suppressPowerStateDiffs will allow a state of "error" or "migrating" even though we don't
// allow them as a user input.
func suppressPowerStateDiffs(_, old, _ string, _ *schema.ResourceData) bool {
	if old == "error" || old == "migrating" {
		return true
	}

	return false
}
