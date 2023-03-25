package vkcs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Dbaas timeouts
const (
	dbInstanceDelay         = 10 * time.Second
	dbInstanceMinTimeout    = 3 * time.Second
	dbDatabaseDelay         = 10 * time.Second
	dbDatabaseMinTimeout    = 3 * time.Second
	dbUserDelay             = 10 * time.Second
	dbUserMinTimeout        = 3 * time.Second
	dbCreateTimeout         = 30 * time.Minute
	dbDeleteTimeout         = 30 * time.Minute
	dbUserCreateTimeout     = 10 * time.Minute
	dbUserDeleteTimeout     = 10 * time.Minute
	dbDatabaseCreateTimeout = 10 * time.Minute
	dbDatabaseDeleteTimeout = 10 * time.Minute
)

type dbInstanceStatus string

var (
	dbInstanceStatusDeleted            dbInstanceStatus = "DELETED"
	dbInstanceStatusBuild              dbInstanceStatus = "BUILD"
	dbInstanceStatusActive             dbInstanceStatus = "ACTIVE"
	dbInstanceStatusError              dbInstanceStatus = "ERROR"
	dbInstanceStatusShutdown           dbInstanceStatus = "SHUTDOWN"
	dbInstanceStatusResize             dbInstanceStatus = "RESIZE"
	dbInstanceStatusDetach             dbInstanceStatus = "DETACH"
	dbInstanceStatusCapabilityApplying dbInstanceStatus = "CAPABILITY_APPLYING"
	dbInstanceStatusBackup             dbInstanceStatus = "BACKUP"
)

type dbCapabilityStatus string

var (
	dbCapabilityStatusActive dbCapabilityStatus = "ACTIVE"
	dbCapabilityStatusError  dbCapabilityStatus = "ERROR"
)

const (
	dbImportedStatus = "IMPORTED"
)

func resourceDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseInstanceCreate,
		ReadContext:   resourceDatabaseInstanceRead,
		DeleteContext: resourceDatabaseInstanceDelete,
		UpdateContext: resourceDatabaseInstanceUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				config := meta.(configer)
				DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
				if err != nil {
					return nil, fmt.Errorf("error creating VKCS database client: %s", err)
				}

				if resourceDatabaseInstanceRead(ctx, d, meta).HasError() {
					return nil, fmt.Errorf("error reading vkcs_db_instance")
				}

				capabilities, err := instanceGetCapabilities(DatabaseV1Client, d.Id()).extract()
				if err != nil {
					return nil, fmt.Errorf("error getting instance capabilities")
				}
				d.Set("capabilities", flattenDatabaseInstanceCapabilities(capabilities))
				d.Set("volume_type", dbImportedStatus)
				if v, ok := d.GetOk("wal_volume"); ok {
					walV, _ := extractDatabaseWalVolume(v.([]interface{}))
					walvolume := walVolume{Size: &walV.Size, VolumeType: dbImportedStatus}
					d.Set("wal_volume", flattenDatabaseInstanceWalVolume(walvolume))
				}

				return []*schema.ResourceData{d}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(dbCreateTimeout),
			Delete: schema.DefaultTimeout(dbDeleteTimeout),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Region to create resource in.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the instance. Changing this creates a new instance",
			},

			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Computed:    false,
				Description: "The ID of flavor for the instance.",
			},

			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    false,
				Computed:    false,
				Description: "Size of the instance volume.",
			},

			"volume_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Computed:    false,
				Description: "The type of the instance volume. Changing this creates a new instance.",
			},

			"wal_volume": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    false,
							Description: "Size of the instance wal volume.",
						},
						"volume_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The type of the instance wal volume.",
						},
						"autoexpand": {
							Type:          schema.TypeBool,
							Optional:      true,
							ForceNew:      false,
							Deprecated:    "Please, use wal_disk_autoexpand block instead",
							ConflictsWith: []string{"wal_disk_autoexpand.0.autoexpand"},
							Description:   "Indicates whether wal volume autoresize is enabled. ***Deprecated***. Please, use wal_disk_autoexpand block instead.",
						},
						"max_disk_size": {
							Type:          schema.TypeInt,
							Optional:      true,
							ForceNew:      false,
							Deprecated:    "Please, use wal_disk_autoexpand block instead",
							ConflictsWith: []string{"wal_disk_autoexpand.0.max_disk_size"},
							Description:   "Maximum disk size for wal volume autoresize. ***Deprecated***. Please, use wal_disk_autoexpand block instead.",
						},
					},
				},
				Description: "Object that represents wal volume of the instance. Changing this creates a new instance.",
			},

			"datastore": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Version of the datastore. Changing this creates a new instance.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Type of the datastore. Changing this creates a new instance.",
						},
					},
				},
				Description: "Object that represents datastore of the instance. Changing this creates a new instance.",
			},

			"network": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The id of the network. Changing this creates a new instance.",
						},
						"port": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The port id of the network. Changing this creates a new instance. ***Deprecated*** This argument is deprecated, please do not use it.",
							Deprecated:  "This argument is deprecated, please do not use it.",
						},
						"fixed_ip_v4": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Description: "The IPv4 address. Changing this creates a new instance. " +
								"**Note** This argument conflicts with \"replica_of\". Setting both at the same time causes \"fixed_ip_v4\" to be ignored.",
						},
					},
				},
				Description: "Object that represents network of the instance. Changing this creates a new instance.",
			},

			"configuration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    false,
				Description: "The id of the configuration attached to instance.",
			},

			"replica_of": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    false,
				Description: "ID of the instance, that current instance is replica of.",
			},

			"root_enabled": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      false,
				ConflictsWith: []string{"replica_of"},
				Description:   "Indicates whether root user is enabled for the instance.",
			},

			"root_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Computed:      true,
				ForceNew:      false,
				ConflictsWith: []string{"replica_of"},
				Description:   "Password for the root user of the instance. If this field is empty and root user is enabled, then after creation of the instance this field will contain auto-generated root user password.",
			},

			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The name of the availability zone of the instance. Changing this creates a new instance.",
			},

			"floating_ip_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Indicates whether floating ip is created for instance. Changing this creates a new instance.",
			},

			"keypair": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Name of the keypair to be attached to instance. Changing this creates a new instance.",
			},

			"disk_autoexpand": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"autoexpand": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    false,
							Description: "Indicates whether autoresize is enabled.",
						},
						"max_disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    false,
							Description: "Maximum disk size for autoresize.",
						},
					},
				},
				Description: "Object that represents autoresize properties of the instance.",
			},

			"wal_disk_autoexpand": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"autoexpand": {
							Type:          schema.TypeBool,
							Optional:      true,
							ForceNew:      false,
							ConflictsWith: []string{"wal_volume.0.autoexpand"},
							Description:   "Indicates whether wal volume autoresize is enabled.",
						},
						"max_disk_size": {
							Type:          schema.TypeInt,
							Optional:      true,
							ForceNew:      false,
							ConflictsWith: []string{"wal_volume.0.max_disk_size"},
							Description:   "Maximum disk size for wal volume autoresize.",
						},
					},
				},
				Description: "Object that represents autoresize properties of the instance wal volume.",
			},

			"capabilities": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the capability to apply.",
						},
						"settings": {
							Type:        schema.TypeMap,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "Map of key-value settings of the capability.",
						},
					},
				},
				Description: "Object that represents capability applied to instance. There can be several instances of this object (see example).",
			},
			"restore_point": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "ID of the backup.",
						},
						"target": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Used only for restoring from postgresql PITR backups. Timestamp of needed backup in format \"2021-10-06 01:02:00\". You can specify \"latest\" to use most recent backup.",
						},
					},
				},
				Description: "Object that represents backup to restore instance from. **New since v.0.1.4**.",
			},

			"backup_schedule": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the schedule.",
						},
						"start_hours": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Hours part of timestamp of initial backup.",
						},
						"start_minutes": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Minutes part of timestamp of initial backup.",
						},
						"interval_hours": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Time interval between backups, specified in hours. Available values: 3, 6, 8, 12, 24.",
						},
						"keep_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of backups to be stored.",
						},
					},
				},
				Description: "Object that represents configuration of PITR backup. This functionality is available only for postgres datastore. **New since v.0.1.4**.",
			},

			// Computed values
			"ip": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "IP address of the instance.",
			},
		},
		Description: "Provides a db instance resource. This can be used to create, modify and delete db instance.",
	}
}

func resourceDatabaseInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	size := d.Get("size").(int)
	createOpts := &dbInstanceCreateOpts{
		FlavorRef:         d.Get("flavor_id").(string),
		Name:              d.Get("name").(string),
		Volume:            &volume{Size: &size, VolumeType: d.Get("volume_type").(string)},
		ReplicaOf:         d.Get("replica_of").(string),
		AvailabilityZone:  d.Get("availability_zone").(string),
		FloatingIPEnabled: d.Get("floating_ip_enabled").(bool),
		Keypair:           d.Get("keypair").(string),
	}

	message := "unable to determine vkcs_db_instance"

	if v, ok := d.GetOk("restore_point"); ok {
		restorepoint, err := extractDatabaseRestorePoint(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s restore_point", message)
		}
		createOpts.RestorePoint = &restorepoint
	}

	if v, ok := d.GetOk("datastore"); ok {
		datastore, err := extractDatabaseDatastore(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s datastore", message)
		}
		createOpts.Datastore = &datastore
	}

	if replicaOf, ok := d.GetOk("replica_of"); ok {
		if createOpts.Datastore.Type == PostgresPro {
			return diag.Errorf("replica_of field is forbidden for PostgresPro")
		}
		createOpts.ReplicaOf = replicaOf.(string)
	}

	if v, ok := d.GetOk("network"); ok {
		createOpts.Nics, err = extractDatabaseNetworks(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s network", message)
		}
	}

	if v, ok := d.GetOk("disk_autoexpand"); ok {
		autoExpandOpts, err := extractDatabaseAutoExpand(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s disk_autoexpand", message)
		}
		var autoexpand int
		if autoExpandOpts.AutoExpand {
			autoexpand = 1
		} else {
			autoexpand = 0
		}
		createOpts.AutoExpand = &autoexpand
		createOpts.MaxDiskSize = autoExpandOpts.MaxDiskSize
	}

	if v, ok := d.GetOk("wal_volume"); ok {
		walVolumeOpts, err := extractDatabaseWalVolume(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s wal_volume", message)
		}
		createOpts.Walvolume = &walVolume{
			Size:       &walVolumeOpts.Size,
			VolumeType: walVolumeOpts.VolumeType,
		}
		walAutoExpandOpts, err := extractDatabaseAutoExpand(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s wal_disk_autoexpand", message)
		}
		var walAutoexpand int
		if walAutoExpandOpts.AutoExpand {
			walAutoexpand = 1
		} else {
			walAutoexpand = 0
		}
		createOpts.Walvolume.AutoExpand = walAutoexpand
		createOpts.Walvolume.MaxDiskSize = walAutoExpandOpts.MaxDiskSize
	}

	if v, ok := d.GetOk("wal_disk_autoexpand"); ok {
		walAutoExpandOpts, err := extractDatabaseAutoExpand(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s wal_disk_autoexpand", message)
		}
		var walAutoexpand int
		if walAutoExpandOpts.AutoExpand {
			walAutoexpand = 1
		} else {
			walAutoexpand = 0
		}
		createOpts.Walvolume.AutoExpand = walAutoexpand
		createOpts.Walvolume.MaxDiskSize = walAutoExpandOpts.MaxDiskSize
	}

	if v, ok := d.GetOk("backup_schedule"); ok {
		backupSchedule, err := extractDatabaseBackupSchedule(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s backup_schedule", message)
		}
		createOpts.BackupSchedule = &backupSchedule
	}

	var checkCapabilities *[]instanceCapabilityOpts
	if capabilities, ok := d.GetOk("capabilities"); ok {
		capabilitiesOpts, err := extractDatabaseCapabilities(capabilities.([]interface{}))
		if err != nil {
			return diag.Errorf("%s capability", message)
		}
		createOpts.Capabilities = capabilitiesOpts
		checkCapabilities = &capabilitiesOpts
	} else {
		checkCapabilities = nil
	}

	log.Printf("[DEBUG] vkcs_db_instance create options: %#v", createOpts)

	inst := dbInstance{}
	inst.Instance = createOpts
	instance, err := instanceCreate(DatabaseV1Client, inst).extract()
	if err != nil {
		return diag.Errorf("error creating vkcs_db_instance: %s", err)
	}

	// Wait for the instance to become available.
	log.Printf("[DEBUG] Waiting for vkcs_db_instance %s to become available", instance.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(dbInstanceStatusBuild), string(dbInstanceStatusBackup)},
		Target:     []string{string(dbInstanceStatusActive)},
		Refresh:    databaseInstanceStateRefreshFunc(DatabaseV1Client, instance.ID, checkCapabilities),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      dbInstanceDelay,
		MinTimeout: dbInstanceMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for vkcs_db_instance %s to become ready: %s", instance.ID, err)
	}

	if configuration, ok := d.GetOk("configuration_id"); ok {
		log.Printf("[DEBUG] Attaching configuration %s to vkcs_db_instance %s", configuration, instance.ID)
		var attachConfigurationOpts instanceAttachConfigurationGroupOpts
		attachConfigurationOpts.Instance.Configuration = configuration.(string)
		err := instanceAttachConfigurationGroup(DatabaseV1Client, instance.ID, &attachConfigurationOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error attaching configuration group %s to vkcs_db_instance %s: %s",
				configuration, instance.ID, err)
		}
	}

	if rootEnabled, ok := d.GetOk("root_enabled"); ok {
		if rootEnabled.(bool) {
			rootPassword := d.Get("root_password")
			var rootUserEnableOpts instanceRootUserEnableOpts
			if rootPassword != "" {
				rootUserEnableOpts.Password = rootPassword.(string)
			}
			rootUser, err := instanceRootUserEnable(DatabaseV1Client, instance.ID, &rootUserEnableOpts).extract()
			if err != nil {
				return diag.Errorf("error creating root user for instance: %s: %s", instance.ID, err)
			}
			d.Set("root_password", rootUser.Password)
		}
	}

	// Store the ID now
	d.SetId(instance.ID)

	return resourceDatabaseInstanceRead(ctx, d, meta)
}

func resourceDatabaseInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	instance, err := instanceGet(DatabaseV1Client, d.Id()).extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_db_instance"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_db_instance %s: %#v", d.Id(), instance)

	d.Set("name", instance.Name)
	d.Set("flavor_id", instance.Flavor.ID)
	d.Set("datastore", flattenDatabaseInstanceDatastore(*instance.DataStore))
	if _, ok := d.GetOk("disk_autoexpand"); ok {
		d.Set("disk_autoexpand", flattenDatabaseInstanceAutoExpand(instance.AutoExpand, instance.MaxDiskSize))
	}
	d.Set("region", getRegion(d, config))
	d.Set("size", instance.Volume.Size)
	d.Set("configuration_id", instance.ConfigurationID)
	if instance.WalVolume != nil && instance.WalVolume.VolumeID != "" {
		var walVolumeType string
		if v, ok := d.GetOk("wal_volume"); ok {
			walV, _ := extractDatabaseWalVolume(v.([]interface{}))
			walVolumeType = walV.VolumeType
		}
		walvolume := walVolume{Size: instance.WalVolume.Size, VolumeType: walVolumeType}
		d.Set("wal_volume", flattenDatabaseInstanceWalVolume(walvolume))

		if _, ok := d.GetOk("wal_disk_autoexpand"); ok {
			d.Set("wal_disk_autoexpand", flattenDatabaseInstanceAutoExpand(instance.WalVolume.AutoExpand, instance.WalVolume.MaxDiskSize))
		}
	}
	if instance.ReplicaOf != nil {
		d.Set("replica_of", instance.ReplicaOf.ID)
	} else {
		isRootEnabledResult := instanceRootUserGet(DatabaseV1Client, d.Id())
		isRootEnabled, err := isRootEnabledResult.extract()
		if err != nil {
			return diag.Errorf("error checking if root user is enabled for instance: %s: %s", d.Id(), err)
		}
		if isRootEnabled {
			d.Set("root_enabled", true)
		}
	}

	backupSchedule, err := instanceGetBackupSchedule(DatabaseV1Client, d.Id()).extract()
	if err != nil {
		return diag.Errorf("error getting backup schedule for instance: %s: %s", d.Id(), err)
	}
	if backupSchedule != nil {
		flattened := flattenDatabaseBackupSchedule(*backupSchedule)
		d.Set("backup_schedule", flattened)
	} else {
		d.Set("backup_schedule", nil)
	}

	d.Set("ip", instance.IP)

	if _, ok := d.GetOk("replica_of"); !ok {
		return nil
	}
	// Check if user set both "replica_of" and "network.fixed_ip_v4"
	var diags diag.Diagnostics

	rawNetworks := d.Get("network").([]interface{})
	for i, n := range rawNetworks {
		rawNetwork := n.(map[string]interface{})
		rawPath := fmt.Sprintf("network.%d.fixed_ip_v4", i)
		if v := rawNetwork["fixed_ip_v4"].(string); v != "" && d.HasChange(rawPath) {
			path := cty.Path{
				cty.GetAttrStep{Name: "network"},
				cty.IndexStep{Key: cty.NumberIntVal(int64(i))},
				cty.GetAttrStep{Name: "fixed_ip_v4"},
			}
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Field conflicts with attribute \"replica_of\".",
				Detail: "Setting \"fixed_ip_v4\" and \"replica_of\" at the same time " +
					"causes the \"fixed_ip_v4\" field to be ignored.",
				AttributePath: path,
			})
		}
	}

	return diags
}

func resourceDatabaseInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(dbInstanceStatusBuild)},
		Target:     []string{string(dbInstanceStatusActive)},
		Refresh:    databaseInstanceStateRefreshFunc(DatabaseV1Client, d.Id(), nil),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      dbInstanceDelay,
		MinTimeout: dbInstanceMinTimeout,
	}

	if d.HasChange("configuration_id") {
		old, new := d.GetChange("configuration_id")

		err := instanceDetachConfigurationGroup(DatabaseV1Client, d.Id()).ExtractErr()
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("Detaching configuration %s from vkcs_db_instance %s", old, d.Id())

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_db_instance %s to become ready: %s", d.Id(), err)
		}

		if new != "" {
			var attachConfigurationOpts instanceAttachConfigurationGroupOpts
			attachConfigurationOpts.Instance.Configuration = new.(string)
			err := instanceAttachConfigurationGroup(DatabaseV1Client, d.Id(), &attachConfigurationOpts).ExtractErr()
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("Attaching configuration %s to vkcs_db_instance %s", new, d.Id())

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("error waiting for vkcs_db_instance %s to become ready: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("size") {
		_, new := d.GetChange("size")
		var resizeVolumeOpts instanceResizeVolumeOpts
		resizeVolumeOpts.Resize.Volume.Size = new.(int)
		err := instanceAction(DatabaseV1Client, d.Id(), &resizeVolumeOpts).ExtractErr()
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("Resizing volume from vkcs_db_instance %s", d.Id())

		stateConf.Pending = []string{string(dbInstanceStatusResize)}
		stateConf.Target = []string{string(dbInstanceStatusActive)}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_db_instance %s to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("flavor_id") {
		var resizeOpts instanceResizeOpts
		resizeOpts.Resize.FlavorRef = d.Get("flavor_id").(string)
		err := instanceAction(DatabaseV1Client, d.Id(), &resizeOpts).ExtractErr()
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("Resizing flavor from vkcs_db_instance %s", d.Id())

		stateConf.Pending = []string{string(dbInstanceStatusResize)}
		stateConf.Target = []string{string(dbInstanceStatusActive)}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_db_instance %s to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("replica_of") {
		old, new := d.GetChange("replica_of")
		if old != "" && new == "" {
			detachReplicaOpts := &instanceDetachReplicaOpts{}
			detachReplicaOpts.Instance.ReplicaOf = old.(string)
			err := instanceDetachReplica(DatabaseV1Client, d.Id(), detachReplicaOpts).ExtractErr()
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("Detach replica from vkcs_db_instance %s", d.Id())

			stateConf.Pending = []string{string(dbInstanceStatusDetach)}
			stateConf.Target = []string{string(dbInstanceStatusActive)}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("error waiting for vkcs_db_instance %s to become ready: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("root_enabled") {
		_, new := d.GetChange("root_enabled")
		if new == true {
			rootPassword := d.Get("root_password")
			var rootUserEnableOpts instanceRootUserEnableOpts
			if rootPassword != "" {
				rootUserEnableOpts.Password = rootPassword.(string)
			}

			rootUser, err := instanceRootUserEnable(DatabaseV1Client, d.Id(), &rootUserEnableOpts).extract()
			if err != nil {
				return diag.Errorf("error creating root user for instance: %s: %s", d.Id(), err)
			}
			d.Set("root_password", rootUser.Password)
		} else {
			err = instanceRootUserDisable(DatabaseV1Client, d.Id()).ExtractErr()
			if err != nil {
				return diag.Errorf("error deleting root_user for instance %s: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("disk_autoexpand") {
		_, new := d.GetChange("disk_autoexpand")
		autoExpandProperties, err := extractDatabaseAutoExpand(new.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_instance disk_autoexpand")
		}
		var autoExpandOpts instanceUpdateAutoExpandOpts
		if autoExpandProperties.AutoExpand {
			autoExpandOpts.Instance.VolumeAutoresizeEnabled = 1
		} else {
			autoExpandOpts.Instance.VolumeAutoresizeEnabled = 0
		}
		autoExpandOpts.Instance.VolumeAutoresizeMaxSize = autoExpandProperties.MaxDiskSize
		err = instanceUpdateAutoExpand(DatabaseV1Client, d.Id(), &autoExpandOpts).ExtractErr()
		if err != nil {
			return diag.FromErr(err)
		}

		stateConf.Pending = []string{string(dbInstanceStatusBuild)}
		stateConf.Target = []string{string(dbInstanceStatusActive)}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_db_instance %s to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("wal_volume") {
		old, new := d.GetChange("wal_volume")
		walVolumeOptsNew, err := extractDatabaseWalVolume(new.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_instance wal_volume")
		}

		walVolumeOptsOld, err := extractDatabaseWalVolume(old.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_instance wal_volume")
		}

		if walVolumeOptsNew.Size != walVolumeOptsOld.Size {
			var resizeWalVolumeOpts instanceResizeWalVolumeOpts
			resizeWalVolumeOpts.Resize.Volume.Size = walVolumeOptsNew.Size
			resizeWalVolumeOpts.Resize.Volume.Kind = "wal"
			err = instanceAction(DatabaseV1Client, d.Id(), &resizeWalVolumeOpts).ExtractErr()
			if err != nil {
				return diag.FromErr(err)
			}

			stateConf.Pending = []string{string(dbInstanceStatusResize)}
			stateConf.Target = []string{string(dbInstanceStatusActive)}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("error waiting for vkcs_db_instance %s to become ready: %s", d.Id(), err)
			}
		}

	}

	if d.HasChange("wal_disk_autoexpand") {
		_, new := d.GetChange("wal_disk_autoexpand")
		walAutoExpandProperties, err := extractDatabaseAutoExpand(new.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_instance wal_disk_autoexpand")
		}
		var walAutoExpandOpts instanceUpdateAutoExpandWalOpts
		if walAutoExpandProperties.AutoExpand {
			walAutoExpandOpts.Instance.WalVolume.VolumeAutoresizeEnabled = 1
		} else {
			walAutoExpandOpts.Instance.WalVolume.VolumeAutoresizeEnabled = 0
		}
		walAutoExpandOpts.Instance.WalVolume.VolumeAutoresizeMaxSize = walAutoExpandProperties.MaxDiskSize
		err = instanceUpdateAutoExpand(DatabaseV1Client, d.Id(), &walAutoExpandOpts).ExtractErr()
		if err != nil {
			return diag.FromErr(err)
		}

		stateConf.Pending = []string{string(dbInstanceStatusBuild)}
		stateConf.Target = []string{string(dbInstanceStatusActive)}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_db_instance %s to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("capabilities") {
		_, newCapabilities := d.GetChange("capabilities")
		newCapabilitiesOpts, err := extractDatabaseCapabilities(newCapabilities.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_instance capability")
		}
		var applyCapabilityOpts instanceApplyCapabilityOpts
		applyCapabilityOpts.ApplyCapability.Capabilities = newCapabilitiesOpts

		err = instanceAction(DatabaseV1Client, d.Id(), &applyCapabilityOpts).ExtractErr()

		if err != nil {
			return diag.Errorf("error applying capability to vkcs_db_instance %s: %s", d.Id(), err)
		}

		applyCapabilityInstanceConf := &resource.StateChangeConf{
			Pending:    []string{string(dbInstanceStatusCapabilityApplying), string(dbInstanceStatusBuild)},
			Target:     []string{string(dbInstanceStatusActive)},
			Refresh:    databaseInstanceStateRefreshFunc(DatabaseV1Client, d.Id(), &newCapabilitiesOpts),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      dbInstanceDelay,
			MinTimeout: dbInstanceMinTimeout,
		}
		log.Printf("[DEBUG] Waiting for instance to become ready after applying capability")
		_, err = applyCapabilityInstanceConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error applying capability to vkcs_db_instance %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("backup_schedule") {
		_, newBackupSchedule := d.GetChange("backup_schedule")
		backupScheduleUpdateOpts, err := extractDatabaseBackupSchedule(newBackupSchedule.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_instance backup_schedule")
		}

		err = instanceUpdateBackupSchedule(DatabaseV1Client, d.Id(), &backupScheduleUpdateOpts).ExtractErr()

		if err != nil {
			return diag.Errorf("error updating backup schedule for vkcs_db_instance %s: %s", d.Id(), err)
		}

		stateConf.Pending = []string{string(dbInstanceStatusBuild), string(dbInstanceStatusBackup)}
		stateConf.Target = []string{string(dbInstanceStatusActive)}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_db_instance %s to become ready: %s", d.Id(), err)
		}
	}

	return resourceDatabaseInstanceRead(ctx, d, meta)
}

func resourceDatabaseInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	err = instanceDelete(DatabaseV1Client, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_db_instance"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(dbInstanceStatusActive), string(dbInstanceStatusShutdown)},
		Target:     []string{string(dbInstanceStatusDeleted)},
		Refresh:    databaseInstanceStateRefreshFunc(DatabaseV1Client, d.Id(), nil),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      dbInstanceDelay,
		MinTimeout: dbInstanceMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for vkcs_db_instance %s to delete: %s", d.Id(), err)
	}

	return nil
}
