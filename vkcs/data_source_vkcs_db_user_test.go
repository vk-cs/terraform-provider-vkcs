package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDatabaseDataSourceUser_basic(t *testing.T) {
	resourceName := "vkcs_db_user.basic"
	datasourceName := "data.vkcs_db_user.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatabaseUserBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDatabaseUserID(datasourceName),
					resource.TestCheckResourceAttrPair(resourceName, "name", datasourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceDatabaseUserID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find user data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("User data source ID not set")
		}

		return nil
	}
}

var testAccDataSourceDatabaseUserBasic = fmt.Sprintf(`
%s

data "vkcs_db_user" "basic" {
	id = "${vkcs_db_user.basic.id}"
}
`, testAccDatabaseUserBasic)
