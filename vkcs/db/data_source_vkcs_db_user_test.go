package db_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDatabaseDataSourceUser_basic(t *testing.T) {
	resourceName := "vkcs_db_user.basic"
	datasourceName := "data.vkcs_db_user.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDataSourceDatabaseUserBasic, map[string]string{"TestAccDatabaseUserBasic": acctest.AccTestRenderConfig(testAccDatabaseUserBasic)}),
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

const testAccDataSourceDatabaseUserBasic = `
{{.TestAccDatabaseUserBasic}}

data "vkcs_db_user" "basic" {
	id = vkcs_db_user.basic.id
}
`
