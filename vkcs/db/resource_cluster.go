package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/instances"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
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
	dbClusterStatusError              dbClusterStatus = "ERROR"
)

const (
	DBClusterInstanceRoleLeader string = "leader"
)

func ResourceDatabaseCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseClusterCreate,
		ReadContext:   resourceDatabaseClusterRead,
		DeleteContext: resourceDatabaseClusterDelete,
		UpdateContext: resourceDatabaseClusterUpdate,
		CustomizeDiff: resourceDatabaseCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				config := meta.(clients.Config)
				DatabaseV1Client, err := config.DatabaseV1Client(util.GetRegion(d, config))
				if err != nil {
					return nil, fmt.Errorf("error creating VKCS database client: %s", err)
				}

				if resourceDatabaseClusterRead(ctx, d, meta).HasError() {
					return nil, fmt.Errorf("error reading vkcs_cluster")
				}

				capabilities, err := clusters.GetCapabilities(DatabaseV1Client, d.Id()).Extract()
				if err != nil {
					return nil, fmt.Errorf("error getting cluster capabilities")
				}
				d.Set("capabilities", flattenDatabaseInstanceCapabilities(capabilities))
				d.Set("volume_type", dbImportedStatus)
				if v, ok := d.GetOk("wal_volume"); ok {
					walV, _ := extractDatabaseWalVolume(v.([]interface{}))
					walvolume := instances.WalVolume{Size: &walV.Size, VolumeType: dbImportedStatus}
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
							Description: "The id of the network. Changing this creates a new cluster. _note_ Although this argument is marked as optional, it is actually required at the moment. Not setting a value for it may cause an error.",
						},
						"port": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The port id of the network. Changing this creates a new cluster.",
							Deprecated:  "This argument is deprecated, please do not use it.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The id of the subnet. Changing this creates a new cluster.",
						},
						"security_groups": {
							Type:        schema.TypeSet,
							Optional:    true,
							ForceNew:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "An array of one or more security group IDs to associate with the cluster instances. Changing this creates a new cluster.",
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
				Description: "Password for the root user of the cluster. When enabling root, password is autogenerated, use this field to obtain it.",
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
				Description: "Object that represents backup to restore cluster from.",
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
				Description: "Object that represents configuration of PITR backup. This functionality is available only for postgres datastore.",
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

			"cloud_monitoring_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Description: "Enable cloud monitoring for the cluster. Changing this for Redis or MongoDB creates a new instance.",
			},

			"vendor_options": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"restart_confirmed": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Boolean to confirm autorestart of the cluster's instances if it is required to apply configuration group changes.",
						},
					},
				},
				Description: "Map of additional vendor-specific options. Supported options are described below.",
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
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	createOpts := &clusters.CreateOpts{
		Name:                   d.Get("name").(string),
		FloatingIPEnabled:      d.Get("floating_ip_enabled").(bool),
		CloudMonitoringEnabled: d.Get("cloud_monitoring_enabled").(bool),
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
	clusterInstances := make([]clusters.InstanceCreateOpts, clusterSize)
	volumeSize := d.Get("volume_size").(int)
	createDBInstanceOpts := clusters.InstanceCreateOpts{
		Keypair:          d.Get("keypair").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		FlavorRef:        d.Get("flavor_id").(string),
		Volume:           &instances.Volume{Size: &volumeSize, VolumeType: d.Get("volume_type").(string)},
	}

	if v, ok := d.GetOk("network"); ok {
		createDBInstanceOpts.Nics, createDBInstanceOpts.SecurityGroups, err = extractDatabaseNetworks(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s network", message)
		}
	}

	if v, ok := d.GetOk("wal_volume"); ok {
		walVolumeOpts, err := extractDatabaseWalVolume(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s wal_volume", message)
		}
		createDBInstanceOpts.Walvolume = &instances.WalVolume{
			Size:       &walVolumeOpts.Size,
			VolumeType: walVolumeOpts.VolumeType,
		}
	}

	for i := 0; i < clusterSize; i++ {
		clusterInstances[i] = createDBInstanceOpts
	}

	createOpts.Instances = clusterInstances

	if v, ok := d.GetOk("backup_schedule"); ok {
		backupSchedule, err := extractDatabaseBackupSchedule(v.([]interface{}))
		if err != nil {
			return diag.Errorf("%s backup_schedule", message)
		}
		createOpts.BackupSchedule = &backupSchedule
	}

	var checkCapabilities *[]instances.CapabilityOpts
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
	clust := clusters.Cluster{}
	clust.Cluster = createOpts

	cluster, err := clusters.Create(DatabaseV1Client, clust).Extract()
	if err != nil {
		return diag.Errorf("error creating vkcs_db_cluster: %s", err)
	}

	// Store the ID now
	d.SetId(cluster.ID)

	// Wait for the cluster to become available.
	log.Printf("[DEBUG] Waiting for vkcs_db_cluster %s to become available", cluster.ID)

	stateConf := &retry.StateChangeConf{
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

		var attachConfigurationOpts clusters.AttachConfigurationGroupOpts
		vendorOptionsRaw := d.Get("vendor_options").(*schema.Set)
		if vendorOptionsRaw.Len() > 0 {
			vendorOptions := util.ExpandVendorOptions(vendorOptionsRaw.List())
			if v, ok := vendorOptions["restart_confirmed"]; ok && v.(bool) {
				restartConfirmed := true
				attachConfigurationOpts.ConfigurationAttach.RestartConfirmed = &restartConfirmed
			}
		}
		attachConfigurationOpts.ConfigurationAttach.ConfigurationID = configuration.(string)

		err := clusters.ClusterAction(DatabaseV1Client, cluster.ID, &attachConfigurationOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error attaching configuration group %s to vkcs_db_cluster %s: %s",
				configuration, cluster.ID, err)
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{string(dbClusterStatusUpdating)},
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
	}

	diags := make(diag.Diagnostics, 0)

	if rootEnabled, ok := d.GetOk("root_enabled"); ok {
		if rootEnabled.(bool) {
			updateCtx := &dbResourceUpdateContext{
				Ctx:       ctx,
				Client:    DatabaseV1Client,
				D:         d,
				StateConf: nil,
			}
			diags = append(diags, databaseClusterActionEnableRoot(updateCtx)...)
			if diags.HasError() {
				return diags
			}
		}
	}

	return append(diags, resourceDatabaseClusterRead(ctx, d, meta)...)
}

func resourceDatabaseClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	cluster, err := clusters.Get(DatabaseV1Client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error retrieving vkcs_db_cluster"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_db_cluster %s: %#v", d.Id(), cluster)

	d.Set("name", cluster.Name)
	d.Set("datastore", flattenDatabaseInstanceDatastore(*cluster.DataStore))
	d.Set("loadbalancer_id", cluster.LoadbalancerID)

	d.Set("configuration_id", cluster.ConfigurationID)
	if _, ok := d.GetOk("disk_autoexpand"); ok {
		d.Set("disk_autoexpand", flattenDatabaseInstanceAutoExpand(cluster.AutoExpand, cluster.MaxDiskSize))
	}

	if len(cluster.Instances) > 0 {
		d.Set("flavor_id", cluster.Instances[0].Flavor.ID)
		d.Set("cluster_size", len(cluster.Instances))
		d.Set("volume_size", cluster.Instances[0].Volume.Size)

		if cluster.Instances[0].WalVolume != nil && cluster.Instances[0].WalVolume.VolumeID != "" {
			var walVolumeType string
			if v, ok := d.GetOk("wal_volume"); ok {
				walV, _ := extractDatabaseWalVolume(v.([]interface{}))
				walVolumeType = walV.VolumeType
			}
			walvolume := instances.WalVolume{Size: cluster.Instances[0].WalVolume.Size, VolumeType: walVolumeType}
			d.Set("wal_volume", flattenDatabaseClusterWalVolume(walvolume))

			if _, ok := d.GetOk("wal_disk_autoexpand"); ok {
				d.Set("wal_disk_autoexpand", flattenDatabaseInstanceAutoExpand(cluster.WalAutoExpand, cluster.WalMaxDiskSize))
			}
		}

		d.Set("instances", flattenDatabaseClusterInstances(cluster.Instances))
	}

	backupSchedule, err := clusters.GetBackupSchedule(DatabaseV1Client, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("error getting backup schedule for cluster: %s: %s", d.Id(), err)
	}
	if backupSchedule != nil {
		flattened := flattenDatabaseBackupSchedule(*backupSchedule)
		d.Set("backup_schedule", flattened)
	} else {
		d.Set("backup_schedule", nil)
	}

	if !d.HasChangesExcept() {
		return nil
	}

	var diags diag.Diagnostics

	rawNetworks := d.Get("network").([]interface{})
	diags = checkDBNetworks(rawNetworks, cty.Path{}, diags)
	return diags
}

func resourceDatabaseClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	dbClient, err := config.DatabaseV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	clusterID := d.Id()
	stateConf := &retry.StateChangeConf{
		Pending:    []string{string(dbClusterStatusBuild)},
		Target:     []string{string(dbClusterStatusActive)},
		Refresh:    databaseClusterStateRefreshFunc(dbClient, d.Id(), nil),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      dbInstanceDelay,
		MinTimeout: dbInstanceMinTimeout,
	}
	updateCtx := &dbResourceUpdateContext{
		Ctx:       ctx,
		Client:    dbClient,
		D:         d,
		StateConf: stateConf,
	}

	if d.HasChange("configuration_id") {
		err = databaseClusterActionUpdateConfiguration(updateCtx)
		if err != nil {
			return databaseClusterUpdateProcessError(err, clusterID)
		}
	}

	if d.HasChange("volume_size") {
		err = databaseClusterActionResizeVolume(updateCtx, "")
		if err != nil {
			return databaseClusterUpdateProcessError(err, clusterID)
		}
	}

	if d.HasChange("flavor_id") {
		err = databaseClusterActionResizeFlavor(updateCtx, "")
		if err != nil {
			return databaseClusterUpdateProcessError(err, clusterID)
		}
	}

	if d.HasChange("disk_autoexpand") {
		err = databaseClusterUpdateDiskAutoexpand(updateCtx)
		if err != nil {
			return databaseClusterUpdateProcessError(err, clusterID)
		}
	}

	if d.HasChange("wal_volume") {
		err = databaseClusterActionResizeWalVolume(updateCtx, "")
		if err != nil {
			return databaseClusterUpdateProcessError(err, clusterID)
		}
	}

	if d.HasChange("wal_disk_autoexpand") {
		err = databaseClusterUpdateWalDiskAutoexpand(updateCtx)
		if err != nil {
			return databaseClusterUpdateProcessError(err, clusterID)
		}
	}

	if d.HasChange("capabilities") {
		err = databaseClusterActionApplyCapabilities(updateCtx)
		if err != nil {
			return databaseClusterUpdateProcessError(err, clusterID)
		}
	}

	if d.HasChange("cluster_size") {
		old, new := d.GetChange("cluster_size")
		if sizeChange := new.(int) - old.(int); sizeChange > 0 {
			err = databaseClusterActionGrow(updateCtx, "")
		} else if sizeChange < 0 {
			err = databaseClusterActionShrink(updateCtx, "")
		}
		if err != nil {
			return databaseClusterUpdateProcessError(err, clusterID)
		}
	}

	if d.HasChange("backup_schedule") {
		_, newBackupSchedule := d.GetChange("backup_schedule")
		backupScheduleUpdateOpts, err := extractDatabaseBackupSchedule(newBackupSchedule.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_cluster backup_schedule")
		}

		err = clusters.UpdateBackupSchedule(dbClient, clusterID, &backupScheduleUpdateOpts).ExtractErr()

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

	if d.HasChange("cloud_monitoring_enabled") {
		err = databaseClusterUpdateCloudMonitoring(updateCtx)
		if err != nil {
			return databaseClusterUpdateProcessError(err, clusterID)
		}
	}

	diags := make(diag.Diagnostics, 0)

	if d.HasChange("root_enabled") {
		_, new := d.GetChange("root_enabled")
		if new == true {
			err := databaseClusterActionEnableRoot(updateCtx)
			if err.HasError() {
				return err
			} else {
				diags = append(diags, err...)
			}
		} else {
			d.Set("root_enabled", false)
			d.Set("root_password", "")
		}
	}

	return append(diags, resourceDatabaseClusterRead(ctx, d, meta)...)
}

func resourceDatabaseClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	err = clusters.Delete(DatabaseV1Client, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_db_cluster"))
	}

	stateConf := &retry.StateChangeConf{
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

func databaseClusterUpdateProcessError(err error, clusterID string) diag.Diagnostics {
	baseErr := err
	if unwrappedErr := errors.Unwrap(err); unwrappedErr != nil {
		baseErr = unwrappedErr
	}

	newErrMsg := baseErr.Error()
	switch baseErr {
	case errDBClusterNotFound:
		newErrMsg = fmt.Sprintf("error retrieving vkcs_db_cluster %s", clusterID)
	case errDBClusterUpdateWait:
		newErrMsg = fmt.Sprintf("error waiting for vkcs_db_cluster %s to become ready", clusterID)

	case errDBClusterUpdateDiskAutoexpand:
		newErrMsg = fmt.Sprintf("error updating disk_autoexpand for vkcs_db_cluster %s", clusterID)
	case errDBClusterUpdateDiskAutoexpandExtract:
		newErrMsg = fmt.Sprintf("unable to determine disk_autoexpand from vkcs_db_cluster %s", clusterID)
	case errDBClusterUpdateWalDiskAutoexpand:
		newErrMsg = fmt.Sprintf("error updating wal_disk_autoexpand for vkcs_db_cluster %s", clusterID)
	case errDBClusterUpdateWalDiskAutoexpandExtract:
		newErrMsg = fmt.Sprintf("unable to determine wal_disk_autoexpand from vkcs_db_cluster %s", clusterID)
	case errDBClusterUpdateCloudMonitoring:
		newErrMsg = fmt.Sprintf("error updating cloud_monitoring_enabled for vkcs_db_cluster %s", clusterID)

	case errDBClusterActionUpdateConfiguration:
		newErrMsg = fmt.Sprintf("error updating configuration for vkcs_db_cluster %s", clusterID)
	case errDBClusterActionApplyCapabitilies:
		newErrMsg = fmt.Sprintf("error updating capabilities for vkcs_db_cluster %s", clusterID)
	case errDBClusterActionApplyCapabilitiesExtract:
		newErrMsg = fmt.Sprintf("error extracting capabilities for vkcs_db_cluster %s", clusterID)
	case errDBClusterActionResizeWalVolumeExtract:
		newErrMsg = fmt.Sprintf("unable to determine wal_volume from vkcs_db_cluster %s", clusterID)
	case errDBClusterActionGrow:
		newErrMsg = fmt.Sprintf("error growing vkcs_db_cluster %s", clusterID)
	case errDBClusterActionShrink:
		newErrMsg = fmt.Sprintf("error shrinking vkcs_db_cluster %s", clusterID)
	case errDBClusterActionShrinkWrongOptions:
		newErrMsg = fmt.Sprintf("invalid shrink options for vkcs_db_cluster %s", clusterID)
	case errDBClusterActionShrinkInstancesExtract:
		newErrMsg = fmt.Sprintf("error determining instances to shrink vkcs_db_cluster %s", clusterID)
	case errDBClusterActionResizeVolume:
		newErrMsg = fmt.Sprintf("error resizing volume for vkcs_db_cluster %s", clusterID)
	case errDBClusterActionResizeWalVolume:
		newErrMsg = fmt.Sprintf("error resizing wal_volume for vkcs_db_cluster %s", clusterID)
	case errDBClusterActionResizeFlavor:
		newErrMsg = fmt.Sprintf("error changing flavor for vkcs_db_cluster %s", clusterID)
	}

	errMsg := strings.Replace(err.Error(), baseErr.Error(), newErrMsg, 1)
	return diag.Errorf(errMsg)
}
