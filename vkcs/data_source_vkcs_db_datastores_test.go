package vkcs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatabaseDatastoresDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
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

const testAccDatabaseDatastoresDataSourceConfig = `
data "vkcs_db_datastores" "datastores" {}
`
