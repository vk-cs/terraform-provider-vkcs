package keymanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKeyManagerContainerDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKeyManagerContainerDataSourceBasic, map[string]string{"TestAccKeyManagerContainerBasic": acctest.AccTestRenderConfig(testAccKeyManagerContainerBasic, map[string]string{"TestAccKeyManagerContainer": testAccKeyManagerContainer})}),
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
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKeyManagerContainerDataSourceAcls, map[string]string{"TestAccKeyManagerContainerAcls": acctest.AccTestRenderConfig(testAccKeyManagerContainerAcls, map[string]string{"TestAccKeyManagerContainer": testAccKeyManagerContainer})}),
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

const testAccKeyManagerContainerDataSourceBasic = `
{{.TestAccKeyManagerContainerBasic}}

data "vkcs_keymanager_container" "container_1" {
  name = vkcs_keymanager_container.container_1.name
}
`

const testAccKeyManagerContainerDataSourceAcls = `
{{.TestAccKeyManagerContainerAcls}}

data "vkcs_keymanager_container" "container_1" {
  name = vkcs_keymanager_container.container_1.name
}
`
