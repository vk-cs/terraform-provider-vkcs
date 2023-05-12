package db_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDatabaseDatastoresDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
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
