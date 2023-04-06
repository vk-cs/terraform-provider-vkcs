package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatabaseDatastoreCapabilitiesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseDatastoreCapabilitiesDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_db_datastore_capabilities.capabilities", "datastore_name", "mysql"),
				),
			},
		},
	})
}

const testAccDatabaseDatastoreCapabilitiesDataSourceBasic = `
data "vkcs_db_datastore" "datastore" {
	name = "mysql"
}

data "vkcs_db_datastore_capabilities" "capabilities" {
	datastore_name = data.vkcs_db_datastore.datastore.name
	datastore_version_id = data.vkcs_db_datastore.datastore.versions.0.id
}
`
