package baremetal_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccBareMetalOSesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOSesDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						resourceName := "data.vkcs_baremetal_oses.basic"
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("root module has no resource called %s", resourceName)
						}

						oses, ok := rs.Primary.Attributes["oses.#"]
						if !ok {
							return fmt.Errorf("oses attribute is missing")
						}

						osesQuantity, err := strconv.Atoi(oses)
						if err != nil {
							return fmt.Errorf("error parsing oses (%s) into integer: %s", oses, err)
						}

						if osesQuantity == 0 {
							return fmt.Errorf("no oses found, this is probably a bug")
						}

						return nil
					},
				),
			},
		},
	})
}

const testAccOSesDataSourceBasic = `
data "vkcs_baremetal_oses" "basic" {}
`
