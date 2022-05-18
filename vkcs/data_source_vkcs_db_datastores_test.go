package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatabaseDataSourceDatastores_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDatabase(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseDataSourceDatastoresBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_db_datastores.datastores", "datastores.#", "1"),
					resource.TestCheckResourceAttr("data.vkcs_db_datastores.datastores", "datastores.0.name", osDBDatastoreType),
				),
			},
		},
	})
}

var testAccDatabaseDataSourceDatastoresBasic = fmt.Sprintf(`
data "vkcs_db_datastores" "datastores" {
	filter = "%s"
}
`, osDBDatastoreType)
