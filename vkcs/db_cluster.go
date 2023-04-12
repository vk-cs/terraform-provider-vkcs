package vkcs

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenDatabaseClusterWalVolume(w walVolume) []map[string]interface{} {
	walvolume := make([]map[string]interface{}, 1)
	walvolume[0] = make(map[string]interface{})
	walvolume[0]["size"] = w.Size
	walvolume[0]["volume_type"] = w.VolumeType
	return walvolume
}

func flattenDatabaseClusterInstances(insts []dbClusterInstanceResp) []map[string]interface{} {
	instances := make([]map[string]interface{}, len(insts))
	for i, inst := range insts {
		instances[i] = flattenDatabaseClusterInstance(inst)
	}

	return instances
}

func flattenDatabaseClusterInstance(inst dbClusterInstanceResp) map[string]interface{} {
	instance := make(map[string]interface{})
	instance["instance_id"] = inst.ID
	instance["ip"] = inst.IP
	instance["role"] = inst.Role

	return instance
}

func flattenDatabaseClusterShards(shardsInsts map[string][]dbClusterInstanceResp) (r []map[string]interface{}) {
	for id, insts := range shardsInsts {
		r = append(r, flattenDatabaseClusterShard(id, insts))
	}
	return
}

func flattenDatabaseClusterShard(id string, shardInsts []dbClusterInstanceResp) map[string]interface{} {
	shard := make(map[string]interface{})
	shard["shard_id"] = id
	shard["size"] = len(shardInsts)
	shard["flavor_id"] = shardInsts[0].Flavor.ID
	shard["volume_size"] = shardInsts[0].Volume.Size
	shard["volume_type"] = shardInsts[0].Volume.VolumeType
	if walVolume := shardInsts[0].WalVolume; walVolume != nil {
		shard["wal_volume"] = flattenDatabaseClusterWalVolume(*walVolume)
	}
	shard["instances"] = flattenDatabaseClusterShardInstances(shardInsts)
	return shard
}

func getDatabaseClusterShardInstances(insts []dbClusterInstanceResp) map[string][]dbClusterInstanceResp {
	shardsInstances := make(map[string][]dbClusterInstanceResp)
	for _, inst := range insts {
		shardsInstances[inst.ShardID] = append(shardsInstances[inst.ShardID], inst)
	}
	return shardsInstances
}

func flattenDatabaseClusterShardInstances(insts []dbClusterInstanceResp) (r []map[string]interface{}) {
	for _, inst := range insts {
		r = append(r, flattenDatabaseClusterShardInstance(inst))
	}
	return
}

func flattenDatabaseClusterShardInstance(inst dbClusterInstanceResp) map[string]interface{} {
	instance := make(map[string]interface{})
	instance["instance_id"] = inst.ID
	instance["ip"] = inst.IP
	return instance
}

func expandDatabaseClusterShrinkOptions(v []interface{}) []string {
	opts := make([]string, len(v))
	for i, opt := range v {
		opts[i] = opt.(string)
	}
	return opts
}

func databaseClusterDetermineShrinkedInstances(toDelete int, shrinkOptions []string, instances []dbClusterInstanceResp, shardID string) ([]dbClusterShrinkOpts, error) {
	ids := []dbClusterShrinkOpts{}
	foundIDs := 0
	if len(shrinkOptions) == 0 {
		for _, instance := range instances {
			if instance.Role != dbClusterInstanceRoleLeader && instance.ShardID == shardID {
				ids = append(ids, dbClusterShrinkOpts{ID: instance.ID})
				foundIDs++
				if foundIDs == toDelete {
					break
				}
			}
		}
		if foundIDs != toDelete {
			return nil, fmt.Errorf("not enough instances to shrink cluster")
		}
	} else {
		err := databaseClusterValidateShrinkOptions(shrinkOptions, instances, shardID)
		if err != nil {
			return nil, fmt.Errorf("invalid shrink options: %s", err)
		}
		for _, instance := range instances {
			needToRemain := false
			for _, opt := range shrinkOptions {
				if instance.ID == opt {
					needToRemain = true
				}
			}
			if !needToRemain {
				ids = append(ids, dbClusterShrinkOpts{ID: instance.ID})
				foundIDs++
			}
		}
		if foundIDs != toDelete {
			return nil, fmt.Errorf("invalid shrink options: not enough instances to delete")
		}
	}

	return ids, nil
}

