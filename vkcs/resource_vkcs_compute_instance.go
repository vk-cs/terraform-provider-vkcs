package vkcs

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

	"github.com/gophercloud/gophercloud"
	volumesV3 "github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/availabilityzones"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/schedulerhints"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/shelveunshelve"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/startstop"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/tags"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	flavorsutils "github.com/gophercloud/utils/openstack/compute/v2/flavors"
	imagesutils "github.com/gophercloud/utils/openstack/imageservice/v2/images"
	"github.com/gophercloud/utils/terraform/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceComputeInstance() *schema.Resource {
	return &schema.Resource{
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
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
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"availability_zone": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: suppressAvailabilityZoneDetailDiffs,
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
			},
			"network": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"fixed_ip_v4": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"fixed_ip_v6": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_network": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			"config_drive": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"admin_pass": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  false,
			},
			"access_ip_v4": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: false,
			},
			"access_ip_v6": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: false,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"block_device": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"volume_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"destination_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"boot_index": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"delete_on_termination": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"guest_format": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"device_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disk_bus": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"scheduler_hints": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
				Set: resourceComputeSchedulerHintsHash,
			},
			"personality": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file": {
							Type:     schema.TypeString,
							Required: true,
						},
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: resourceComputeInstancePersonalityHash,
			},
			"stop_before_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"all_metadata": {
				Type:     schema.TypeMap,
				Computed: true,
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
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"all_tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vendor_options": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ignore_resize_confirmation": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"detach_ports_before_destroy": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceComputeInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	imageClient, err := config.ImageV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS image client: %s", err)
	}

	computeClient.Microversion = computeAPIMicroVersion
	var createOpts servers.CreateOptsBuilder
	var availabilityZone string
	var networks interface{}

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
	if err := checkBlockDeviceConfig(d); err != nil {
		return diag.FromErr(err)
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
		blockDevices, err := resourceInstanceBlockDevicesV2(d, vL.([]interface{}))
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
		server, err = bootfromvolume.Create(computeClient, createOpts).Extract()
	} else {
		server, err = servers.Create(computeClient, createOpts).Extract()
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

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    ServerStateRefreshFunc(computeClient, server.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	err = resource.RetryContext(ctx, stateConf.Timeout, func() *resource.RetryError {
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			log.Printf("[DEBUG] Retrying after error: %s", err)
			return checkForRetryableError(err)
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
		err = startstop.Stop(computeClient, d.Id()).ExtractErr()
		if err != nil {
			return diag.Errorf("Error stopping VKCS instance: %s", err)
		}
		stopStateConf := &resource.StateChangeConf{
			//Pending:    []string{"ACTIVE"},
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

	return resourceComputeInstanceRead(ctx, d, meta)
}

func resourceComputeInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	server, err := servers.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "server"))
	}

	log.Printf("[DEBUG] Retrieved Server %s: %+v", d.Id(), server)

	d.Set("name", server.Name)

	// Get the instance network and address information
	networks, err := flattenInstanceNetworks(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// Determine the best IPv4 and IPv6 addresses to access the instance with
	hostv4, hostv6 := getInstanceAccessAddresses(d, networks)

	// AccessIPv4/v6 isn't standard in OpenStack, but there have been reports
	// of them being used in some environments.
	if server.AccessIPv4 != "" && hostv4 == "" {
		hostv4 = server.AccessIPv4
	}

	if server.AccessIPv6 != "" && hostv6 == "" {
		hostv6 = server.AccessIPv6
	}

	d.Set("network", networks)
	d.Set("access_ip_v4", hostv4)
	d.Set("access_ip_v6", hostv6)

	// Determine the best IP address to use for SSH connectivity.
	// Prefer IPv4 over IPv6.
	var preferredSSHAddress string
	if hostv4 != "" {
		preferredSSHAddress = hostv4
	} else if hostv6 != "" {
		preferredSSHAddress = hostv6
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
	flavor, err := flavors.Get(computeClient, flavorID).Extract()
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
	err = servers.Get(computeClient, d.Id()).ExtractInto(&serverWithAZ)
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "server"))
	}
	// Set the availability zone
	d.Set("availability_zone", serverWithAZ.AvailabilityZone)

	// Set the region
	d.Set("region", getRegion(d, config))

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
	instanceTags, err := tags.List(computeClient, server.ID).Extract()
	if err != nil {
		log.Printf("[DEBUG] Unable to get tags for vkcs_compute_instance: %s", err)
	} else {
		ComputeInstanceReadTags(d, instanceTags)
	}

	return nil
}

func resourceComputeInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	var updateOpts servers.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if updateOpts != (servers.UpdateOpts{}) {
		_, err := servers.Update(computeClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating VKCS server: %s", err)
		}
	}

	if d.HasChange("power_state") {
		powerStateOldRaw, powerStateNewRaw := d.GetChange("power_state")
		powerStateOld := powerStateOldRaw.(string)
		powerStateNew := powerStateNewRaw.(string)
		if strings.ToLower(powerStateNew) == "shelved_offloaded" {
			err = shelveunshelve.Shelve(computeClient, d.Id()).ExtractErr()
			if err != nil {
				return diag.Errorf("Error shelve VKCS instance: %s", err)
			}
			shelveStateConf := &resource.StateChangeConf{
				//Pending:    []string{"ACTIVE"},
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
			err = startstop.Stop(computeClient, d.Id()).ExtractErr()
			if err != nil {
				return diag.Errorf("Error stopping VKCS instance: %s", err)
			}
			stopStateConf := &resource.StateChangeConf{
				//Pending:    []string{"ACTIVE"},
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
				err = shelveunshelve.Unshelve(computeClient, d.Id(), unshelveOpt).ExtractErr()
				if err != nil {
					return diag.Errorf("Error unshelving VKCS instance: %s", err)
				}
			} else {
				err = startstop.Start(computeClient, d.Id()).ExtractErr()
				if err != nil {
					return diag.Errorf("Error starting VKCS instance: %s", err)
				}
			}
			startStateConf := &resource.StateChangeConf{
				//Pending:    []string{"SHUTOFF"},
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
			err := servers.DeleteMetadatum(computeClient, d.Id(), key).ExtractErr()
			if err != nil {
				return diag.Errorf("Error deleting metadata (%s) from server (%s): %s", key, d.Id(), err)
			}
		}

		// Update existing metadata and add any new metadata.
		metadataOpts := make(servers.MetadataOpts)
		for k, v := range newMetadata.(map[string]interface{}) {
			metadataOpts[k] = v.(string)
		}

		_, err := servers.UpdateMetadata(computeClient, d.Id(), metadataOpts).Extract()
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
			err := secgroups.RemoveServer(computeClient, d.Id(), g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				if _, ok := err.(gophercloud.ErrDefault404); ok {
					continue
				}

				return diag.Errorf("Error removing security group (%s) from VKCS server (%s): %s", g, d.Id(), err)
			}
			log.Printf("[DEBUG] Removed security group (%s) from instance (%s)", g, d.Id())
		}

		for _, g := range secgroupsToAdd.List() {
			err := secgroups.AddServer(computeClient, d.Id(), g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				return diag.Errorf("Error adding security group (%s) to VKCS server (%s): %s", g, d.Id(), err)
			}
			log.Printf("[DEBUG] Added security group (%s) to instance (%s)", g, d.Id())
		}
	}

	if d.HasChange("admin_pass") {
		if newPwd, ok := d.Get("admin_pass").(string); ok {
			err := servers.ChangeAdminPassword(computeClient, d.Id(), newPwd).ExtractErr()
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
			vendorOptions := expandVendorOptions(vendorOptionsRaw.List())
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
		err = servers.Resize(computeClient, d.Id(), resizeOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("Error resizing VKCS server: %s", err)
		}

		// Wait for the instance to finish resizing.
		log.Printf("[DEBUG] Waiting for instance (%s) to finish resizing", d.Id())

		// Resize instance without confirmation if specified by user.
		if ignoreResizeConfirmation {
			stateConf := &resource.StateChangeConf{
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
			stateConf := &resource.StateChangeConf{
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
			err = servers.ConfirmResize(computeClient, d.Id()).ExtractErr()
			if err != nil {
				return diag.Errorf("Error confirming resize of VKCS server: %s", err)
			}

			stateConf = &resource.StateChangeConf{
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
	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	if d.Get("stop_before_destroy").(bool) {
		err = startstop.Stop(computeClient, d.Id()).ExtractErr()
		if err != nil {
			log.Printf("[WARN] Error stopping vkcs_compute_instance: %s", err)
		} else {
			stopStateConf := &resource.StateChangeConf{
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
		vendorOptions := expandVendorOptions(vendorOptionsRaw.List())
		detachPortBeforeDestroy = vendorOptions["detach_ports_before_destroy"].(bool)
	}
	if detachPortBeforeDestroy {
		allInstanceNetworks, err := getAllInstanceNetworks(d, meta)
		if err != nil {
			log.Printf("[WARN] Unable to get vkcs_compute_instance ports: %s", err)
		} else {
			for _, network := range allInstanceNetworks {
				if network.Port != "" {
					stateConf := &resource.StateChangeConf{
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
		err = servers.ForceDelete(computeClient, d.Id()).ExtractErr()
		if err != nil {
			return diag.FromErr(checkDeleted(d, err, "Error force deleting vkcs_compute_instance"))
		}
	} else {
		log.Printf("[DEBUG] Deleting VKCS Instance %s", d.Id())
		err = servers.Delete(computeClient, d.Id()).ExtractErr()
		if err != nil {
			return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_compute_instance"))
		}
	}

	log.Printf("[DEBUG] Deleting VKCS Instance %s", d.Id())
	err = servers.Delete(computeClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_compute_instance"))
	}
	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
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

func resourceComputeInstanceImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var serverWithAttachments struct {
		VolumesAttached []map[string]interface{} `json:"os-extended-volumes:volumes_attached"`
	}

	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("error creating VKCS compute client: %s", err)
	}

	results := make([]*schema.ResourceData, 1)
	diagErr := resourceComputeInstanceRead(ctx, d, meta)
	if diagErr != nil {
		return nil, fmt.Errorf("error reading vkcs_compute_instance %s: %v", d.Id(), diagErr)
	}

	raw := servers.Get(computeClient, d.Id())
	if raw.Err != nil {
		return nil, checkDeleted(d, raw.Err, "vkcs_compute_instance")
	}

	if err := raw.ExtractInto(&serverWithAttachments); err != nil {
		log.Printf("[DEBUG] unable to unmarshal raw struct to serverWithAttachments: %s", err)
	}

	log.Printf("[DEBUG] Retrieved vkcs_compute_instance %s volume attachments: %#v",
		d.Id(), serverWithAttachments)

	bds := []map[string]interface{}{}
	if len(serverWithAttachments.VolumesAttached) > 0 {
		blockStorageClient, err := config.BlockStorageV3Client(getRegion(d, config))
		if err != nil {
			return nil, fmt.Errorf("error creating VKCS volume client: %s", err)
		}
		var volMetaData = struct {
			VolumeImageMetadata map[string]interface{} `json:"volume_image_metadata"`
			ID                  string                 `json:"id"`
			Size                int                    `json:"size"`
			Bootable            string                 `json:"bootable"`
		}{}
		for i, b := range serverWithAttachments.VolumesAttached {
			rawVolume := volumesV3.Get(blockStorageClient, b["id"].(string))
			if err := rawVolume.ExtractInto(&volMetaData); err != nil {
				log.Printf("[DEBUG] unable to unmarshal raw struct to volume metadata: %s", err)
			}

			log.Printf("[DEBUG] retrieved volume%+v", volMetaData)
			v := map[string]interface{}{
				"delete_on_termination": true,
				"uuid":                  volMetaData.VolumeImageMetadata["image_id"],
				"boot_index":            i,
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

		d.Set("block_device", bds)
	}

	metadata, err := servers.Metadata(computeClient, d.Id()).Extract()
	if err != nil {
		return nil, fmt.Errorf("unable to read metadata for vkcs_compute_instance %s: %s", d.Id(), err)
	}

	d.Set("metadata", metadata)

	results[0] = d

	return results, nil
}

// ServerStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an VKCS instance.
func ServerStateRefreshFunc(client *gophercloud.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := servers.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
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

func resourceInstanceBlockDevicesV2(_ *schema.ResourceData, bds []interface{}) ([]bootfromvolume.BlockDevice, error) {
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
			image, err := images.Get(computeClient, imageID).Extract()
			if err != nil {
				if _, ok := err.(gophercloud.ErrDefault404); ok {
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

func checkBlockDeviceConfig(d *schema.ResourceData) error {
	if vL, ok := d.GetOk("block_device"); ok {
		for _, v := range vL.([]interface{}) {
			vM := v.(map[string]interface{})

			if vM["source_type"] != "blank" && vM["uuid"] == "" {
				return fmt.Errorf("you must specify a uuid for %s block device types", vM["source_type"])
			}

			if vM["source_type"] == "image" && vM["destination_type"] == "volume" {
				if vM["volume_size"] == 0 {
					return fmt.Errorf("you must specify a volume_size when creating a volume from an image")
				}
			}

			if vM["source_type"] == "blank" && vM["destination_type"] == "local" {
				if vM["volume_size"] == 0 {
					return fmt.Errorf("you must specify a volume_size when creating a blank block device")
				}
			}
		}
	}

	return nil
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
