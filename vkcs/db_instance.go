package vkcs

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/mitchellh/mapstructure"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Datastore names
const (
	Redis       = "redis"
	MongoDB     = "mongodb"
	PostgresPro = "postgrespro"
	Galera      = "galera_mysql"
	Postgres    = "postgresql"
	Clickhouse  = "clickhouse"
	MySQL       = "mysql"
	Tarantool   = "tarantool"
)

func getClusterDatastores() []string {
	return []string{Galera, Postgres, Tarantool}
}

func getClusterWithShardsDatastores() []string {
	return []string{Clickhouse}
}

func extractDatabaseRestorePoint(v []interface{}) (restorePoint, error) {
	var R restorePoint
	in := v[0].(map[string]interface{})
	err := mapstructure.Decode(in, &R)
	if err != nil {
		return R, err
	}
	return R, nil
}

func extractDatabaseDatastore(v []interface{}) (dataStoreShort, error) {
	var D dataStoreShort
	in := v[0].(map[string]interface{})
	err := mapStructureDecoder(&D, &in, decoderConfig)
	if err != nil {
		return D, err
	}
	return D, nil
}

func flattenDatabaseInstanceDatastore(d dataStoreShort) []map[string]interface{} {
	datastore := make([]map[string]interface{}, 1)
	datastore[0] = make(map[string]interface{})
	datastore[0]["type"] = d.Type
	datastore[0]["version"] = d.Version
	return datastore
}

func extractDatabaseNetworks(v []interface{}) ([]networkOpts, []string, error) {
	Networks := make([]networkOpts, len(v))
	var SecurityGroups []string
	for i, network := range v {
		var N networkOpts
		networkMap := network.(map[string]interface{})
		sg, ok := networkMap["security_groups"]
		if ok {
			SecurityGroups = expandToStringSlice(sg.(*schema.Set).List())
		}
		err := mapstructure.Decode(networkMap, &N)
		if err != nil {
			return nil, nil, err
		}
		Networks[i] = N
	}
	return Networks, SecurityGroups, nil
}

func extractDatabaseAutoExpand(v []interface{}) (instanceAutoExpandOpts, error) {
	var A instanceAutoExpandOpts
	in := v[0].(map[string]interface{})
	err := mapstructure.Decode(in, &A)
	if err != nil {
		return A, err
	}
	return A, nil
}

func flattenDatabaseInstanceAutoExpand(autoExpandInt int, maxDiskSize int) []map[string]interface{} {
	autoExpand := make([]map[string]interface{}, 1)
	autoExpand[0] = make(map[string]interface{})
	if autoExpandInt != 0 {
		autoExpand[0]["autoexpand"] = true
	} else {
		autoExpand[0]["autoexpand"] = false
	}
	autoExpand[0]["max_disk_size"] = maxDiskSize
	return autoExpand
}

func extractDatabaseWalVolume(v []interface{}) (walVolumeOpts, error) {
	var W walVolumeOpts
	in := v[0].(map[string]interface{})
	err := mapstructure.Decode(in, &W)
	if err != nil {
		return W, err
	}
	return W, nil
}

func flattenDatabaseInstanceWalVolume(w walVolume) []map[string]interface{} {
	walvolume := make([]map[string]interface{}, 1)
	walvolume[0] = make(map[string]interface{})
	walvolume[0]["size"] = w.Size
	walvolume[0]["volume_type"] = w.VolumeType
	return walvolume
}

func extractDatabaseCapabilities(v []interface{}) ([]instanceCapabilityOpts, error) {
	capabilities := make([]instanceCapabilityOpts, len(v))
	for i, capability := range v {
		var C instanceCapabilityOpts
		err := mapstructure.Decode(capability.(map[string]interface{}), &C)
		if err != nil {
			return nil, err
		}
		capabilities[i] = C
	}
	return capabilities, nil
}

func flattenDatabaseInstanceCapabilities(c []databaseCapability) []map[string]interface{} {
	capabilities := make([]map[string]interface{}, len(c))
	for i, capability := range c {
		capabilities[i] = make(map[string]interface{})
		capabilities[i]["name"] = capability.Name
		capabilities[i]["settings"] = capability.Params
	}
	return capabilities
}

func extractDatabaseBackupSchedule(v []interface{}) (backupSchedule, error) {
	var B backupSchedule
	in := v[0].(map[string]interface{})
	err := mapStructureDecoder(&B, &in, decoderConfig)
	if err != nil {
		return B, err
	}
	return B, nil
}

func flattenDatabaseBackupSchedule(b backupSchedule) []map[string]interface{} {
	backupschedule := make([]map[string]interface{}, 1)
	backupschedule[0] = make(map[string]interface{})
	backupschedule[0]["name"] = b.Name
	backupschedule[0]["start_hours"] = b.StartHours
	backupschedule[0]["start_minutes"] = b.StartMinutes
	backupschedule[0]["interval_hours"] = b.IntervalHours
	backupschedule[0]["keep_count"] = b.KeepCount

	return backupschedule
}

func databaseInstanceStateRefreshFunc(client databaseClient, instanceID string, capabilitiesOpts *[]instanceCapabilityOpts) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		i, err := instanceGet(client, instanceID).extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return i, "DELETED", nil
			}
			return nil, "", err
		}

		if i.Status == string(dbInstanceStatusError) {
			return i, i.Status, fmt.Errorf("there was an error creating the database instance")
		}

		if i.Status == string(dbInstanceStatusActive) {
			if capabilitiesOpts != nil {
				instCapabilities, err := instanceGetCapabilities(client, instanceID).extract()
				if err != nil {
					return nil, "", fmt.Errorf("error getting instance capabilities: %s", err)
				}

				capabilitiesReady, err := checkDBMSCapabilities(*capabilitiesOpts, instCapabilities)
				if err != nil {
					return nil, "", err
				}
				if capabilitiesReady {
					return i, string(dbInstanceStatusActive), nil
				} else {
					return i, string(dbInstanceStatusBuild), nil
				}
			}
		}
		return i, i.Status, nil
	}
}

func checkDBMSCapabilities(neededCapabilities []instanceCapabilityOpts, actualCapabilities []databaseCapability) (bool, error) {
	// this is workaround for situation when capabilities are applied sequentially and not all of them are returned by api
	if len(neededCapabilities) != len(actualCapabilities) {
		return false, nil
	}
	for _, neededCap := range neededCapabilities {
		found := false
		for _, actualCap := range actualCapabilities {
			if neededCap.Name == actualCap.Name {
				found = true
				if actualCap.Status == string(dbCapabilityStatusError) {
					return false, fmt.Errorf("error applying capabilities")
				}
				if actualCap.Status != string(dbCapabilityStatusActive) {
					return false, nil
				}
			}
		}
		if !found {
			return false, fmt.Errorf("error applying capabilities")
		}
	}
	return true, nil
}

func getDBMSResource(client databaseClient, dbmsID string) (interface{}, error) {
	instanceResource, err := instanceGet(client, dbmsID).extract()
	if err == nil {
		return instanceResource, nil
	}
	clusterResource, err := dbClusterGet(client, dbmsID).extract()
	if err == nil {
		return clusterResource, nil
	}
	return nil, err
}