func databaseClusterValidateShrinkOptions(shrinkOptions []string, instances []dbClusterInstanceResp, shardID string) error {
	for _, opt := range shrinkOptions {
		optIsValid := false
		for _, instance := range instances {
			if instance.ID == opt && instance.ShardID == shardID {
				optIsValid = true
			}
		}
		if !optIsValid {
			if shardID != "" {
				return fmt.Errorf("shard %s does not have instance: %s", shardID, opt)
			}
			return fmt.Errorf("cluster does not have instance: %s", opt)
		}
	}
	return nil
}

func databaseClusterExpandShards(d *schema.ResourceData) (r []map[string]interface{}) {
	shardsRaw := d.Get("shard").([]interface{})
	for _, shRaw := range shardsRaw {
		r = append(r, shRaw.(map[string]interface{}))
	}
	return
}

func shardIndex(d *schema.ResourceData, shardID string) (int, error) {
	shards := databaseClusterExpandShards(d)
	for i, sh := range shards {
		if sh["shard_id"] == shardID {
			return i, nil
		}
	}
	return 0, fmt.Errorf("shard with %s not found", shardID)
}

func shardPathPrefix(d *schema.ResourceData, shardID string) (string, error) {
	if shardID != "" {
		i, err := shardIndex(d, shardID)
		if err != nil {
			return "", fmt.Errorf("%w: %s", errDBClusterShardNotFound, err)
		}
		return fmt.Sprintf("shard.%d.", i), nil
	}
	return "", nil
}

func databaseClusterCheckDeleted(d *schema.ResourceData, err error) error {
	if _, ok := err.(gophercloud.ErrDefault404); ok {
		d.SetId("")
		return nil
	}
	return fmt.Errorf("%w: %s", errDBClusterNotFound, err)
}

type dbResourceUpdateContext struct {
	Ctx       context.Context
	Client    databaseClient
	D         *schema.ResourceData
	StateConf *resource.StateChangeConf
}

func (uCtx *dbResourceUpdateContext) WaitForStateContext() error {
	_, err := uCtx.StateConf.WaitForStateContext(uCtx.Ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", errDBClusterUpdateWait, err)
	}
	return nil
}

var (
	errDBClusterNotFound      = errors.New("cluster not found")
	errDBClusterShardNotFound = errors.New("unable to determine shard")
	errDBClusterUpdateWait    = errors.New("error waiting for cluster to become ready")

	errDBClusterUpdateDiskAutoexpand           = errors.New("error updating disk_autoexpand")
	errDBClusterUpdateDiskAutoexpandExtract    = errors.New("unable to determine disk_autoexpand")
	errDBClusterUpdateWalDiskAutoexpand        = errors.New("error updating wal_disk_autoexpand")
	errDBClusterUpdateWalDiskAutoexpandExtract = errors.New("unable to determine wal_disk_autoexpand")
	errDBClusterUpdateCloudMonitoring          = errors.New("error updating cloud_monitoring_enabled")

	errDBClusterActionUpdateConfiguration      = errors.New("error updating configuration for cluster")
	errDBClusterActionApplyCapabitilies        = errors.New("error applying capabilities")
	errDBClusterActionApplyCapabilitiesExtract = errors.New("error extracting capabilities")
	errDBClusterActionResizeWalVolumeExtract   = errors.New("unable to determine wal_volume")
	errDBClusterActionGrow                     = errors.New("error growing cluster")
	errDBClusterActionShrink                   = errors.New("error shrinking cluster")
	errDBClusterActionShrinkWrongOptions       = errors.New("invalid shrink options")
	errDBClusterActionShrinkInstancesExtract   = errors.New("error determining instances to shrink")
	errDBClusterActionResizeVolume             = errors.New("error resizing volume")
	errDBClusterActionResizeWalVolume          = errors.New("error resizing wal_volume")
	errDBClusterActionResizeFlavor             = errors.New("error resizing flavor")
)

