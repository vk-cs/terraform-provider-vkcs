package keymanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKeyManagerContainerDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
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

func TestAccKeymanagerContainerDataSource_migrateToFramework(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"vkcs": {
						VersionConstraint: "0.3.0",
						Source:            "vk-cs/vkcs",
					},
				},
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
			{
				ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
				Config:                   acctest.AccTestRenderConfig(testAccKeyManagerContainerDataSourceBasic, map[string]string{"TestAccKeyManagerContainerBasic": acctest.AccTestRenderConfig(testAccKeyManagerContainerBasic, map[string]string{"TestAccKeyManagerContainer": testAccKeyManagerContainer})}),
				PlanOnly:                 true,
			},
		},
	})
}

func TestAccKeyManagerContainerDataSource_acls(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
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
