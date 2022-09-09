package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKeyManagerContainerDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerDataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_container.container_1", "id",
						"vkcs_keymanager_container.container_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_container.container_1", "secret_refs",
						"vkcs_keymanager_container.container_1", "secret_refs"),
					resource.TestCheckResourceAttr(
						"data.vkcs_keymanager_container.container_1", "secret_refs.#", "3"),
				),
			},
		},
	})
}

func TestAccKeyManagerContainerDataSource_acls(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerDataSourceAcls(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_container.container_1", "id",
						"vkcs_keymanager_container.container_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_container.container_1", "secret_refs",
						"vkcs_keymanager_container.container_1", "secret_refs"),
					resource.TestCheckResourceAttr(
						"data.vkcs_keymanager_container.container_1", "secret_refs.#", "3"),
					resource.TestCheckResourceAttr("data.vkcs_keymanager_container.container_1", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("data.vkcs_keymanager_container.container_1", "acl.0.read.0.users.#", "2"),
				),
			},
		},
	})
}

func testAccKeyManagerContainerDataSourceBasic() string {
	return fmt.Sprintf(`
%s

data "vkcs_keymanager_container" "container_1" {
  name = vkcs_keymanager_container.container_1.name
}
`, testAccKeyManagerContainerBasic())
}

func testAccKeyManagerContainerDataSourceAcls() string {
	return fmt.Sprintf(`
%s

data "vkcs_keymanager_container" "container_1" {
  name = vkcs_keymanager_container.container_1.name
}
`, testAccKeyManagerContainerAcls())
}