func databaseClusterActionUpdateConfiguration(updateCtx *dbResourceUpdateContext) error {
	old, new := updateCtx.D.GetChange("configuration_id")

	var detachOpts dbClusterDetachConfigurationGroupOpts
	detachOpts.ConfigurationDetach.ConfigurationID = old.(string)

	var attachOpts *dbClusterAttachConfigurationGroupOpts
	if new != "" {
		attachOpts = &dbClusterAttachConfigurationGroupOpts{}
		attachOpts.ConfigurationAttach.ConfigurationID = new.(string)
	}

	return databaseClusterActionUpdateConfigurationBase(updateCtx, &detachOpts, attachOpts)
}

func databaseClusterActionUpdateConfigurationBase(updateCtx *dbResourceUpdateContext, detachOpts *dbClusterDetachConfigurationGroupOpts, attachOpts *dbClusterAttachConfigurationGroupOpts) error {
	dbClient, clusterID := updateCtx.Client, updateCtx.D.Id()

	err := dbClusterAction(dbClient, clusterID, detachOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("%w: %s", errDBClusterActionUpdateConfiguration, err)
	}

	log.Printf("[DEBUG] Detaching configuration %s from cluster %s", detachOpts.ConfigurationDetach.ConfigurationID, clusterID)
	err = updateCtx.WaitForStateContext()
	if err != nil {
		return err
	}

	if attachOpts != nil {
		err := dbClusterAction(dbClient, clusterID, attachOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("%w: %s", errDBClusterActionUpdateConfiguration, err)
		}

		log.Printf("[DEBUG] Attaching configuration %s to cluster %s", attachOpts.ConfigurationAttach.ConfigurationID, clusterID)
		return updateCtx.WaitForStateContext()
	}

	return nil
}

func databaseClusterUpdateDiskAutoexpand(updateCtx *dbResourceUpdateContext) error {
	diskAutoexp := updateCtx.D.Get("disk_autoexpand")
	autoExpandProperties, err := extractDatabaseAutoExpand(diskAutoexp.([]interface{}))
	if err != nil {
		return errDBClusterUpdateDiskAutoexpandExtract
	}

	var autoExpandOpts dbClusterUpdateAutoExpandOpts
	if autoExpandProperties.AutoExpand {
		autoExpandOpts.Cluster.VolumeAutoresizeEnabled = 1
	} else {
		autoExpandOpts.Cluster.VolumeAutoresizeEnabled = 0
	}
	autoExpandOpts.Cluster.VolumeAutoresizeMaxSize = autoExpandProperties.MaxDiskSize

	return databaseClusterUpdateDiskAutoexpandBase(updateCtx, autoExpandOpts)
}

func databaseClusterUpdateDiskAutoexpandBase(updateCtx *dbResourceUpdateContext, autoExpandOpts dbClusterUpdateAutoExpandOpts) error {
	dbClient, clusterID := updateCtx.Client, updateCtx.D.Id()

	err := dbClusterUpdateAutoExpand(dbClient, clusterID, &autoExpandOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("%w: %s", errDBClusterUpdateDiskAutoexpand, err)
	}

	updateCtx.StateConf.Pending = []string{string(dbClusterStatusUpdating)}
	updateCtx.StateConf.Target = []string{string(dbClusterStatusActive)}

	log.Printf("[DEBUG] Waiting for cluster %s to become ready after updating disk_autoexpand", clusterID)
	return updateCtx.WaitForStateContext()
}

