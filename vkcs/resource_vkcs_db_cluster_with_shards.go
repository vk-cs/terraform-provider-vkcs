package vkcs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabaseClusterWithShards() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseClusterWithShardsCreate,
		ReadContext:   resourceDatabaseClusterWithShardsRead,
		DeleteContext: resourceDatabaseClusterWithShardsDelete,
		UpdateContext: resourceDatabaseClusterWithShardsUpdate,
		CustomizeDiff: resourceDatabaseCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				config := meta.(configer)
				DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
				if err != nil {
					return nil, fmt.Errorf("error creating VKCS database client: %s", err)
				}

				if resourceDatabaseClusterWithShardsRead(ctx, d, meta).HasError() {
					return nil, fmt.Errorf("error reading vkcs_db_cluster_with_shards")
				}

				cluster, err := dbClusterGet(DatabaseV1Client, d.Id()).extract()
				if err != nil {
					return nil, fmt.Errorf("error retrieving vkcs_db_cluster_with_shards")
				}

				shardIDs := make(map[string]int)
				shards := make([]map[string]interface{}, 0)
				for _, inst := range cluster.Instances {
					if _, ok := shardIDs[inst.ShardID]; ok {
						shardIDs[inst.ShardID]++
						continue
					}
					shardIDs[inst.ShardID] = 1
					newShard := flattenDatabaseClusterShard(inst)
					if inst.WalVolume != nil {
						newShard["wal_volume"] = flattenDatabaseClusterWalVolume(*inst.WalVolume)
					}
					shards = append(shards, newShard)
				}
				for _, shard := range shards {
					shard["size"] = shardIDs[shard["shard_id"].(string)]
				}
				d.Set("shard", shards)

				capabilities, err := clusterGetCapabilities(DatabaseV1Client, d.Id()).extract()
				if err != nil {
					return nil, fmt.Errorf("error getting cluster capabilities")
				}
				d.Set("capabilities", flattenDatabaseInstanceCapabilities(capabilities))
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
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								if v != Clickhouse {
									errs = append(errs, fmt.Errorf("datastore type must be %v, got: %s", getClusterWithShardsDatastores(), v))
								}
								return
							},
							Description: "Type of the datastore. Changing this creates a new cluster. Type of the datastore must be \"clickhouse\".",
						},
					},
				},
				Description: "Object that represents datastore of the cluster. Changing this creates a new cluster.",
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

			"floating_ip_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Boolean field that indicates whether floating ip is created for cluster. Changing this creates a new cluster.",
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
					},
				},
				Description: "Object that represents backup to restore instance from. **New since v.0.1.4**.",
			},

			"cloud_monitoring_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Description: "Enable cloud monitoring for the cluster. Changing this for Redis or MongoDB creates a new instance.",
			},

			"shard": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shard_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The ID of the shard. Changing this creates a new cluster.",
						},

						"size": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    false,
							Description: "The number of instances in the cluster shard.",
						},

						"flavor_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    false,
							Computed:    false,
							Description: "The ID of flavor for the cluster shard.",
						},
						"volume_size": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    false,
							Computed:    false,
							Description: "Size of the cluster shard instance volume.",
						},

						"volume_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    false,
							Computed:    false,
							Description: "The type of the cluster shard instance volume.",
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
										Description: "The type of the cluster wal volume.",
									},
								},
							},
							Description: "Object that represents wal volume of the cluster.",
						},

						"network": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uuid": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Description: "The id of the network. Changing this creates a new cluster." +
											"**Note** Although this argument is marked as optional, it is actually required at the moment. Not setting a value for it may cause an error.",
									},
									"port": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "The port id of the network. Changing this creates a new cluster. ***Deprecated*** This argument is deprecated, please do not use it.",
										Deprecated:  "This argument is deprecated, please do not use it.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "The id of the subnet. Changing this creates a new cluster. **New since v.0.1.15**.",
									},
								},
								Description: "Object that represents network of the cluster shard. Changing this creates a new cluster.",
							},
						},

						"availability_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    false,
							ForceNew:    true,
							Description: "The name of the availability zone of the cluster shard. Changing this creates a new cluster.",
						},

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
								},
							},
							Description: "Shard instances info. **New since v.0.1.15**.",
						},
					},
				},
				Description: "Object that represents cluster shard. There can be several instances of this object.",
			},
		},
		Description: "Provides a db cluster with shards resource. This can be used to create, modify and delete db cluster with shards for clickhouse datastore.",
	}
}

func resourceDatabaseClusterWithShardsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	createOpts := &dbClusterCreateOpts{
		Name:                   d.Get("name").(string),
		FloatingIPEnabled:      d.Get("floating_ip_enabled").(bool),
		CloudMonitoringEnabled: d.Get("cloud_monitoring_enabled").(bool),
	}

	message := "unable to determine vkcs_db_cluster_with_shards"

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

	var instanceCount int
	shardsRaw := d.Get("shard").([]interface{})
	shardInfo := make([]dbClusterInstanceCreateOpts, len(shardsRaw))
	shardsSize := make([]int, len(shardInfo))

	for i, shardRaw := range shardsRaw {
		shardMap := shardRaw.(map[string]interface{})
		shardSize := shardMap["size"].(int)
		shardsSize[i] = shardSize
		instanceCount += shardSize
		volumeSize := shardMap["volume_size"].(int)
		shardInfo[i].Volume = &volume{Size: &volumeSize, VolumeType: shardMap["volume_type"].(string)}
		shardInfo[i].Nics, _ = extractDatabaseNetworks(shardMap["network"].([]interface{}))
		shardInfo[i].AvailabilityZone = shardMap["availability_zone"].(string)
		shardInfo[i].FlavorRef = shardMap["flavor_id"].(string)
		shardInfo[i].ShardID = shardMap["shard_id"].(string)
		walVolumeV := shardMap["wal_volume"].([]interface{})
		if len(walVolumeV) > 0 {
			walVolumeOpts, err := extractDatabaseWalVolume(walVolumeV)
			if err != nil {
				return diag.Errorf("%s wal_volume", message)
			}
			shardInfo[i].Walvolume = &walVolume{Size: &walVolumeOpts.Size, VolumeType: walVolumeOpts.VolumeType}
		}
	}

	for i := 0; i < len(shardInfo); i++ {
		shardInfo[i].Keypair = d.Get("keypair").(string)
	}
	instances := make([]dbClusterInstanceCreateOpts, instanceCount)
	k := 0
	for i, shardSize := range shardsSize {
		for j := 0; j < shardSize; j++ {
			instances[k] = shardInfo[i]
			k++
		}
	}
	createOpts.Instances = instances

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

	log.Printf("[DEBUG] vkcs_db_cluster_with_shards create options: %#v", createOpts)
	clust := dbCluster{}
	clust.Cluster = createOpts

	cluster, err := dbClusterCreate(DatabaseV1Client, clust).extract()
	if err != nil {
		return diag.Errorf("error creating vkcs_db_cluster_with_shards: %s", err)
	}

	// Wait for the cluster to become available.
	log.Printf("[DEBUG] Waiting for vkcs_db_cluster_with_shards %s to become available", cluster.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(dbClusterStatusBuild)},
		Target:     []string{string(dbClusterStatusActive)},
		Refresh:    databaseClusterStateRefreshFunc(DatabaseV1Client, cluster.ID, checkCapabilities),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      dbInstanceDelay,
		MinTimeout: dbInstanceMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for vkcs_db_cluster_with_shards %s to become ready: %s", cluster.ID, err)
	}

	if configuration, ok := d.GetOk("configuration_id"); ok {
		log.Printf("[DEBUG] Attaching configuration %s to vkcs_db_cluster_with_shards %s", configuration, cluster.ID)
		var attachConfigurationOpts dbClusterAttachConfigurationGroupOpts
		attachConfigurationOpts.ConfigurationAttach.ConfigurationID = configuration.(string)
		err := instanceAttachConfigurationGroup(DatabaseV1Client, cluster.ID, &attachConfigurationOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error attaching configuration group %s to vkcs_db_cluster_with_shards %s: %s",
				configuration, cluster.ID, err)
		}
	}

	// Store the ID now
	d.SetId(cluster.ID)
	return resourceDatabaseClusterWithShardsRead(ctx, d, meta)
}

func resourceDatabaseClusterWithShardsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	cluster, err := dbClusterGet(DatabaseV1Client, d.Id()).extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "error retrieving vkcs_db_cluster_with_shards"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_db_cluster_with_shards %s: %#v", d.Id(), cluster)

	d.Set("name", cluster.Name)
	d.Set("datastore", flattenDatabaseInstanceDatastore(*cluster.DataStore))

	d.Set("configuration_id", cluster.ConfigurationID)
	if _, ok := d.GetOk("disk_autoexpand"); ok {
		d.Set("disk_autoexpand", flattenDatabaseInstanceAutoExpand(cluster.AutoExpand, cluster.MaxDiskSize))
	}
	if _, ok := d.GetOk("wal_disk_autoexpand"); ok {
		d.Set("wal_disk_autoexpand", flattenDatabaseInstanceAutoExpand(cluster.WalAutoExpand, cluster.WalMaxDiskSize))
	}

	shardsRaw := d.Get("shard").([]interface{})
	shards := make([]map[string]interface{}, len(shardsRaw))
	shardsInstances := getDatabaseClusterShardInstances(cluster.Instances)
	for i, shardRaw := range shardsRaw {
		shard := shardRaw.(map[string]interface{})
		shardID := shard["shard_id"].(string)
		shard["instances"] = shardsInstances[shardID]
		shards[i] = shard
	}
	d.Set("shard", shards)

	return nil
}

func resourceDatabaseClusterWithShardsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		log.Printf("Detaching configuration %s from vkcs_db_cluster_with_shards %s", old, d.Id())

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_db_cluster_with_shards %s to become ready: %s", d.Id(), err)
		}

		if new != "" {
			var attachConfigurationOpts dbClusterAttachConfigurationGroupOpts
			attachConfigurationOpts.ConfigurationAttach.ConfigurationID = new.(string)
			err := dbClusterAction(DatabaseV1Client, d.Id(), &attachConfigurationOpts).ExtractErr()
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("Attaching configuration %s to vkcs_db_cluster_with_shards %s", new, d.Id())

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("error waiting for vkcs_db_cluster_with_shards %s to become ready: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("disk_autoexpand") {
		_, new := d.GetChange("disk_autoexpand")
		autoExpandProperties, err := extractDatabaseAutoExpand(new.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_cluster_with_shards disk_autoexpand")
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
			return diag.Errorf("error waiting for vkcs_db_cluster_with_shards %s to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("capabilities") {
		_, newCapabilities := d.GetChange("capabilities")
		newCapabilitiesOpts, err := extractDatabaseCapabilities(newCapabilities.([]interface{}))
		if err != nil {
			return diag.Errorf("unable to determine vkcs_db_cluster_with_shards capability")
		}
		var applyCapabilityOpts dbClusterApplyCapabilityOpts
		applyCapabilityOpts.ApplyCapability.Capabilities = newCapabilitiesOpts

		err = dbClusterAction(DatabaseV1Client, d.Id(), &applyCapabilityOpts).ExtractErr()

		if err != nil {
			return diag.Errorf("error applying capability to vkcs_db_cluster_with_shards %s: %s", d.Id(), err)
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
			return diag.Errorf("error applying capability to vkcs_db_cluster_with_shards %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("cloud_monitoring_enabled") {
		_, new := d.GetChange("cloud_monitoring_enabled")
		var cloudMonitoringOpts updateCloudMonitoringOpts
		cloudMonitoringOpts.CloudMonitoring.Enable = new.(bool)

		err := dbClusterAction(DatabaseV1Client, d.Id(), &cloudMonitoringOpts).ExtractErr()
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("Update cloud_monitoring_enabled in vkcs_db_cluster_with_shards %s", d.Id())
	}

	return resourceDatabaseClusterWithShardsRead(ctx, d, meta)
}

func resourceDatabaseClusterWithShardsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	err = dbClusterDelete(DatabaseV1Client, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_db_cluster_with_shards"))
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
		return diag.Errorf("error waiting for vkcs_db_cluster_with_shards %s to delete: %s", d.Id(), err)
	}

	return nil
}
