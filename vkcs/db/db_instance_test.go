//go:build db_acc_test
// +build db_acc_test

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractInstanceDatastore(t *testing.T) {
	datastore := []interface{}{
		map[string]interface{}{
			"version": "foo",
			"type":    "bar",
		},
	}

	expected := dataStoreShort{
		Version: "foo",
		Type:    "bar",
	}

	actual, _ := extractDatabaseDatastore(datastore)
	assert.Equal(t, expected, actual)
}

func TestExtractDatabaseNetworks(t *testing.T) {
	network := []interface{}{
		map[string]interface{}{
			"uuid":        "foobar",
			"port":        "",
			"fixed_ip_v4": "",
		},
	}

	expected := []networkOpts{
		{
			UUID: "foobar",
		},
	}

	actual, _ := extractDatabaseNetworks(network)
	assert.Equal(t, expected, actual)
}

func TestExtractDatabaseAutoExpand(t *testing.T) {
	autoExpand := []interface{}{
		map[string]interface{}{
			"autoexpand":    true,
			"max_disk_size": 1000,
		},
	}

	expected := instanceAutoExpandOpts{
		AutoExpand:  true,
		MaxDiskSize: 1000,
	}

	actual, _ := extractDatabaseAutoExpand(autoExpand)
	assert.Equal(t, expected, actual)
}

func TestExtractDatabaseWalVolume(t *testing.T) {
	walVolume := []interface{}{
		map[string]interface{}{
			"size":          10,
			"volume_type":   "ms1",
			"autoexpand":    true,
			"max_disk_size": 1000,
		},
	}

	expected := walVolumeOpts{
		Size:        10,
		VolumeType:  "ms1",
		AutoExpand:  true,
		MaxDiskSize: 1000,
	}

	actual, _ := extractDatabaseWalVolume(walVolume)
	assert.Equal(t, expected, actual)
}

func TestExtractDatabaseCapabilities(t *testing.T) {
	capabilities := []interface{}{
		map[string]interface{}{
			"name": "node_exporter",
			"settings": map[string]string{
				"listen_port": "9100",
			},
		},
		map[string]interface{}{
			"name": "mysqld_exporter",
			"settings": map[string]string{
				"listen_port": "9104",
			},
		},
	}

	expected := []instanceCapabilityOpts{
		{
			Name: "node_exporter",
			Params: map[string]string{
				"listen_port": "9100",
			},
		},
		{
			Name: "mysqld_exporter",
			Params: map[string]string{
				"listen_port": "9104",
			},
		},
	}

	actual, _ := extractDatabaseCapabilities(capabilities)
	assert.Equal(t, expected, actual)
}

func TestExtractDatabaseRestorePoint(t *testing.T) {
	restorepoint := []interface{}{
		map[string]interface{}{
			"backup_id": "foo",
			"target":    "bar",
		},
	}

	expected := restorePoint{
		BackupRef: "foo",
		Target:    "bar",
	}

	actual, _ := extractDatabaseRestorePoint(restorepoint)
	assert.Equal(t, expected, actual)
}
