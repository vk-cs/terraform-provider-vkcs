//go:build db_acc_test
// +build db_acc_test

package vkcs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/db/v1/databases"
)

func TestExtractDatabaseUserDatabases(t *testing.T) {
	dbs := []interface{}{"db1", "db2"}

	expected := databases.BatchCreateOpts{
		databases.CreateOpts{
			Name: "db1",
		},
		databases.CreateOpts{
			Name: "db2",
		},
	}

	actual, _ := extractDatabaseUserDatabases(dbs)
	assert.Equal(t, expected, actual)
}

func TestFlattenDatabaseUserDatabases(t *testing.T) {
	dbs := []databases.Database{
		{
			Name: "db1",
		},
		{
			Name: "db2",
		},
	}

	expected := []interface{}{"db1", "db2"}
	actual := flattenDatabaseUserDatabases(dbs)
	assert.Equal(t, expected, actual)
}
