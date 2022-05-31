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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"datastore": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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
						},
					},
				},
			},

			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: false,
				ForceNew: false,
			},

			"root_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},

			"root_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
				ForceNew:  false,
			},

			"floating_ip_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: false,
				ForceNew: true,
			},

			"keypair": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: false,
				ForceNew: true,
			},

			"disk_autoexpand": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"autoexpand": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: false,
						},
						"max_disk_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: false,
						},
					},
				},
			},

			"wal_disk_autoexpand": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"autoexpand": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: false,
						},
						"max_disk_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: false,
						},
					},
				},
			},

			"capabilities": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"settings": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
					},
				},
			},

			"restore_point": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"shard": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shard_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: false,
						},

						"flavor_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
							Computed: false,
						},
						"volume_size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: false,
							Computed: false,
						},

						"volume_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
							Computed: false,
						},

						"wal_volume": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: false,
									},
									"volume_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: false,
									},
								},
							},
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
									},
									"port": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},

						"availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: false,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceDatabaseClusterWithShardsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	createOpts := &dbClusterCreateOpts{
		Name:              d.Get("name").(string),
		FloatingIPEnabled: d.Get("floating_ip_enabled").(bool),
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
