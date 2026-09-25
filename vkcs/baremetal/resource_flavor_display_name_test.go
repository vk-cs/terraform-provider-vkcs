package baremetal_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
)

func TestAccBareMetalFlavorDisplayNameResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccBareMetalFlavorDisplayNamePreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFlavorDisplayNameResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_baremetal_flavor_display_name.basic", "display_name", "TerraformAccTestBasic"),
				),
			},
			{
				Config: testAccFlavorDisplayNameResourceUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_baremetal_flavor_display_name.basic", "display_name", "TerraformAccTestUpdated"),
				),
			},
		},
	})
}

func testAccBareMetalFlavorDisplayNamePreCheck(t *testing.T) {
	t.Helper()

	opts := clients.ConfigOpts{}
	if _, err := opts.LoadAndValidate(); err != nil {
		t.Fatalf("Error loading VKCS configuration for acceptance test: %s", err)
	}
}

const testAccFlavorDisplayNameResourceBasic = `
data "vkcs_baremetal_flavor" "basic" {
  name = "test_flavor2"
}

resource "vkcs_baremetal_flavor_display_name" "basic" {
  id           = data.vkcs_baremetal_flavor.basic.id
  display_name = "TerraformAccTestBasic"
}
`

const testAccFlavorDisplayNameResourceUpdated = `
data "vkcs_baremetal_flavor" "basic" {
  name = "test_flavor2"
}

resource "vkcs_baremetal_flavor_display_name" "basic" {
  id           = data.vkcs_baremetal_flavor.basic.id
  display_name = "TerraformAccTestUpdated"
}
`