func databaseClusterUpdateWalDiskAutoexpand(updateCtx *dbResourceUpdateContext) error {
	walDiskAutoexp := updateCtx.D.Get("wal_disk_autoexpand")
	walAutoExpandProperties, err := extractDatabaseAutoExpand(walDiskAutoexp.([]interface{}))
	if err != nil {
		return errDBClusterUpdateWalDiskAutoexpandExtract
	}

	var walAutoExpandOpts dbClusterUpdateAutoExpandWalOpts
	if walAutoExpandProperties.AutoExpand {
		walAutoExpandOpts.Cluster.WalVolume.VolumeAutoresizeEnabled = 1
	} else {
		walAutoExpandOpts.Cluster.WalVolume.VolumeAutoresizeEnabled = 0
	}
	walAutoExpandOpts.Cluster.WalVolume.VolumeAutoresizeMaxSize = walAutoExpandProperties.MaxDiskSize

	return databaseClusterUpdateWalDiskAutoexpandBase(updateCtx, walAutoExpandOpts)
}

func databaseClusterUpdateWalDiskAutoexpandBase(updateCtx *dbResourceUpdateContext, walAutoExpandOpts dbClusterUpdateAutoExpandWalOpts) error {
	dbClient, clusterID := updateCtx.Client, updateCtx.D.Id()

	err := dbClusterUpdateAutoExpand(dbClient, clusterID, &walAutoExpandOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("%w: %s", errDBClusterUpdateWalDiskAutoexpand, err)
	}

	updateCtx.StateConf.Pending = []string{string(dbClusterStatusUpdating)}
	updateCtx.StateConf.Target = []string{string(dbClusterStatusActive)}

	log.Printf("[DEBUG] Waiting for cluster %s to become ready after updating wal_disk_autoexpand", clusterID)
	return updateCtx.WaitForStateContext()
}

func databaseClusterUpdateCloudMonitoring(updateCtx *dbResourceUpdateContext) error {
	enabled := updateCtx.D.Get("cloud_monitoring_enabled").(bool)
	var cloudMonitoringOpts updateCloudMonitoringOpts
	cloudMonitoringOpts.CloudMonitoring.Enable = enabled
	return databaseClusterUpdateCloudMonitoringBase(updateCtx, cloudMonitoringOpts)
}

func databaseClusterUpdateCloudMonitoringBase(updateCtx *dbResourceUpdateContext, cloudMonitoringOpts updateCloudMonitoringOpts) error {
	clusterID := updateCtx.D.Id()
	err := dbClusterAction(updateCtx.Client, clusterID, &cloudMonitoringOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("%w: %s", errDBClusterUpdateCloudMonitoring, err)
	}
	log.Printf("[DEBUG] Updated cloud_monitoring_enabled in cluster %s", clusterID)
	return nil
}

func databaseClusterActionApplyCapabilities(updateCtx *dbResourceUpdateContext) error {
	dbClient, clusterID := updateCtx.Client, updateCtx.D.Id()

	caps := updateCtx.D.Get("capabilities")
	opts, err := extractDatabaseCapabilities(caps.([]interface{}))
	if err != nil {
		return errDBClusterActionApplyCapabilitiesExtract
	}

	var applyCapabilityOpts dbClusterApplyCapabilityOpts
	applyCapabilityOpts.ApplyCapability.Capabilities = opts

	updateCtx.StateConf.Refresh = databaseClusterStateRefreshFunc(dbClient, clusterID, &opts)

	return databaseClusterActionApplyCapabilitiesBase(updateCtx, applyCapabilityOpts)
}

