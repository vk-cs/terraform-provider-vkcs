package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDatabaseDataSourceInstance_basic(t *testing.T) {
	resourceName := "vkcs_db_instance.basic"
	datasourceName := "data.vkcs_db_instance.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDatabase(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatabaseInstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDatabaseInstanceID(datasourceName),
					resource.TestCheckResourceAttrPair(resourceName, "name", datasourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceDatabaseInstanceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find instance data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Instance data source ID not set")
		}

		return nil
	}
}

var testAccDataSourceDatabaseInstanceBasic = fmt.Sprintf(`
%s

data "vkcs_db_instance" "basic" {
	id = "${vkcs_db_instance.basic.id}"
}
`, testAccDatabaseInstanceBasic)
