package baremetal_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccBareMetalFlavorsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFlavorsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						resourceName := "data.vkcs_baremetal_flavors.basic"
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("root module has no resource called %s", resourceName)
						}

						flavors, ok := rs.Primary.Attributes["flavors.#"]
						if !ok {
							return fmt.Errorf("flavors attribute is missing")
						}

						flavorsQuantity, err := strconv.Atoi(flavors)
						if err != nil {
							return fmt.Errorf("error parsing flavors (%s) into integer: %s", flavors, err)
						}

						if flavorsQuantity == 0 {
							return fmt.Errorf("no flavors found, this is probably a bug")
						}

						return nil
					},
				),
			},
		},
	})
}

const testAccFlavorsDataSourceBasic = `
data "vkcs_baremetal_flavors" "basic" {}
`
