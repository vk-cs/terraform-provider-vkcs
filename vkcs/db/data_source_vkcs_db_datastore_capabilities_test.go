package db_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDatabaseDatastoreCapabilitiesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
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
