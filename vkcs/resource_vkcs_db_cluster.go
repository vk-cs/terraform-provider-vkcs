package vkcs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type dbClusterStatus string

var (
	dbClusterStatusActive             dbClusterStatus = "CLUSTER_ACTIVE"
	dbClusterStatusBuild              dbClusterStatus = "BUILDING"
	dbClusterStatusDeleted            dbClusterStatus = "DELETED"
	dbClusterStatusDeleting           dbClusterStatus = "DELETING"
	dbClusterStatusGrow               dbClusterStatus = "GROWING_CLUSTER"
	dbClusterStatusResize             dbClusterStatus = "RESIZING_CLUSTER"
	dbClusterStatusShrink             dbClusterStatus = "SHRINKING_CLUSTER"
	dbClusterStatusUpdating           dbClusterStatus = "UPDATING_CLUSTER"
	dbClusterStatusCapabilityApplying dbClusterStatus = "CAPABILITY_APPLYING"
	dbClusterStatusBackup             dbClusterStatus = "BACKUP"
)

const (
	dbClusterInstanceRoleLeader string = "leader"
)

func resourceDatabaseCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseClusterCreate,
		ReadContext:   resourceDatabaseClusterRead,
		DeleteContext: resourceDatabaseClusterDelete,
		UpdateContext: resourceDatabaseClusterUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				config := meta.(configer)
				DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
				if err != nil {
					return nil, fmt.Errorf("error creating VKCS database client: %s", err)
				}

				if resourceDatabaseClusterRead(ctx, d, meta).HasError() {
					return nil, fmt.Errorf("error reading vkcs_cluster")
				}

				capabilities, err := clusterGetCapabilities(DatabaseV1Client, d.Id()).extract()
				if err != nil {
					return nil, fmt.Errorf("error getting cluster capabilities")
				}
				d.Set("capabilities", flattenDatabaseInstanceCapabilities(capabilities))
				d.Set("volume_type", dbImportedStatus)
				if v, ok := d.GetOk("wal_volume"); ok {
					walV, _ := extractDatabaseWalVolume(v.([]interface{}))
					walvolume := walVolume{Size: &walV.Size, VolumeType: dbImportedStatus}
					d.Set("wal_volume", flattenDatabaseClusterWalVolume(walvolume))
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
				Description: "The name of the cluster. Changing this creates a new cluster.",
			},

			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Computed:    false,
				Description: "The ID of flavor for the cluster.",
			},

			"cluster_size": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    false,
				Description: "The number of instances in the cluster.",
			},

			"volume_size": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    false,
				Computed:    false,
				Description: "Size of the cluster instance volume.",
			},

			"volume_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Computed:    false,
				Description: "The type of the cluster instance volume. Changing this creates a new cluster.",
			},

			"wal_volume": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
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
							ForceNew:    false,
							Description: "The type of the cluster wal volume. Changing this creates a new cluster.",
						},
					},
				},
				Description: "Object that represents wal volume of the cluster. Changing this creates a new cluster.",
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
							Description: "Version of the datastore. Changing this creates a new cluster.",
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(getClusterDatastores(), true),
							Description:  "Type of the datastore. Changing this creates a new cluster. Type of the datastore can either be \"galera_mysql\", \"postgresql\" or \"tarantool\".",
						},
					},
				},
				Description: "Object that represents datastore of the cluster. Changing this creates a new cluster.",
			},

			"loadbalancer_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the loadbalancer attached to the cluster.",
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
							Description: "The id of the network. Changing this creates a new cluster.",
						},
						"port": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The port id of the network. Changing this creates a new cluster.",
						},
					},
				},
				Description: "Object that represents network of the cluster. Changing this creates a new cluster.",
			},

			"configuration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    false,
				Description: "The id of the configuration attached to cluster.",
			},

			"root_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Description: "Indicates whether root user is enabled for the cluster.",
			},

			"root_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Computed:    true,
				ForceNew:    false,
				Description: "Password for the root user of the cluster.",
			},

			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The name of the availability zone of the cluster. Changing this creates a new cluster.",
			},

			"floating_ip_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Indicates whether floating ip is created for cluster. Changing this creates a new cluster.",
			},

			"keypair": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Name of the keypair to be attached to cluster. Changing this creates a new cluster.",
			},

			"disk_autoexpand": {
				Type:     schema.TypeList,
				Optional: true,
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
				Description: "Object that represents autoresize properties of the cluster.",
			},

			"wal_disk_autoexpand": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"autoexpand": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    false,
							Description: "Indicates whether wal volume autoresize is enabled.",
						},
						"max_disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    false,
							Description: "Maximum disk size for wal volume autoresize.",
						},
					},
				},
				Description: "Object that represents autoresize properties of wal volume of the cluster.",
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
				Description: "Object that represents capability applied to cluster. There can be several instances of this object.",
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
							Description: "Used only for restoring from PITR backups. Timestamp of needed backup in format \"2021-10-06 01:02:00\". You can specify \"latest\" to use most recent backup.",
						},
					},
				},
				Description: "Object that represents backup to restore cluster from. **New since v.0.1.4**.",
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
				Description: "Object that represents configuration of PITR backup. This functionality is available only for postgres datastore. **New since v.0.1.4**",
			},

			"shrink_options": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Used only for shrinking cluster. List of IDs of instances that should remain after shrink. If no options are supplied, shrink operation will choose first non-leader instance to delete.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
			},

			// Computed values
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the instance.",
						},
						"ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "IP address of the instance.",
						},
						"role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The role of the instance in cluster.",
						},
					},
				},
				Description: "Cluster instances info.",
			},
		},
		Description: "Provides a db cluster resource. This can be used to create, modify and delete db cluster for galera_mysql, postgresql, tarantool datastores.",
	}
}

func resourceDatabaseClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	createOpts := &dbClusterCreateOpts{
		Name:              d.Get("name").(string),
		FloatingIPEnabled: d.Get("floating_ip_enabled").(bool),
	}

	message := "unable to determine vkcs_db_cluster"

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

	if v, ok := d.GetOk("disk_autoexpand"); ok {
		autoExpandOpts, err := extractDatabaseAutoExpand(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s disk_autoexpand", message)
		}
		if autoExpandOpts.AutoExpand {
			createOpts.AutoExpand = 1
		} else {
			createOpts.AutoExpand = 0
		}
		createOpts.MaxDiskSize = autoExpandOpts.MaxDiskSize
	}

	if v, ok := d.GetOk("wal_disk_autoexpand"); ok {
		walAutoExpandOpts, err := extractDatabaseAutoExpand(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s wal_disk_autoexpand", message)
		}
		if walAutoExpandOpts.AutoExpand {
			createOpts.WalAutoExpand = 1
		} else {
			createOpts.WalAutoExpand = 0
		}
		createOpts.WalMaxDiskSize = walAutoExpandOpts.MaxDiskSize
	}

	clusterSize := d.Get("cluster_size").(int)
	instances := make([]dbClusterInstanceCreateOpts, clusterSize)
	volumeSize := d.Get("volume_size").(int)
	createDBInstanceOpts := dbClusterInstanceCreateOpts{
		Keypair:          d.Get("keypair").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		FlavorRef:        d.Get("flavor_id").(string),
		Volume:           &volume{Size: &volumeSize, VolumeType: d.Get("volume_type").(string)},
	}

	if v, ok := d.GetOk("network"); ok {
		createDBInstanceOpts.Nics, err = extractDatabaseNetworks(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s network", message)
		}
	}

	if v, ok := d.GetOk("wal_volume"); ok {
		walVolumeOpts, err := extractDatabaseWalVolume(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s wal_volume", message)
		}
		createDBInstanceOpts.Walvolume = &walVolume{
			Size:       &walVolumeOpts.Size,
			VolumeType: walVolumeOpts.VolumeType,
		}
	}

	for i := 0; i < clusterSize; i++ {
		instances[i] = createDBInstanceOpts
	}

	createOpts.Instances = instances

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

	log.Printf("[DEBUG] vkcs_db_cluster create options: %#v", createOpts)
	clust := dbCluster{}
	clust.Cluster = createOpts

	cluster, err := dbClusterCreate(DatabaseV1Client, clust).extract()
	if err != nil {
		return diag.Errorf("error creating vkcs_db_cluster: %s", err)
	}

	// Wait for the cluster to become available.
	log.Printf("[DEBUG] Waiting for vkcs_db_cluster %s to become available", cluster.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(dbClusterStatusBuild), string(dbClusterStatusBackup)},
		Target:     []string{string(dbClusterStatusActive)},
		Refresh:    databaseClusterStateRefreshFunc(DatabaseV1Client, cluster.ID, checkCapabilities),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      dbInstanceDelay,
		MinTimeout: dbInstanceMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for vkcs_db_cluster %s to become ready: %s", cluster.ID, err)
	}

	if configuration, ok := d.GetOk("configuration_id"); ok {
		log.Printf("[DEBUG] Attaching configuration %s to vkcs_db_cluster %s", configuration, cluster.ID)
		var attachConfigurationOpts dbClusterAttachConfigurationGroupOpts
		attachConfigurationOpts.ConfigurationAttach.ConfigurationID = configuration.(string)
		err := instanceAttachConfigurationGroup(DatabaseV1Client, cluster.ID, &attachConfigurationOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error attaching configuration group %s to vkcs_db_cluster %s: %s",
				configuration, cluster.ID, err)
		}
	}

	// Store the ID now
	d.SetId(cluster.ID)
	return resourceDatabaseClusterRead(ctx, d, meta)
}

func resourceDatabaseClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	cluster, err := dbClusterGet(DatabaseV1Client, d.Id()).extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_db_cluster"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_db_cluster %s: %#v", d.Id(), cluster)

	d.Set("name", cluster.Name)
	d.Set("datastore", flattenDatabaseInstanceDatastore(*cluster.DataStore))
	d.Set("flavor_id", cluster.Instances[0].Flavor.ID)
	d.Set("cluster_size", len(cluster.Instances))
	d.Set("volume_size", cluster.Instances[0].Volume.Size)
	d.Set("loadbalancer_id", cluster.LoadbalancerID)

	d.Set("configuration_id", cluster.ConfigurationID)
	if _, ok := d.GetOk("disk_autoexpand"); ok {
		d.Set("disk_autoexpand", flattenDatabaseInstanceAutoExpand(cluster.AutoExpand, cluster.MaxDiskSize))
	}
	if cluster.Instances[0].WalVolume != nil && cluster.Instances[0].WalVolume.VolumeID != "" {
		var walVolumeType string
		if v, ok := d.GetOk("wal_volume"); ok {
			walV, _ := extractDatabaseWalVolume(v.([]interface{}))
			walVolumeType = walV.VolumeType
		}
		walvolume := walVolume{Size: cluster.Instances[0].WalVolume.Size, VolumeType: walVolumeType}
		d.Set("wal_volume", flattenDatabaseClusterWalVolume(walvolume))

		if _, ok := d.GetOk("wal_disk_autoexpand"); ok {
			d.Set("wal_disk_autoexpand", flattenDatabaseInstanceAutoExpand(cluster.WalAutoExpand, cluster.WalMaxDiskSize))
		}
	}

	d.Set("instances", flattenDatabaseClusterInstances(cluster.Instances))

	backupSchedule, err := dbClusterGetBackupSchedule(DatabaseV1Client, d.Id()).extract()
	if err != nil {
		return diag.Errorf("error getting backup schedule for cluster: %s: %s", d.Id(), err)
	}
	if backupSchedule != nil {
		flattened := flattenDatabaseBackupSchedule(*backupSchedule)
		d.Set("backup_schedule", flattened)
	} else {
		d.Set("backup_schedule", nil)
	}

	return nil
}

func resourceDatabaseClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(dbClusterStatusBuild)},
		Target:     []string{string(dbClusterStatusActive)},
		Refresh:    databaseClusterStateRefreshFunc(DatabaseV1Client, d.Id(), nil),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      dbInstanceDelay,
		MinTimeout: dbInstanceMinTimeout,
	}

	if d.HasChange("configuration_id") {
		old, new := d.GetChange("configuration_id")

		var detachConfigurationOpts dbClusterDetachConfigurationGroupOpts
		detachConfigurationOpts.ConfigurationDetach.ConfigurationID = old.(string)
		err := dbClusterAction(DatabaseV1Client, d.Id(), &detachConfigurationOpts).ExtractErr()
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("Detaching configuration %s from vkcs_db_cluster %s", old, d.Id())

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_db_cluster %s to become ready: %s", d.Id(), err)
		}

		if new != "" {
			var attachConfigurationOpts dbClusterAttachConfigurationGroupOpts
			attachConfigurationOpts.ConfigurationAttach.ConfigurationID = new.(string)
			err := dbClusterAction(DatabaseV1Client, d.Id(), &attachConfigurationOpts).ExtractErr()
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("Attaching configuration %s to vkcs_db_cluster %s", new, d.Id())

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("error waiting for vkcs_db_cluster %s to become ready: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("volume_size") {
		_, new := d.GetChange("volume_size")
		var resizeVolumeOpts dbClusterResizeVolumeOpts
		resizeVolumeOpts.Resize.Volume.Size = new.(int)
		err := dbClusterAction(DatabaseV1Client, d.Id(), &resizeVolumeOpts).ExtractErr()
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("Resizing volume from vkcs_db_cluster %s", d.Id())

		stateConf.Pending = []string{string(dbClusterStatusResize)}
		stateConf.Target = []string{string(dbClusterStatusActive)}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_db_cluster %s to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("flavor_id") {
		var resizeOpts dbClusterResizeOpts
		resizeOpts.Resize.FlavorRef = d.Get("flavor_id").(string)
		err := dbClusterAction(DatabaseV1Client, d.Id(), &resizeOpts).ExtractErr()
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("Resizing flavor from vkcs_db_cluster %s", d.Id())

		stateConf.Pending = []string{string(dbClusterStatusResize)}
		stateConf.Target = []string{string(dbClusterStatusActive)}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_db_cluster %s to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("disk_autoexpand") {
		_, new := d.GetChange("disk_autoexpand")
		autoExpandProperties, err := extractDatabaseAutoExpand(new.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_cluster disk_autoexpand")
		}
		var autoExpandOpts dbClusterUpdateAutoExpandOpts
		if autoExpandProperties.AutoExpand {
			autoExpandOpts.Cluster.VolumeAutoresizeEnabled = 1
		} else {
			autoExpandOpts.Cluster.VolumeAutoresizeEnabled = 0
		}
		autoExpandOpts.Cluster.VolumeAutoresizeMaxSize = autoExpandProperties.MaxDiskSize
		err = dbClusterUpdateAutoExpand(DatabaseV1Client, d.Id(), &autoExpandOpts).ExtractErr()
		if err != nil {
			return diag.FromErr(err)
		}

		stateConf.Pending = []string{string(dbClusterStatusUpdating)}
		stateConf.Target = []string{string(dbClusterStatusActive)}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_db_cluster %s to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("wal_volume") {
		old, new := d.GetChange("wal_volume")
		walVolumeOptsNew, err := extractDatabaseWalVolume(new.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_cluster wal_volume")
		}

		walVolumeOptsOld, err := extractDatabaseWalVolume(old.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_cluster wal_volume")
		}

		if walVolumeOptsNew.Size != walVolumeOptsOld.Size {
			var resizeWalVolumeOpts dbClusterResizeWalVolumeOpts
			resizeWalVolumeOpts.Resize.Volume.Size = walVolumeOptsNew.Size
			resizeWalVolumeOpts.Resize.Volume.Kind = "wal"
			err = dbClusterAction(DatabaseV1Client, d.Id(), &resizeWalVolumeOpts).ExtractErr()
			if err != nil {
				return diag.FromErr(err)
			}

			stateConf.Pending = []string{string(dbClusterStatusResize)}
			stateConf.Target = []string{string(dbClusterStatusActive)}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("error waiting for vkcs_db_cluster %s to become ready: %s", d.Id(), err)
			}
		}

	}
	if d.HasChange("wal_disk_autoexpand") {
		_, new := d.GetChange("wal_disk_autoexpand")
		walAutoExpandProperties, err := extractDatabaseAutoExpand(new.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_cluster wal_disk_autoexpand")
		}
		var walAutoExpandOpts dbClusterUpdateAutoExpandWalOpts
		if walAutoExpandProperties.AutoExpand {
			walAutoExpandOpts.Cluster.WalVolume.VolumeAutoresizeEnabled = 1
		} else {
			walAutoExpandOpts.Cluster.WalVolume.VolumeAutoresizeEnabled = 0
		}
		walAutoExpandOpts.Cluster.WalVolume.VolumeAutoresizeMaxSize = walAutoExpandProperties.MaxDiskSize
		err = dbClusterUpdateAutoExpand(DatabaseV1Client, d.Id(), &walAutoExpandOpts).ExtractErr()
		if err != nil {
			return diag.FromErr(err)
		}

		stateConf.Pending = []string{string(dbClusterStatusUpdating)}
		stateConf.Target = []string{string(dbClusterStatusActive)}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_db_cluster %s to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("capabilities") {
		_, newCapabilities := d.GetChange("capabilities")
		newCapabilitiesOpts, err := extractDatabaseCapabilities(newCapabilities.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_cluster capability")
		}
		var applyCapabilityOpts dbClusterApplyCapabilityOpts
		applyCapabilityOpts.ApplyCapability.Capabilities = newCapabilitiesOpts

		err = dbClusterAction(DatabaseV1Client, d.Id(), &applyCapabilityOpts).ExtractErr()

		if err != nil {
			return diag.Errorf("error applying capability to vkcs_db_cluster %s: %s", d.Id(), err)
		}

		applyCapabilityClusterConf := &resource.StateChangeConf{
			Pending:    []string{string(dbClusterStatusCapabilityApplying), string(dbClusterStatusBuild)},
			Target:     []string{string(dbClusterStatusActive)},
			Refresh:    databaseClusterStateRefreshFunc(DatabaseV1Client, d.Id(), &newCapabilitiesOpts),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      dbInstanceDelay,
			MinTimeout: dbInstanceMinTimeout,
		}
		log.Printf("[DEBUG] Waiting for cluster to become ready after applying capability")
		_, err = applyCapabilityClusterConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error applying capability to vkcs_db_cluster %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("cluster_size") {
		old, new := d.GetChange("cluster_size")
		if new.(int) > old.(int) {
			opts := make([]dbClusterGrowOpts, new.(int)-old.(int))

			volumeSize := d.Get("volume_size").(int)
			growOpts := dbClusterGrowOpts{
				Keypair:          d.Get("keypair").(string),
				AvailabilityZone: d.Get("availability_zone").(string),
				FlavorRef:        d.Get("flavor_id").(string),
				Volume:           &volume{Size: &volumeSize, VolumeType: d.Get("volume_type").(string)},
			}
			if v, ok := d.GetOk("wal_volume"); ok {
				walVolumeOpts, err := extractDatabaseWalVolume(v.([]interface{}))
				if err != nil {
					return diag.Errorf("unable to determine vkcs_db_cluster wal_volume")
				}
				growOpts.Walvolume = &walVolume{
					Size:       &walVolumeOpts.Size,
					VolumeType: walVolumeOpts.VolumeType,
				}
			}
			for i := 0; i < len(opts); i++ {
				opts[i] = growOpts
			}
			growClusterOpts := dbClusterGrowClusterOpts{
				Grow: opts,
			}
			err = dbClusterAction(DatabaseV1Client, d.Id(), &growClusterOpts).ExtractErr()

			if err != nil {
				return diag.Errorf("error growing vkcs_db_cluster %s: %s", d.Id(), err)
			}
			stateConf.Pending = []string{string(dbClusterStatusGrow)}
			stateConf.Target = []string{string(dbClusterStatusActive)}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("error waiting for vkcs_db_cluster %s to become ready: %s", d.Id(), err)
			}
		} else {
			rawShrinkOptions := d.Get("shrink_options").([]interface{})
			shrinkOptions := expandDatabaseClusterShrinkOptions(rawShrinkOptions)
			if len(shrinkOptions) > 0 && len(shrinkOptions) != new.(int) {
				return diag.Errorf("invalid shrink options: number of instances in shrink options should equal new cluster size")
			}

			cluster, err := dbClusterGet(DatabaseV1Client, d.Id()).extract()
			if err != nil {
				return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_db_cluster"))
			}

			toDelete := old.(int) - new.(int)
			ids, err := databaseClusterDetermineShrinkedInstances(toDelete, shrinkOptions, cluster.Instances)
			if err != nil {
				return diag.Errorf("error determining instances to shrink: %s", err)
			}

			shrinkClusterOpts := dbClusterShrinkClusterOpts{
				Shrink: ids,
			}

			err = dbClusterAction(DatabaseV1Client, d.Id(), &shrinkClusterOpts).ExtractErr()

			if err != nil {
				return diag.Errorf("error growing vkcs_db_cluster %s: %s", d.Id(), err)
			}
			stateConf.Pending = []string{string(dbClusterStatusShrink)}
			stateConf.Target = []string{string(dbClusterStatusActive)}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("error waiting for vkcs_db_cluster %s to become ready: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("backup_schedule") {
		_, newBackupSchedule := d.GetChange("backup_schedule")
		backupScheduleUpdateOpts, err := extractDatabaseBackupSchedule(newBackupSchedule.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_cluster backup_schedule")
		}

		err = dbClusterUpdateBackupSchedule(DatabaseV1Client, d.Id(), &backupScheduleUpdateOpts).ExtractErr()

		if err != nil {
			return diag.Errorf("error updating backup schedule for vkcs_db_cluster %s: %s", d.Id(), err)
		}

		stateConf.Pending = []string{string(dbClusterStatusUpdating), string(dbClusterStatusBackup)}
		stateConf.Target = []string{string(dbClusterStatusActive)}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_db_cluster %s to become ready: %s", d.Id(), err)
		}
	}

	return resourceDatabaseClusterRead(ctx, d, meta)
}

func resourceDatabaseClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	err = dbClusterDelete(DatabaseV1Client, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_db_cluster"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(dbClusterStatusActive), string(dbClusterStatusDeleting)},
		Target:     []string{string(dbClusterStatusDeleted)},
		Refresh:    databaseClusterStateRefreshFunc(DatabaseV1Client, d.Id(), nil),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      dbInstanceDelay,
		MinTimeout: dbInstanceMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for vkcs_db_cluster %s to delete: %s", d.Id(), err)
	}

	return nil
}
