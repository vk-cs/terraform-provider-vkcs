package db_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDatabaseDatastoreCapabilitiesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
				Config:                   testAccDatabaseDatastoreCapabilitiesDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_db_datastore_capabilities.capabilities", "datastore_name", "mysql"),
				),
			},
		},
	})
}

func TestAccDatabaseDatastoreCapabilitiesDataSource_migrateToFramework(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"vkcs": {
						VersionConstraint: "0.2.2",
						Source:            "vk-cs/vkcs",
					},
				},
				Config: testAccDatabaseDatastoreCapabilitiesDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_db_datastore_capabilities.capabilities", "datastore_name", "mysql"),
				),
			},
			{
				ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
				Config:                   testAccDatabaseDatastoreCapabilitiesDataSourceBasic,
				PlanOnly:                 true,
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