func databaseClusterActionApplyCapabilitiesBase(updateCtx *dbResourceUpdateContext, applyCapabilityOpts dbClusterApplyCapabilityOpts) error {
	dbClient, clusterID := updateCtx.Client, updateCtx.D.Id()

	err := dbClusterAction(dbClient, clusterID, &applyCapabilityOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("%w: %s", errDBClusterActionApplyCapabitilies, err)
	}

	updateCtx.StateConf.Pending = []string{string(dbClusterStatusCapabilityApplying), string(dbClusterStatusBuild)}
	updateCtx.StateConf.Target = []string{string(dbClusterStatusActive)}

	log.Printf("[DEBUG] Waiting for cluster %s to become ready after applying capability", clusterID)
	return updateCtx.WaitForStateContext()
}

func databaseClusterActionGrow(updateCtx *dbResourceUpdateContext, shardID string) error {
	d := updateCtx.D
	pathPrefix, err := shardPathPrefix(d, shardID)
	if err != nil {
		return err
	}

	volumeSize := d.Get(pathPrefix + "volume_size").(int)
	growOpts := dbClusterGrowOpts{
		Keypair:          d.Get("keypair").(string),
		AvailabilityZone: d.Get(pathPrefix + "availability_zone").(string),
		FlavorRef:        d.Get(pathPrefix + "flavor_id").(string),
		Volume:           &volume{Size: &volumeSize, VolumeType: d.Get(pathPrefix + "volume_type").(string)},
		ShardID:          shardID,
	}

	if v, ok := d.GetOk(pathPrefix + "wal_volume"); ok {
		walVolumeOpts, err := extractDatabaseWalVolume(v.([]interface{}))
		if err != nil {
			return errDBClusterActionResizeWalVolumeExtract
		}
		growOpts.Walvolume = &walVolume{
			Size:       &walVolumeOpts.Size,
			VolumeType: walVolumeOpts.VolumeType,
		}
	}

	var old, new interface{}
	if shardID != "" {
		old, new = d.GetChange(pathPrefix + "size")
	} else {
		old, new = d.GetChange("cluster_size")
	}
	growSize := new.(int) - old.(int)

	if shardID != "" {
		updateCtx.StateConf.Pending = []string{string(dbClusterStatusGrow), string(dbClusterStatusBuild)}
	} else {
		updateCtx.StateConf.Pending = []string{string(dbClusterStatusGrow)}
	}
	updateCtx.StateConf.Target = []string{string(dbClusterStatusActive)}

	return databaseClusterActionGrowBase(updateCtx, growOpts, growSize)
}

func databaseClusterActionGrowBase(updateCtx *dbResourceUpdateContext, growOpts dbClusterGrowOpts, growSize int) error {
	clusterID := updateCtx.D.Id()
	opts := make([]dbClusterGrowOpts, growSize)
	for i := 0; i < growSize; i++ {
		opts[i] = growOpts
	}
	growClusterOpts := dbClusterGrowClusterOpts{Grow: opts}

	err := dbClusterAction(updateCtx.Client, clusterID, &growClusterOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("%w: %s", errDBClusterActionGrow, err)
	}

	log.Printf("[DEBUG] Growing cluster %s", clusterID)
	return updateCtx.WaitForStateContext()
}

func databaseClusterActionShrink(updateCtx *dbResourceUpdateContext, shardID string) error {
	d := updateCtx.D
	pathPrefix, err := shardPathPrefix(d, shardID)
	if err != nil {
		return err
	}

	var old, new interface{}
	if shardID != "" {
		old, new = d.GetChange(pathPrefix + "size")
	} else {
		old, new = d.GetChange("cluster_size")
	}
	newSize, shrinkSize := new.(int), old.(int)-new.(int)

	rawShrinkOptions := d.Get(pathPrefix + "shrink_options").([]interface{})
	shrinkOptions := expandDatabaseClusterShrinkOptions(rawShrinkOptions)
	if len(shrinkOptions) > 0 && len(shrinkOptions) != newSize {
		return fmt.Errorf("%w: number of instances in shrink options should equal new size",
			errDBClusterActionShrinkWrongOptions)
	}

	cluster, err := dbClusterGet(updateCtx.Client, d.Id()).extract()
	if err != nil {
		return databaseClusterCheckDeleted(d, err)
	}

	ids, err := databaseClusterDetermineShrinkedInstances(shrinkSize, shrinkOptions, cluster.Instances, shardID)
	if err != nil {
		return fmt.Errorf("%w: %s", errDBClusterActionShrinkInstancesExtract, err)
	}

	if shardID != "" {
		updateCtx.StateConf.Pending = []string{string(dbClusterStatusShrink), string(dbClusterStatusBuild)}
	} else {
		updateCtx.StateConf.Pending = []string{string(dbClusterStatusShrink)}
	}
	updateCtx.StateConf.Target = []string{string(dbClusterStatusActive)}

	return databaseClusterActionShrinkBase(updateCtx, ids)
}

