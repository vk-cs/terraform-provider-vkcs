package db_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDatabaseDatastoresDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseDatastoresDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.vkcs_db_datastores.datastores", "datastores.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttrSet("data.vkcs_db_datastores.datastores", "datastores.0.id"),
					resource.TestCheckResourceAttrSet("data.vkcs_db_datastores.datastores", "datastores.0.name"),
				),
			},
		},
	})
}

func TestAccDatabaseDatastoresDataSource_migrateToFramework(t *testing.T) {
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
				Config: testAccDatabaseDatastoresDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.vkcs_db_datastores.datastores", "datastores.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
			{
				ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
				Config:                   testAccDatabaseDatastoresDataSourceConfig,
				PlanOnly:                 true,
			},
		},
	})
}

const testAccDatabaseDatastoresDataSourceConfig = `
data "vkcs_db_datastores" "datastores" {}
`
