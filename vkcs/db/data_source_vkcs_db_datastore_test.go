package db_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDatabaseDatastoreDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseDatastoreDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_db_datastore.datastore", "name", "mysql"),
					resource.TestMatchResourceAttr("data.vkcs_db_datastore.datastore", "versions.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

func TestAccDatabaseDatastoreDataSource_migrateToFramework(t *testing.T) {
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
				Config: testAccDatabaseDatastoreDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_db_datastore.datastore", "name", "mysql"),
					resource.TestMatchResourceAttr("data.vkcs_db_datastore.datastore", "versions.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
			{
				ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
				Config:                   testAccDatabaseDatastoreDataSourceConfig,
				PlanOnly:                 true,
			},
		},
	})
}

const testAccDatabaseDatastoreDataSourceConfig = `
data "vkcs_db_datastore" "datastore" {
	name = "mysql"
}
`
