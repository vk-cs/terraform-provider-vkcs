package keymanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKeyManagerContainer_importBasic(t *testing.T) {
	resourceName := "vkcs_keymanager_container.container_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKeyManagerContainerBasic, map[string]string{"TestAccKeyManagerContainer": testAccKeyManagerContainer}),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKeyManagerContainer_importACLs(t *testing.T) {
	resourceName := "vkcs_keymanager_container.container_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKeyManagerContainerAcls, map[string]string{"TestAccKeyManagerContainer": testAccKeyManagerContainer}),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