func databaseClusterActionShrinkBase(updateCtx *dbResourceUpdateContext, shrinkOpts []dbClusterShrinkOpts) error {
	clusterID := updateCtx.D.Id()
	shrinkClusterOpts := dbClusterShrinkClusterOpts{
		Shrink: shrinkOpts,
	}

	err := dbClusterAction(updateCtx.Client, clusterID, &shrinkClusterOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("%w: %s", errDBClusterActionShrink, err)
	}

	log.Printf("[DEBUG] Shrinking cluster %s", clusterID)
	return updateCtx.WaitForStateContext()
}

func databaseClusterActionResizeVolume(updateCtx *dbResourceUpdateContext, shardID string) error {
	d := updateCtx.D
	pathPrefix, err := shardPathPrefix(d, shardID)
	if err != nil {
		return err
	}

	_, volumeSize := d.GetChange(pathPrefix + "volume_size")
	var resizeVolumeOpts dbClusterResizeVolumeOpts
	resizeVolumeOpts.Resize.Volume.Size = volumeSize.(int)
	resizeVolumeOpts.Resize.ShardID = shardID

	updateCtx.StateConf.Pending = []string{string(dbClusterStatusResize)}
	updateCtx.StateConf.Target = []string{string(dbClusterStatusActive)}

	return databaseClusterActionResizeVolumeBase(updateCtx, resizeVolumeOpts)
}

func databaseClusterActionResizeVolumeBase(updateCtx *dbResourceUpdateContext, opts dbClusterResizeVolumeOpts) error {
	clusterID := updateCtx.D.Id()
	err := dbClusterAction(updateCtx.Client, clusterID, &opts).ExtractErr()
	if err != nil {
		return fmt.Errorf("%w: %s", errDBClusterActionResizeVolume, err)
	}
	log.Printf("[DEBUG] Resizing volume from cluster %s", clusterID)
	return updateCtx.WaitForStateContext()
}

func databaseClusterActionResizeWalVolume(updateCtx *dbResourceUpdateContext, shardID string) error {
	d := updateCtx.D
	pathPrefix, err := shardPathPrefix(d, shardID)
	if err != nil {
		return err
	}

	old, new := d.GetChange(pathPrefix + "wal_volume")
	walVolumeOptsNew, err := extractDatabaseWalVolume(new.([]interface{}))
	if err != nil {
		return errDBClusterActionResizeWalVolumeExtract
	}

	walVolumeOptsOld, err := extractDatabaseWalVolume(old.([]interface{}))
	if err != nil {
		return errDBClusterActionResizeWalVolumeExtract
	}

	if walVolumeOptsNew.Size != walVolumeOptsOld.Size {
		var resizeWalVolumeOpts dbClusterResizeWalVolumeOpts
		resizeWalVolumeOpts.Resize.Volume.Size = walVolumeOptsNew.Size
		resizeWalVolumeOpts.Resize.Volume.Kind = "wal"
		resizeWalVolumeOpts.Resize.ShardID = shardID

		updateCtx.StateConf.Pending = []string{string(dbClusterStatusResize)}
		updateCtx.StateConf.Target = []string{string(dbClusterStatusActive)}

		return databaseClusterActionResizeWalVolumeBase(updateCtx, resizeWalVolumeOpts)
	}

	return nil
}

