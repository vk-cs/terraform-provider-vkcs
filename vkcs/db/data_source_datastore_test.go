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

const testAccDatabaseDatastoreDataSourceConfig = `
data "vkcs_db_datastore" "datastore" {
	name = "mysql"
}
`
