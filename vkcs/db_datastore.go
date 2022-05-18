package vkcs

func flattenDatabaseDatastoreVersions(vs *[]datastoreVersion) []map[string]interface{} {
	versions := make([]map[string]interface{}, len(*vs))
	for i, v := range *vs {
		versions[i] = make(map[string]interface{})
		versions[i]["name"] = v.Name
		versions[i]["id"] = v.ID
	}

	return versions
}

func flattenDatabaseDatastore(d *datastoreResp) map[string]interface{} {
	datastore := make(map[string]interface{})
	datastore["name"] = d.Name
	datastore["id"] = d.ID
	datastore["minimum_cpu"] = d.MinimumCPU
	datastore["minimum_ram"] = d.MinimumRAM
	datastore["versions"] = flattenDatabaseDatastoreVersions(d.Versions)
	datastore["volume_types"] = d.VolumeTypes
	datastore["cluster_volume_types"] = d.ClusterVolumeTypes
	return datastore
}
