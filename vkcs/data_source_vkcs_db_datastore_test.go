package vkcs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatabaseDatastoreDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
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

const testAccDatabaseDatastoreDataSourceConfig = `
data "vkcs_db_datastore" "datastore" {
	name = "mysql"
}
`