func databaseClusterActionResizeWalVolumeBase(updateCtx *dbResourceUpdateContext, opts dbClusterResizeWalVolumeOpts) error {
	clusterID := updateCtx.D.Id()
	err := dbClusterAction(updateCtx.Client, clusterID, &opts).ExtractErr()
	if err != nil {
		return fmt.Errorf("%w: %s", errDBClusterActionResizeWalVolume, err)
	}
	log.Printf("[DEBUG] Resizing wal_folume from cluster %s", clusterID)
	return updateCtx.WaitForStateContext()
}

func databaseClusterActionResizeFlavor(updateCtx *dbResourceUpdateContext, shardID string) error {
	d := updateCtx.D
	pathPrefix, err := shardPathPrefix(d, shardID)
	if err != nil {
		return err
	}

	var resizeOpts dbClusterResizeOpts
	resizeOpts.Resize.FlavorRef = d.Get(pathPrefix + "flavor_id").(string)
	resizeOpts.Resize.ShardID = shardID

	updateCtx.StateConf.Pending = []string{string(dbClusterStatusResize)}
	updateCtx.StateConf.Target = []string{string(dbClusterStatusActive)}

	return databaseClusterActionResizeFlavorBase(updateCtx, resizeOpts)
}

func databaseClusterActionResizeFlavorBase(updateCtx *dbResourceUpdateContext, opts dbClusterResizeOpts) error {
	clusterID := updateCtx.D.Id()
	err := dbClusterAction(updateCtx.Client, clusterID, &opts).ExtractErr()
	if err != nil {
		return fmt.Errorf("%w: %s", errDBClusterActionResizeFlavor, err)
	}
	log.Printf("[DEBUG] Resizing flavor from cluster %s", clusterID)
	return updateCtx.WaitForStateContext()
}

func getClusterStatus(c *dbClusterResp) string {
	instancesStatus := string(dbInstanceStatusActive)
	for _, inst := range c.Instances {
		if inst.Status == string(dbInstanceStatusError) {
			return inst.Status
		}
		if inst.Status == string(dbInstanceStatusBuild) || inst.Status == string(dbInstanceStatusResize) {
			instancesStatus = inst.Status
		}
	}
	if c.Task.Name == "NONE" {
		switch instancesStatus {
		case string(dbInstanceStatusActive):
			return string(dbClusterStatusActive)
		case string(dbInstanceStatusBuild):
			return string(dbClusterStatusBuild)
		case string(dbInstanceStatusResize):
			return string(dbClusterStatusResize)
		}
	}

	return c.Task.Name
}

func databaseClusterStateRefreshFunc(client databaseClient, clusterID string, capabilitiesOpts *[]instanceCapabilityOpts) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, err := dbClusterGet(client, clusterID).extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return c, "DELETED", nil
			}
			return nil, "", err
		}

		clusterStatus := getClusterStatus(c)
		if clusterStatus == "error" {
			return c, clusterStatus, fmt.Errorf("there was an error creating the database cluster")
		}
		if clusterStatus == string(dbClusterStatusActive) {
			if capabilitiesOpts != nil {
				for _, i := range c.Instances {
					instCapabilities, err := instanceGetCapabilities(client, i.ID).extract()
					if err != nil {
						return nil, "", fmt.Errorf("error getting cluster instance capabilities: %s", err)
					}
					capabilitiesReady, err := checkDBMSCapabilities(*capabilitiesOpts, instCapabilities)
					if err != nil {
						return nil, "", err
					}
					if !capabilitiesReady {
						return c, string(dbClusterStatusBuild), nil
					}
				}
				return c, string(dbClusterStatusActive), nil
			}
		}

		return c, clusterStatus, nil
	}
}
