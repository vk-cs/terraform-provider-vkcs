package vkcs

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func flattenDatabaseClusterWalVolume(w walVolume) []map[string]interface{} {
	walvolume := make([]map[string]interface{}, 1)
	walvolume[0] = make(map[string]interface{})
	walvolume[0]["size"] = w.Size
	walvolume[0]["volume_type"] = w.VolumeType
	return walvolume
}

func flattenDatabaseClusterShard(inst dbClusterInstanceResp) map[string]interface{} {
	newShard := make(map[string]interface{})
	newShard["shard_id"] = inst.ShardID
	newShard["flavor_id"] = inst.Flavor.ID
	newShard["volume_size"] = inst.Volume.Size
	newShard["volume_type"] = dbImportedStatus
	return newShard
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

func getDatabaseClusterShardInstances(insts []dbClusterInstanceResp) map[string][]map[string]interface{} {
	shardsInstances := make(map[string][]map[string]interface{})
	for _, inst := range insts {
		flattenInst := flattenDatabaseClusterShardInstance(inst)
		shardsInstances[inst.ShardID] = append(shardsInstances[inst.ShardID], flattenInst)
	}

	return shardsInstances
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

func databaseClusterDetermineShrinkedInstances(toDelete int, shrinkOptions []string, instances []dbClusterInstanceResp) ([]dbClusterShrinkOpts, error) {
	ids := []dbClusterShrinkOpts{}
	foundIDs := 0
	if len(shrinkOptions) == 0 {
		for _, instance := range instances {
			if instance.Role != dbClusterInstanceRoleLeader {
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
		err := databaseClusterValidateShrinkOptions(shrinkOptions, instances)
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

func databaseClusterValidateShrinkOptions(shrinkOptions []string, instances []dbClusterInstanceResp) error {
	for _, opt := range shrinkOptions {
		optIsValid := false
		for _, instance := range instances {
			if instance.ID == opt {
				optIsValid = true
			}
		}
		if !optIsValid {
			return fmt.Errorf("cluster does not have instance: %s", opt)
		}
	}
	return nil
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
